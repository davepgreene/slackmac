package utils

import (
	"github.com/satori/go.uuid"
)

// CreateCorrelationID generates a new uuid
func CreateCorrelationID() string {
	u, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return u.String()
}
