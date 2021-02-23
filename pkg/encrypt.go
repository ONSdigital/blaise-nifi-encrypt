package pkg

import (
	"context"
	"os"

	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/encryption"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/models"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/storage/google"
	"github.com/rs/zerolog/log"
)

func loadConfig(name, location string) models.Encrypt {
	var (
		encryptionDestination string
		keyFile               string
		found                 bool
	)

	if encryptionDestination, found = os.LookupEnv("ENCRYPTION_DESTINATION"); !found {
		log.Fatal().Msg("The ENCRYPTION_DESTINATION variable has not been set")
	}

	log.Info().Msgf("encrypted destination: %s", encryptionDestination)

	if keyFile, found = os.LookupEnv("PUBLIC_KEY"); !found {
		log.Fatal().Msg("The PUBLIC_KEY variable has not been set")
	}

	log.Info().Msgf("public key file: %s", keyFile)

	return models.Encrypt{
		KeyFile:               keyFile,
		FileName:              name,
		Location:              location,
		EncryptionDestination: encryptionDestination,
	}
}

// handles event from item arriving in the encrypt bucket
func HandleEncryptionRequest(ctx context.Context, name, location string) error {
	log.Info().
		Str("location", location).
		Str("file", name).
		Msgf("received encrypt request")

	encryptRequest := loadConfig(name, location)

	r := google.NewStorage(ctx)
	encrypt := encryption.NewService(&r)

	if err := encrypt.EncryptFile(encryptRequest); err != nil {
		log.Warn().Msg("encrypt failed")
		return err
	}

	return nil
}
