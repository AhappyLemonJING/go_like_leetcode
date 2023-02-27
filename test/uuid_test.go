package test

import (
	"fmt"
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestGenerateUUID(t *testing.T) {
	s := uuid.NewV4().String()
	fmt.Printf("s: %v\n", s)
}
