package dex

const MAINNET_FACTORY_ADDR = "EQBfBWT7X2BHg9tXAxzhz2aKiNTU1tpt5NsiK0uSDW_YAJ67" // const

type StackItem struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// DataStructure represents the main JSON object
type DataStructure struct {
	Address string      `json:"address"`
	Method  string      `json:"method"`
	Stack   []StackItem `json:"stack"`
}

// ExecutionResult represents the JSON object
type ExecutionResult struct {
	GasUsed  int         `json:"gas_used"`
	ExitCode int         `json:"exit_code"`
	Stack    []StackItem `json:"stack"`
}
