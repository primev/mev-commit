package backfill

import (
	"encoding/json"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/primev/mev-commit/tools/indexer/pkg/beacon"
)

func asInt64(v any) int64 {
	switch t := v.(type) {
	case float64:
		return int64(t)
	case string:
		if t == "" {
			return 0
		}
		if strings.HasPrefix(t, "0x") || strings.HasPrefix(t, "0X") {
			n, _ := strconv.ParseInt(t[2:], 16, 64)
			return n
		}
		n, _ := strconv.ParseInt(t, 10, 64)
		return n
	default:
		return 0
	}
}

func asString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func asStringSlice(v any) []string {
	if v == nil {
		return nil
	}
	switch arr := v.(type) {
	case []any:
		out := make([]string, 0, len(arr))
		for _, it := range arr {
			out = append(out, asString(it))
		}
		return out
	case []string:
		return arr
	default:
		return nil
	}
}

func asDecimalString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		// convert float64 (JSON numbers) to integer string without decimals
		return strconv.FormatInt(int64(t), 10)
	case nil:
		return ""
	default:
		return asString(v)
	}
}

func ConvertRatedToExecInfo(slot RatedSlotResponse) *beacon.ExecInfo {
	ei := &beacon.ExecInfo{
		BlockNumber: slot.ExecutionBlockNumber,
		Slot:        slot.ConsensusSlot,
	}

	// Validator index
	if slot.ValidatorIndex > 0 {
		idx := slot.ValidatorIndex
		ei.ProposerIdx = &idx
	}

	// Parse timestamp
	if slot.BlockTimestamp != "" {
		if t, err := time.Parse("2006-01-02T15:04:05", slot.BlockTimestamp); err == nil {
			utc := t.UTC()
			ei.Timestamp = &utc
		}
	}

	// Relay information
	if len(slot.Relays) > 0 {
		relayTag := slot.Relays[0]
		ei.RelayTag = &relayTag
	}

	// Builder public key
	if len(slot.BlockBuilderPubkeys) > 0 {
		builderPubkey := slot.BlockBuilderPubkeys[0]
		ei.BuilderPublicKey = &builderPubkey
	}

	// Fee recipient
	if slot.FeeRecipient != "" {
		cleaned := cleanHexString(slot.FeeRecipient)
		ei.FeeRecipient = &cleaned
		ei.ProposerFeeRecHex = &cleaned
	}

	// MEV reward in ETH
	if slot.ExecutionRewardsWei != "" && slot.ExecutionRewardsWei != "0" {
		if mevEth := weiToEth(slot.ExecutionRewardsWei); mevEth > 0 {
			ei.MevRewardEth = &mevEth
			ei.ProposerRewardEth = &mevEth
		}
	} else if slot.BaselineMevWei != "" && slot.BaselineMevWei != "0" {
		if mevEth := weiToEth(slot.BaselineMevWei); mevEth > 0 {
			ei.MevRewardEth = &mevEth
		}
	}

	return ei
}

func cleanHexString(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "\\\\x") {
		return "0x" + s[3:]
	}
	if strings.HasPrefix(s, "\\x") {
		return "0x" + s[2:]
	}
	if !strings.HasPrefix(s, "0x") {
		return "0x" + s
	}
	return s
}

func weiToEth(weiStr string) float64 {
	weiStr = strings.TrimSpace(weiStr)
	if weiStr == "" || weiStr == "0" {
		return 0
	}

	wei, ok := new(big.Int).SetString(weiStr, 10)
	if !ok {
		if strings.HasPrefix(weiStr, "0x") {
			wei, ok = new(big.Int).SetString(weiStr[2:], 16)
			if !ok {
				return 0
			}
		} else {
			return 0
		}
	}

	eth := new(big.Rat).SetFrac(wei, big.NewInt(1e18))
	f, _ := eth.Float64()
	return f
}
