package filter

const (
	// Transaction
	AllTransactionsFilter = Kind(0) + iota
	AllLogsFilter
	CoinTransferFilter
	DeployFilter
	SpectificDeployFilter

	// Log
	AddressLogFilter
	EventLogFilter
)

type Kind byte

func (k Kind) String() string {
	switch k {
	case AllTransactionsFilter:
		return "AllTransactionsFilter"

	case AllLogsFilter:
		return "AllLogsFilter"

	case CoinTransferFilter:
		return "CoinTransferFilter"

	case DeployFilter:
		return "DeployFilter"

	case SpectificDeployFilter:
		return "SpectificDeployFilter"

	case AddressLogFilter:
		return "AddressLogFilter"

	case EventLogFilter:
		return "EventLogFilter"

	default:
		panic("Unknown filter")
	}
}
