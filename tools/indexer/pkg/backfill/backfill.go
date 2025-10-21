package backfill

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/primev/mev-commit/tools/indexer/pkg/beacon"
	"github.com/primev/mev-commit/tools/indexer/pkg/config"
	"github.com/primev/mev-commit/tools/indexer/pkg/database"
)

type SlotData struct {
	Slot            int64
	BlockNumber     int64
	ValidatorPubkey []byte
	ProposerIdx     *int64
}

// RatedSlotResponse represents the response from Rated Network API blocks endpoint
type RatedSlotResponse struct {
	Epoch                      int64    `json:"epoch"`
	ConsensusSlot              int64    `json:"consensusSlot"`
	ConsensusBlockRoot         string   `json:"consensusBlockRoot"`
	ExecutionBlockNumber       int64    `json:"executionBlockNumber"`
	ExecutionBlockHash         string   `json:"executionBlockHash"`
	ValidatorIndex             int64    `json:"validatorIndex"`
	FeeRecipient               string   `json:"feeRecipient"`
	TotalTransactions          int      `json:"totalTransactions"`
	TotalGasUsed               int64    `json:"totalGasUsed"`
	BaseFeePerGas              int64    `json:"baseFeePerGas"`
	BaselineMev                int64    `json:"baselineMev"`
	BaselineMevWei             string   `json:"baselineMevWei"`
	ExecutionRewards           int64    `json:"executionRewards"`
	ExecutionRewardsWei        string   `json:"executionRewardsWei"`
	ConsensusRewards           int64    `json:"consensusRewards"`
	TotalRewards               int64    `json:"totalRewards"`
	BlockTimestamp             string   `json:"blockTimestamp"`
	Relays                     []string `json:"relays"`
	BlockBuilderPubkeys        []string `json:"blockBuilderPubkeys"`
	ExecutionProposerDuty      string   `json:"executionProposerDuty"`
	TotalPriorityFeesValidator int64    `json:"totalPriorityFeesValidator"`
}

// RatedValidatorResponse represents validator details from Rated API
type RatedValidatorResponse struct {
	ValidatorIndex  int64  `json:"validatorIndex"`
	ValidatorPubkey string `json:"validatorPubkey"`
	Pool            string `json:"pool"`
	Network         string `json:"network"`
}

