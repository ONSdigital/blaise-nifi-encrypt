package util

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

func EnvironmentVariableExist(name string) (string, error) {
	value, found := os.LookupEnv(name)

	if !found {
		log.Error().Msgf("The %v environment variable has not been set", name)
		return value, (fmt.Errorf("the %v environment variable has not been set", name))
	}

	if value == "" {
		log.Error().Msgf("The %v environment variable is an empty string", name)
		return value, (fmt.Errorf("the %v environment variable is an empty string", name))
	}

	return value, nil
}
