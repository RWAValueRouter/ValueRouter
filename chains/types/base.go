package types

type ChainType int

const (
	EVM ChainType = iota
)

type Chain interface {
	ValidateParam(string) bool
	ValidateUnsignedTransactionHash(string, string) bool
	GetUnsignedTransactionHash(string) (string, error)
	SendRawTransaction(string, string) (string, error)
}
