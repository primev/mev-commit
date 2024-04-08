package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/primevprotocol/mev-commit/oracle/pkg/settler"
	"github.com/primevprotocol/mev-commit/oracle/pkg/updater"
)

var settlementType = `
DO $$ BEGIN
    CREATE TYPE settlement_type AS ENUM ('reward', 'slash', 'return');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;`

var settlementsTable = `
CREATE TABLE IF NOT EXISTS settlements (
    commitment_index BYTEA PRIMARY KEY,
    transaction TEXT,
    block_number BIGINT,
    builder_address BYTEA,
    type settlement_type,
    amount NUMERIC(24, 0),
    bid_id BYTEA,
    chainhash BYTEA,
    nonce BIGINT,
    settled BOOLEAN,
    decay_percentage BIGINT
);`

var winnersTable = `
CREATE TABLE IF NOT EXISTS winners (
    block_number BIGINT PRIMARY KEY,
    builder_address BYTEA,
    processed BOOLEAN
);`

type Store struct {
	db       *sql.DB
	winnerT  chan struct{}
	settlerT chan struct{}
	returnT  chan struct{}
}

func NewStore(db *sql.DB) (*Store, error) {
	for _, table := range []string{settlementType, settlementsTable, winnersTable} {
		_, err := db.Exec(table)
		if err != nil {
			return nil, err
		}
	}

	return &Store{
		db:       db,
		winnerT:  make(chan struct{}),
		settlerT: make(chan struct{}),
		returnT:  make(chan struct{}),
	}, nil
}

func (s *Store) triggerWinner() {
	select {
	case s.winnerT <- struct{}{}:
	default:
	}
}

func (s *Store) triggerSettler() {
	select {
	case s.settlerT <- struct{}{}:
	default:
	}
}

func (s *Store) triggerReturn() {
	select {
	case s.returnT <- struct{}{}:
	default:
	}
}

func (s *Store) RegisterWinner(ctx context.Context, blockNum int64, winner string) error {
	insertStr := "INSERT INTO winners (block_number, builder_address, processed) VALUES ($1, $2, $3)"

	_, err := s.db.ExecContext(ctx, insertStr, blockNum, winner, false)
	if err != nil {
		return err
	}
	s.triggerWinner()
	return nil
}

func (s *Store) SubscribeWinners(ctx context.Context) <-chan updater.BlockWinner {
	resChan := make(chan updater.BlockWinner)
	go func() {
		defer close(resChan)

	RETRY:
		for {
			results, err := s.db.QueryContext(
				ctx,
				"SELECT block_number, builder_address FROM winners WHERE processed = false",
			)
			if err != nil {
				return
			}
			for results.Next() {
				var bWinner updater.BlockWinner
				err = results.Scan(&bWinner.BlockNumber, &bWinner.Winner)
				if err != nil {
					_ = results.Close()
					continue RETRY
				}
				select {
				case <-ctx.Done():
					_ = results.Close()
					return
				case resChan <- bWinner:
				}
			}
			_ = results.Close()

			select {
			case <-ctx.Done():
				return
			case <-s.winnerT:
			}
		}
	}()

	return resChan
}

func (s *Store) UpdateComplete(ctx context.Context, blockNum int64) error {
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE winners SET processed = true WHERE block_number = $1",
		blockNum,
	)
	if err != nil {
		return err
	}
	s.triggerSettler()
	return nil
}

func (s *Store) AddSettlement(
	ctx context.Context,
	commitmentIdx []byte,
	txHash string,
	blockNum int64,
	amount uint64,
	builder string,
	bidID []byte,
	settlementType settler.SettlementType,
	decayPercentage int64,
) error {
	columns := []string{
		"commitment_index",
		"transaction",
		"block_number",
		"builder_address",
		"type",
		"amount",
		"bid_id",
		"settled",
		"chainhash",
		"nonce",
		"decay_percentage",
	}
	values := []interface{}{
		commitmentIdx,
		txHash,
		blockNum,
		builder,
		settlementType,
		amount,
		bidID,
		false,
		nil,
		0,
		decayPercentage,
	}
	placeholder := make([]string, len(values))
	for i := range columns {
		placeholder[i] = fmt.Sprintf("$%d", i+1)
	}

	insertStr := fmt.Sprintf(
		"INSERT INTO settlements (%s) VALUES (%s)",
		strings.Join(columns, ", "),
		strings.Join(placeholder, ", "),
	)

	_, err := s.db.ExecContext(ctx, insertStr, values...)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) SubscribeSettlements(ctx context.Context) <-chan settler.Settlement {
	resChan := make(chan settler.Settlement)

	go func() {
		defer close(resChan)

	RETRY:
		for {
			queryStr := `
				SELECT commitment_index, transaction, block_number, builder_address, amount, bid_id, type, decay_percentage
				FROM settlements
				WHERE settled = false AND chainhash IS NULL AND type != 'return'
				ORDER BY block_number ASC`

			results, err := s.db.QueryContext(ctx, queryStr)
			if err != nil {
				return
			}

			for results.Next() {
				var s settler.Settlement

				err = results.Scan(
					&s.CommitmentIdx,
					&s.TxHash,
					&s.BlockNum,
					&s.Builder,
					&s.Amount,
					&s.BidID,
					&s.Type,
					&s.DecayPercentage,
				)
				if err != nil {
					_ = results.Close()
					continue RETRY
				}

				select {
				case <-ctx.Done():
					_ = results.Close()
					return
				case resChan <- s:
				}
			}

			_ = results.Close()

			select {
			case <-ctx.Done():
				return
			case <-s.settlerT:
			}
		}
	}()

	return resChan
}

