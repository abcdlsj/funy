package dyloader

type Identity uint8

const (
	FnnyWebService Identity = iota
	FnnyInvokeOnce
)

func (i Identity) String() string {
	switch i {
	case FnnyWebService:
		return "FnnyWebService"
	case FnnyInvokeOnce:
		return "FnnyInvokeOnce"
	}

	return "UNKNOWN"
}
