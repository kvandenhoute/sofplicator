package util

import (
	"fmt"

	"github.com/google/uuid"
)

func GenerateUuid() string {
	return uuid.New().String()
}

func GenerateName(name string, uuid string) string {
	generatedName := fmt.Sprintf("%s-%s", name, uuid)
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
