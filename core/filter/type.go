package filter

type Kind byte

const (
	// Transaction
	CoinTransferFilter = Kind(0) + iota
	DeployFilter
	SpectificDeployFilter
	AllTransactionsFilter
	AllLogsFilter

	// Log
	AddressLogsFilter
	EventLogsFilter
)

func (k Kind) String() string {
	switch k {
	case CoinTransferFilter:
		return "CoinTransferFilter"

	case DeployFilter:
		return "DeployFilter"

	case SpectificDeployFilter:
		return "SpectificDeployFilter"

	case AllTransactionsFilter:
		return "AllTransactionsFilter"

	case AllLogsFilter:
		return "AllLogsFilter"

	case AddressLogsFilter:
		return "AddressLogsFilter"

	case EventLogsFilter:
		return "EventLogsFilter"

	default:
		panic("Unknown filter")
	}
}
