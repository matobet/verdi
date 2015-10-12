package model

import "github.com/satori/go.uuid"

type GUID string

type Status int8

const (
	Down Status = iota
	Up
	Unknown
)

func NewGUID() GUID {
	return GUID(uuid.NewV4().String())
}
