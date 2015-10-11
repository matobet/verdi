package model

type GUID string

type Status int8

const (
	Up Status = iota
	Down
	Unknown
)
