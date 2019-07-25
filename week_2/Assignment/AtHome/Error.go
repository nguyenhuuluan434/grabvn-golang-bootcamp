package main

type customError struct {
	message string
}

func (c customError) Error() string {
	return c.message
}
