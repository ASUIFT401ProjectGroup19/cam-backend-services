package identity

type AccountExists struct{}

func (e AccountExists) Error() string {
	return "an account with that username already exists"
}

type Internal struct{}

func (e Internal) Error() string {
	return "an internal server error occurred"
}

type LoginFailed struct{}

func (e LoginFailed) Error() string {
	return "unable to log in with provided username-password combination"
}

type Unknown struct{}

func (e Unknown) Error() string {
	return "an unknown error occurred"
}
