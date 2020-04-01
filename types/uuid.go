package types

import uuid "github.com/satori/go.uuid"

type UUID string

func NewUUID() UUID {
	return UUID(uuid.NewV4().String())
}

func (id UUID) Valid() bool {
	return string(id) != ""
}

func (id UUID) Equal(o UUID) bool {
	return id == o
}

func (id UUID) String() string {
	return string(id)
}
