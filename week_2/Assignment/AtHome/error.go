package main

type UnknownError struct {
	message string
}

func (c UnknownError) Error() string {
	return c.message
}

type ErrFileNotFound struct {

}

func (ErrFileNotFound) Error() string {

}

type ErrFileAccess struct{

}