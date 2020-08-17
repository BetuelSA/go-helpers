package id

import (
	"github.com/BetuelSA/go-helpers/errors"
	"github.com/google/uuid"
)

// ID is the data type that represents an entity's unique ID
type ID = uuid.UUID

// NewID return an ID
func NewID() ID {
	return ID(uuid.New())
}

// StringToID converts a string to an ID, if it's valid
func StringToID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, errors.BadRequest.Wrapf(err, "invalid ID")
	}
	return ID(id), nil
}
