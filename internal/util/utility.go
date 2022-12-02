package util

import (
	"fmt"

	"github.com/google/uuid"
)

func GenerateName(name string) string {
	id := uuid.New()
	return fmt.Sprintf("%s-%s", name, id.String())
}
