package filter

const (
	// Transaction
	AllTransactionsType = Kind(0) + iota
	AllLogsType
	CoinTransferType
	DeployType
	SpectificDeployType

	// Log
	AddressLogType
	EventLogType
)

type Kind byte

func (k Kind) String() string {
	switch k {
	case AllTransactionsType:
		return "AllTransactionsType"

	case AllLogsType:
		return "AllLogsType"

	case CoinTransferType:
		return "CoinTransferType"

	case DeployType:
		return "DeployType"

	case SpectificDeployType:
		return "SpectificDeployType"

	case AddressLogType:
		return "AddressLogType"

	case EventLogType:
		return "EventLogType"

	default:
		panic("Unknown filter type")
	}
}
