package data

import "github.com/pkg/errors"

type Property struct {
	Name string

	// Value may only be nil, or one of the following types:
	// - int64
	// - bool
	// - string
	// - float64
	//
	// If non-nil, it must be exactly one of these types;
	// it is not sufficient to have the same underlying type.
	Value interface{}
}

var ErrNoSuchEntity = errors.New("No Such Entity")