// FetchSlotsBatch fetches a batch of slots from Rated API
func (r *RatedAPIClient) FetchSlotsBatch(ctx context.Context, startSlot, endSlot int64) ([]RatedSlotResponse, error) {
	if endSlot <= startSlot {
		return nil, nil
	}
	to := endSlot - 1
	if to < startSlot {
		to = startSlot
	}
	pageLimit := endSlot - startSlot
	if pageLimit > 1000 {
		pageLimit = 1000
	}
	// to := endSlot - 1
	nextURL := fmt.Sprintf("%s/blocks?from=%d&to=%d&limit=%d&offset=0", r.baseURL, startSlot, to, pageLimit)

	req, err := retryablehttp.NewRequestWithContext(ctx, "GET", nextURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.apiKey))
	req.Header.Set("Accept", "application/json")

	resp, err := r.httpc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var env map[string]json.RawMessage
	if err := json.Unmarshal(body, &env); err != nil {
		snip := string(body)
		if len(snip) > 512 {
			snip = snip[:512] + "..."
		}
		return nil, fmt.Errorf("decoding envelope: %w; body starts: %s", err, snip)
	}
	resRaw, ok := env["results"]
	if !ok || len(resRaw) == 0 {
		return []RatedSlotResponse{}, nil
	}

	var flat []json.RawMessage
	if err := json.Unmarshal(resRaw, &flat); err != nil {

		var nested [][]json.RawMessage
		if err2 := json.Unmarshal(resRaw, &nested); err2 != nil {
			snip := string(resRaw)
			if len(snip) > 512 {
				snip = snip[:512] + "..."
			}
			return nil, fmt.Errorf("decoding results: %v; nested err: %v; results starts: %s", err, err2, snip)
		}

		for _, grp := range nested {
			flat = append(flat, grp...)
		}
	}

	out := make([]RatedSlotResponse, 0, len(flat))
	for _, raw := range flat {
		var m map[string]any
		if err := json.Unmarshal(raw, &m); err != nil {
			continue
		}
		var s RatedSlotResponse
		s.Epoch = asInt64(m["epoch"])
		s.ConsensusSlot = asInt64(m["consensusSlot"])
		s.ConsensusBlockRoot = asString(m["consensusBlockRoot"])
		s.ExecutionBlockNumber = asInt64(m["executionBlockNumber"])
		s.ExecutionBlockHash = asString(m["executionBlockHash"])
		s.ValidatorIndex = asInt64(m["validatorIndex"])
		s.FeeRecipient = asString(m["feeRecipient"])
		s.TotalTransactions = int(asInt64(m["totalTransactions"]))
		s.TotalGasUsed = asInt64(m["totalGasUsed"])
		s.BaseFeePerGas = asInt64(m["baseFeePerGas"])
		s.BaselineMev = asInt64(m["baselineMev"])
		s.BaselineMevWei = asDecimalString(m["baselineMevWei"])
		s.ExecutionRewards = asInt64(m["executionRewards"])
		s.ExecutionRewardsWei = asDecimalString(m["executionRewardsWei"])
		s.ConsensusRewards = asInt64(m["consensusRewards"])
		s.TotalRewards = asInt64(m["totalRewards"])
		s.BlockTimestamp = asString(m["blockTimestamp"])
		s.Relays = asStringSlice(m["relays"])
		s.BlockBuilderPubkeys = asStringSlice(m["blockBuilderPubkeys"])
		s.ExecutionProposerDuty = asString(m["executionProposerDuty"])
		s.TotalPriorityFeesValidator = asInt64(m["totalPriorityFeesValidator"])
		out = append(out, s)
	}

	return out, nil
}

func (q *QuickNodeClient) FetchValidatorPubkeysBatch(
	ctx context.Context,
	indices []int64,
	missingOut *[]int64,
) (map[int64]string, error) {
	out := make(map[int64]string, len(indices))
	if len(indices) == 0 {
		return out, nil
	}

	// de-dupe and freeze order
	seen := make(map[int64]struct{}, len(indices))
	uniq := make([]int64, 0, len(indices))
	for _, idx := range indices {
		if _, ok := seen[idx]; !ok {
			seen[idx] = struct{}{}
			uniq = append(uniq, idx)
		}
	}

	type chunkRes struct {
		data map[int64]string
		err  error
		seen []int64
	}

	chunks := make([][]int64, 0, (len(uniq)+q.chunkSize-1)/q.chunkSize)
	for i := 0; i < len(uniq); i += q.chunkSize {
		j := i + q.chunkSize
		if j > len(uniq) {
			j = len(uniq)
		}
		chunks = append(chunks, uniq[i:j])
	}

	sem := make(chan struct{}, q.concurrent)
	resCh := make(chan chunkRes, len(chunks))
	var wg sync.WaitGroup

	for _, chunk := range chunks {
		sem <- struct{}{}
		wg.Add(1)
		ch := chunk
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			// Build URL with repeated id params.
			ids := make([]string, len(ch))
			for i, v := range ch {
				ids[i] = strconv.FormatInt(v, 10)
			}
			idQuery := "id=" + strings.Join(ids, "&id=")
			fullURL := fmt.Sprintf("%s/eth/v1/beacon/states/head/validators?%s", q.base, idQuery)

			// Per-request timeout
			reqCtx, cancel := context.WithTimeout(ctx, q.timeout)
			defer cancel()

			req, err := retryablehttp.NewRequestWithContext(reqCtx, "GET", fullURL, nil)
			if err != nil {
				resCh <- chunkRes{err: err, seen: ch}
				return
			}
			req.Header.Set("accept", "application/json")

			resp, err := q.httpc.Do(req)
			if err != nil {
				resCh <- chunkRes{err: err, seen: ch}
				return
			}
			defer resp.Body.Close()

			body, rerr := io.ReadAll(resp.Body)
			if rerr != nil {
				resCh <- chunkRes{err: rerr, seen: ch}
				return
			}
			if resp.StatusCode != 200 {
				snip := string(body)
				if len(snip) > 240 {
					snip = snip[:240] + "..."
				}
				resCh <- chunkRes{err: fmt.Errorf("quicknode %d: %s", resp.StatusCode, snip), seen: ch}
				return
			}

			var parsed struct {
				Data []struct {
					Index     string `json:"index"`
					Validator struct {
						Pubkey string `json:"pubkey"`
					} `json:"validator"`
				} `json:"data"`
			}
			if err := json.Unmarshal(body, &parsed); err != nil {
				resCh <- chunkRes{err: fmt.Errorf("parse quicknode: %w", err), seen: ch}
				return
			}

			m := make(map[int64]string, len(parsed.Data))
			for _, row := range parsed.Data {
				idx, _ := strconv.ParseInt(row.Index, 10, 64)
				pk := strings.TrimSpace(row.Validator.Pubkey)
				if idx > 0 && pk != "" {
					// ensure 0x prefix
					if !strings.HasPrefix(pk, "0x") {
						pk = "0x" + pk
					}
					m[idx] = pk
				}
			}
			resCh <- chunkRes{data: m, seen: ch}
		}()
	}

	wg.Wait()
	close(resCh)

	errs := make([]error, 0)
	seenAll := make(map[int64]struct{}, len(uniq))
	for _, idx := range uniq {
		seenAll[idx] = struct{}{}
	}
	found := make(map[int64]struct{})

	for r := range resCh {
		if r.err != nil {
			errs = append(errs, r.err)
			continue
		}
		for k, v := range r.data {
			out[k] = v
			found[k] = struct{}{}
		}
	}

	if missingOut != nil {
		miss := make([]int64, 0)
		for idx := range seenAll {
			if _, ok := found[idx]; !ok {
				miss = append(miss, idx)
			}
		}
		*missingOut = miss
	}

	if len(errs) > 0 {
		return out, fmt.Errorf("quicknode batch had %d error(s), e.g. %v", len(errs), errs[0])
	}
	return out, nil
}

