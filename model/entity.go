package model

type (
	Entity interface {
		Validate() error
		SetID(int)
	}
)
