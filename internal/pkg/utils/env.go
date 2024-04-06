package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func ReadStringEnvKey(envKey string, envIsRequired bool) (string, error) {
	var err error
	value, valueIsExist := os.LookupEnv(envKey)

	if !valueIsExist && envIsRequired {
		err = errors.New(fmt.Sprintf("env with key %s is not found", envKey))
		return "", err
	}

	return value, nil
}

func ReadIntEnvKey(envKey string, envIsRequired bool) (int, error) {
	var err error
	value, valueIsExist := os.LookupEnv(envKey)

	if !valueIsExist && envIsRequired {
		err = errors.New(fmt.Sprintf("env with key %s is not found", envKey))
		return 0, err
	}

	intValue, err := strconv.Atoi(value)

	if err != nil {
		err = errors.New(fmt.Sprintf("failed to convert string to int for %s value", value))
		return 0, err
	}

	return intValue, nil
}
