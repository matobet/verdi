package model

import "github.com/satori/go.uuid"

type GUID string

func (id GUID) String() string {
	return string(id)
}

type Status int8

const (
	Down Status = iota
	Up
	Unknown
)

//go:generate stringer -type=Status

var GlobalClusterID = GUID(uuid.Nil.String())

func NewGUID() GUID {
	return GUID(uuid.NewV4().String())
}
