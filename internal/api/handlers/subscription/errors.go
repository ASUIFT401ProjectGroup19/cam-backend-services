package subscription

type Internal struct{}

func (e Internal) Error() string {
	return "an internal server error occurred"
}
