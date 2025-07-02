package utils

import "github.com/google/uuid"

func StringToUuid(val string) uuid.UUID {
	id, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil
	}
	return id
}
