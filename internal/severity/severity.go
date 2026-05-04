package severity

type Severity string

const (
	None     Severity = "none"
	Info     Severity = "info"
	Warning  Severity = "warning"
	High     Severity = "high"
	Critical Severity = "critical"
)

func Rank(s Severity) int {
	switch s {
	case None:
		return -1
	case Info:
		return 0
	case Warning:
		return 1
	case High:
		return 2
	case Critical:
		return 3
	default:
		return 0
	}
}

func GTE(actual Severity, threshold Severity) bool {
	return Rank(actual) >= Rank(threshold)
}

func FromScore(score int) Severity {
	switch {
	case score >= 9:
		return Critical
	case score >= 7:
		return High
	case score >= 4:
		return Warning
	default:
		return Info
	}
}
