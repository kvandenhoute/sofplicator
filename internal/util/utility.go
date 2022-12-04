package util

import (
	"fmt"

	"github.com/google/uuid"
)

func GenerateName(name string) string {
	id := uuid.New()
	generatedName := fmt.Sprintf("%s-%s", name, id.String())
	return generatedName[:Min(len(generatedName), 63)]
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