func (s *Store) SubscribeReturns(ctx context.Context, limit int) <-chan settler.Return {
	resChan := make(chan settler.Return)

	go func() {
		defer close(resChan)

	RETRY:
		for {
			queryStr := `
				SELECT DISTINCT bid_id, block_number
				FROM settlements
				WHERE settled = false AND chainhash IS NULL AND type = 'return'
					AND block_number < (SELECT MAX(block_number) FROM settlements WHERE settled = true)
				ORDER BY block_number ASC`

			results, err := s.db.QueryContext(ctx, queryStr)
			if err != nil {
				fmt.Println("error", err)
				return
			}

			returns := make([][]byte, 0, limit)

			copyReturns := func() [][32]byte {
				bidIDs := make([][32]byte, len(returns))
				for idx, bidID := range returns {
					bidIDs[idx] = [32]byte{}
					copy(bidIDs[idx][:], bidID)
				}
				return bidIDs
			}

			for results.Next() {
				var r []byte
				err = results.Scan(&r, new(int64))
				if err != nil {
					_ = results.Close()
					continue RETRY
				}

				returns = append(returns, r)
				if len(returns) == limit {
					select {
					case <-ctx.Done():
						_ = results.Close()
						return
					case resChan <- settler.Return{BidIDs: copyReturns()}:
						returns = returns[:0]
					}
				}
			}

			if len(returns) > 0 {
				select {
				case <-ctx.Done():
					_ = results.Close()
					return
				case resChan <- settler.Return{BidIDs: copyReturns()}:
					returns = returns[:0]
				}
			}

			_ = results.Close()

			select {
			case <-ctx.Done():
				return
			case <-s.settlerT:
			}
		}
	}()

	return resChan
}

func (s *Store) SettlementInitiated(
	ctx context.Context,
	bidIDs [][]byte,
	txHash common.Hash,
	nonce uint64,
) error {
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE settlements SET chainhash = $1, nonce = $2 WHERE bid_id = ANY($3::BYTEA[])",
		txHash.Bytes(),
		nonce,
		pq.Array(bidIDs),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) MarkSettlementComplete(ctx context.Context, nonce uint64) (int, error) {
	result, err := s.db.ExecContext(
		ctx,
		"UPDATE settlements SET settled = true WHERE settled = false AND nonce < $1 AND chainhash IS NOT NULL",
		nonce,
	)
	if err != nil {
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	s.triggerReturn()
	return int(count), nil
}

func (s *Store) LastNonce() (int64, error) {
	var lastNonce int64
	err := s.db.QueryRow("SELECT MAX(nonce) FROM settlements").Scan(&lastNonce)
	if err != nil {
		return 0, err
	}
	return lastNonce, nil
}

func (s *Store) PendingTxnCount() (int, error) {
	var count int
	err := s.db.QueryRow(
		"SELECT COUNT(DISTINCT chainhash) FROM settlements WHERE chainhash IS NOT NULL AND settled = false",
	).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

type BlockInfo struct {
	BlockNumber     int64
	Builder         string
	NoOfCommitments int
	NoOfBids        int
	TotalAmount     sql.NullString
	NoOfRewards     int
	TotalRewards    sql.NullString
	NoOfSlashes     int
	TotalSlashes    sql.NullString
	NoOfSettlements int
}

func (s *Store) ProcessedBlocks(limit, offset int) ([]BlockInfo, error) {
	var blocks []BlockInfo
	rows, err := s.db.Query(`
		SELECT
			winners.block_number,
			winners.builder_address,
			COUNT(settlements.commitment_index) AS commitment_count,
			COUNT(DISTINCT settlements.bid_id) AS bid_count,
			(SELECT SUM(amount) FROM (
				SELECT DISTINCT ON (bid_id) bid_id, amount
				FROM settlements sub_settlements
				WHERE sub_settlements.block_number = winners.block_number
				ORDER BY bid_id, block_number
			) AS distinct_amounts) AS total_amount,
			COUNT(settlements.type = 'reward' OR NULL) AS reward_count,
			SUM(settlements.amount) FILTER (WHERE settlements.type = 'reward') AS total_rewards,
			COUNT(settlements.type = 'slash' OR NULL) AS slash_count,
			SUM(settlements.amount) FILTER (WHERE settlements.type = 'slash') AS total_slashes,
			COUNT(settlements.settled) FILTER (WHERE settlements.settled = true) AS settled_count
		FROM
			winners
		LEFT JOIN
			settlements ON settlements.block_number = winners.block_number
		WHERE
			winners.processed = true
		GROUP BY
			winners.block_number, winners.builder_address
		ORDER BY
			winners.block_number DESC
		LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b BlockInfo
		err := rows.Scan(
			&b.BlockNumber,
			&b.Builder,
			&b.NoOfCommitments,
			&b.NoOfBids,
			&b.TotalAmount,
			&b.NoOfRewards,
			&b.TotalRewards,
			&b.NoOfSlashes,
			&b.TotalSlashes,
			&b.NoOfSettlements,
		)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}

type CommitmentStats struct {
	TotalCount                int
	BidCount                  int
	RewardCount               int
	SlashCount                int
	SettlementsCompletedCount int
}

func (s *Store) CommitmentStats() (CommitmentStats, error) {
	var stats CommitmentStats
	err := s.db.QueryRow(`
		SELECT
			COUNT(*),
			COUNT(DISTINCT bid_id),
			COUNT(type = 'reward' OR NULL),
			COUNT(type = 'slash' OR NULL),
			COUNT(settled) FILTER (WHERE settled = true)
		FROM
			settlements
	`).Scan(
		&stats.TotalCount,
		&stats.BidCount,
		&stats.RewardCount,
		&stats.SlashCount,
		&stats.SettlementsCompletedCount,
	)
	if err != nil {
		return stats, err
	}
	return stats, nil
}
