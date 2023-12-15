package share

type ProcessState uint8

const (
	Create ProcessState = iota
	Queued
	Built
	Loaded
	Deployed
	Error
	Deleted
)

func (p ProcessState) String() string {
	switch p {
	case Create:
		return "CREATE"
	case Queued:
		return "QUEUED"
	case Built:
		return "BUILT"
	case Loaded:
		return "LOADED"
	case Deployed:
		return "DEPLOYED"
	case Error:
		return "ERROR"
	case Deleted:
		return "DELETED"
	}

	return "UNKNOWN"
}
