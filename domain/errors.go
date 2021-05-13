package domain

type ErrorType uint64

type GinError struct {
	Error string `json:"error" bson:"error"`
}
