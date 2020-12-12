package models

type Error struct {
	Message string `json:"message"`
}

type ServError struct {
	Code    string
	Message string
}

func (e ServError) Error() string {
	return e.Code
}