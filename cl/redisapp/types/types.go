package types

type ExecutionHead struct {
	BlockHeight uint64
	BlockHash   []byte
	BlockTime   uint64
}

type BuildStep int

const (
	StepBuildBlock BuildStep = iota
	StepFinalizeBlock
)

func (s BuildStep) String() string {
	switch s {
	case StepBuildBlock:
		return "BuildBlock"
	case StepFinalizeBlock:
		return "FinalizeBlock"
	default:
		return "Unknown"
	}
}

type BlockBuildState struct {
	CurrentStep      BuildStep `json:"current_step"`
	PayloadID        string    `json:"payload_id,omitempty"`
	ExecutionPayload string    `json:"execution_payload,omitempty"`
}

type RedisMsgType string

const (
	RedisMsgTypePending RedisMsgType = "0"
	RedisMsgTypeNew     RedisMsgType = ">"
)
