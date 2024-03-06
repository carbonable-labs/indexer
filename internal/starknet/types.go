package starknet

import "time"

type Transaction struct {
	TransactionHash string   `json:"transaction_hash"`
	Version         string   `json:"version"`
	MaxFeePerGas    string   `json:"max_fee_per_gas"`
	Nonce           string   `json:"nonce"`
	SenderAddress   string   `json:"sender_address"`
	Type            string   `json:"type"`
	Signatures      []string `json:"signatures"`
	Calldata        []string `json:"calldata"`
}

type L2ToL1Message struct {
	FromAddress string   `json:"from_address"`
	ToAddress   string   `json:"to_address"`
	Payload     []string `json:"payload"`
}

type Event struct {
	RecordedAt  time.Time `json:"recorded_at"`
	EventId     string    `json:"event_id"`
	FromAddress string    `json:"from_address"`
	Keys        []string  `json:"keys"`
	Data        []string  `json:"data"`
}

type ExecutionResources struct {
	ActualFee              string `json:"actual_fee"`
	BuiltinInstanceCounter struct {
		RangeCheckBuiltin uint `json:"range_check_builtin"`
		PedersenBuiltin   uint `json:"pedersen_builtin"`
		BitwiseBuiltin    uint `json:"bitwise_builtin"`
		OutputBuiltin     uint `json:"output_builtin"`
		EcdsaBuiltin      uint `json:"ecdsa_builtin"`
		EcOpBuiltin       uint `json:"ec_op_builtin"`
	} `json:"builtin_instance_counter"`
	NSteps       uint `json:"n_steps"`
	NMemoryHoles uint `json:"n_memory_holes"`
}

type TransactionReceipt struct {
	ExecutionStatus    string             `json:"execution_status"`
	TransactionHash    string             `json:"transaction_hash"`
	ActualFee          string             `json:"actual_fee"`
	Version            string             `json:"version"`
	L2ToL1Messages     []L2ToL1Message    `json:"l2_to_l1_messages"`
	Events             []Event            `json:"events"`
	ExecutionResources ExecutionResources `json:"execution_resources"`
	TransactionIndex   uint               `json:"transaction_index"`
}

type GetBlockResponse struct {
	BlockHash           string               `json:"block_hash"`
	ParentBlockHash     string               `json:"parent_block_hash"`
	StateRoot           string               `json:"state_root"`
	Status              string               `json:"status"`
	GasPrice            string               `json:"gas_price"`
	SequencerAddress    string               `json:"sequencer_address"`
	StarknetVersion     string               `json:"starknet_version"`
	Transactions        []Transaction        `json:"transactions"`
	TransactionReceipts []TransactionReceipt `json:"transaction_receipts"`
	BlockNumber         uint64               `json:"block_number"`
	Timestamp           uint64               `json:"timestamp"`
}

type SlotUri struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	ExternalUrl string `json:"external_url"`
	YoutubeUrl  string `json:"youtube_url"`
	Attributes  []struct {
		Value       interface{} `json:"value,string"`
		DisplayType string      `json:"display_type"`
		TraitType   string      `json:"trait_type"`
	} `json:"attributes"`
}
