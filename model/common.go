package model

import "github.com/satori/go.uuid"

type (
	GUID string

	//go:generate stringer -type=Status
	Status int8
)

const (
	Down Status = iota
	Up
	Unknown
)

var GlobalClusterID = GUID(uuid.Nil.String())

func (id GUID) String() string {
	return string(id)
}

func NewGUID() GUID {
	return GUID(uuid.NewV4().String())
}
