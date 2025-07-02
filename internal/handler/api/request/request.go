package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type (
	GetByIdReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
)

func (c GetByIdReq) Validate() error {
	err := validation.ValidateStruct(
		&c,
		validation.Field(&c.Username, validation.Required),
		validation.Field(&c.Password, validation.Required),
	)

	return err
}