// RunAllWithRatedAPI runs the optimized backfill using Rated Network API
func RunAllWithRatedAPI(ctx context.Context, db *database.DB, httpc *retryablehttp.Client,
	cfg *config.Config,
) error {
	qn := NewQuickNodeClient(httpc, cfg.QuickNodeBase)
	logger := slog.With("component", "backfill-rated-optimized")
	logger.Info("Starting optimized streaming backfill with Rated Network API")

	if err := ctx.Err(); err != nil {
		return err
	}

	ratedClient := NewRatedAPIClient(httpc, cfg.RatedAPIKey)

	lastSlotNumber, err := db.GetMinSlotNumber(ctx)
	if err != nil {
		lastSlotNumber = 0
	}

	startSlot := lastSlotNumber - cfg.BackfillLookback
	if startSlot < 0 {
		startSlot = 0
	}

	ratedBatchSize := int64(cfg.BackfillBatch)
	totalSlots := lastSlotNumber - startSlot
	totalBatches := (totalSlots + ratedBatchSize - 1) / ratedBatchSize

	logger.Info("Starting Rated API optimized backfill",
		"start_slot", startSlot,
		"end_slot", lastSlotNumber,
		"total_slots", totalSlots,
		"batch_size", ratedBatchSize,
		"total_batches", totalBatches,
		"estimated_duration", fmt.Sprintf("~%d minutes", totalBatches*3/2),
	)

	var totalProcessed int64
	var totalFailed int64
	startTime := time.Now()

	// Process in batches
	for batchIdx := int64(0); batchIdx < totalBatches; batchIdx++ {
		batchStart := startSlot + batchIdx*ratedBatchSize
		batchEnd := batchStart + ratedBatchSize
		if batchEnd > lastSlotNumber {
			batchEnd = lastSlotNumber
		}

		logger.Info("Processing batch",
			"batch", batchIdx+1, "of", totalBatches,
			"range", fmt.Sprintf("[%d,%d)", batchStart, batchEnd),
			"progress", fmt.Sprintf("%.1f%%", float64(batchIdx)/float64(totalBatches)*100),
		)

		// Fetch slots batch
		tBatch := time.Now()
		slots, err := ratedClient.FetchSlotsBatch(ctx, batchStart, batchEnd)
		if err != nil {
			logger.Error("Failed to fetch batch", "batch", batchIdx, "error", err)

			// Retry once with backoff
			time.Sleep(5 * time.Second)
			slots, err = ratedClient.FetchSlotsBatch(ctx, batchStart, batchEnd)
			if err != nil {
				logger.Error("Retry failed, skipping batch", "batch", batchIdx)
				totalFailed += batchEnd - batchStart
				continue
			}
		}

		valid := make([]RatedSlotResponse, 0, len(slots))
		idxSet := make(map[int64]struct{})
		for _, s := range slots {
			if s.ExecutionBlockNumber == 0 {
				continue
			}
			valid = append(valid, s)
			if s.ValidatorIndex > 0 {
				idxSet[s.ValidatorIndex] = struct{}{}
			}
		}
		indices := make([]int64, 0, len(idxSet))
		for i := range idxSet {
			indices = append(indices, i)
		}

		missing := []int64{}
		// batch-fetch pubkeys from QuickNode (chunked)
		pubByIdx, err := qn.FetchValidatorPubkeysBatch(ctx, indices, &missing)
		if err != nil {
			logger.Error("QuickNode fetch validators failed", "error", err, "indices", len(indices))
			pubByIdx = map[int64]string{}
		}

		execInfos := make([]*beacon.ExecInfo, 0, len(valid))
		for _, s := range valid {
			ei := ConvertRatedToExecInfo(s)
			if pk, ok := pubByIdx[s.ValidatorIndex]; ok && pk != "" {
				if b, err := hex.DecodeString(strings.TrimPrefix(pk, "0x")); err == nil {
					ei.ValidatorPubkey = b
				}
			}
			execInfos = append(execInfos, ei)
		}

		// upsert blocks (StarRocks PK)
		if err := db.BatchUpsertBlocksFromExec(ctx, execInfos); err != nil {
			logger.Error("Batch block insert failed", "error", err)
			for _, ei := range execInfos {
				if err := db.UpsertBlockFromExec(ctx, ei); err != nil {
					logger.Error("Individual block insert failed", "slot", ei.Slot, "error", err)
				}
			}
		}

		// Process entire batch
		pairs := make([]struct {
			Slot   int64
			Pubkey string
		}, 0, len(execInfos))
		for _, ei := range execInfos {
			if len(ei.ValidatorPubkey) == 0 {
				continue
			}
			hexpk := hex.EncodeToString(ei.ValidatorPubkey)
			if !strings.HasPrefix(hexpk, "0x") {
				hexpk = "0x" + hexpk
			}
			pairs = append(pairs, struct {
				Slot   int64
				Pubkey string
			}{Slot: ei.Slot, Pubkey: hexpk})
		}
		if err := db.UpsertBlockPubkeysDirect(ctx, pairs); err != nil {
			logger.Error("upsert block pubkeys direct failed", "error", err)
		}

		batchDuration := time.Since(tBatch)

		logger.Info("Batch completed",
			"batch", batchIdx+1,
			"slots_in_batch", len(slots),
			"duration", batchDuration,
			"slots_per_second", float64(len(slots))/batchDuration.Seconds(),
		)

		// Check for cancellation
		if ctx.Err() != nil {
			logger.Info("Backfill cancelled")
			break
		}
	}

	totalDuration := time.Since(startTime)
	logger.Info("Backfill completed",
		"total_processed", totalProcessed,
		"total_failed", totalFailed,
		"total_duration", totalDuration,
		"average_slots_per_second", float64(totalProcessed)/totalDuration.Seconds(),
		"success_rate", fmt.Sprintf("%.2f%%", float64(totalProcessed)/float64(totalSlots)*100),
	)

	return nil
}
