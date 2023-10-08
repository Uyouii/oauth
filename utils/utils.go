package utils

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func GetHexUuid() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	byteData, err := uuid.MarshalBinary()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(byteData), nil
}

func GetUuid() string {
	uuid, err := GetHexUuid()
	if err != nil {
		return fmt.Sprintf("%x%x", rand.Int63(), time.Now().UnixNano())
	}
	return uuid
}

func GetRandomSecret() string {
	return GetUuid() + GetUuid()
}
