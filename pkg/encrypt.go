package pkg

import (
	"context"

	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/datadeliverystatus"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/encryption"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/models"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/storage/google"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/util"
	"github.com/rs/zerolog/log"
)

func loadConfig(name string, location string) (models.Encrypt, error) {
	encryptionDestination, err := util.EnvironmentVariableExist("ENCRYPTION_DESTINATION")
	if err != nil {
		return models.Encrypt{}, err
	}

	log.Info().Msgf("encrypted destination: %s", encryptionDestination)

	keyFile, err := util.EnvironmentVariableExist("PUBLIC_KEY")
	if err != nil {
		return models.Encrypt{}, err
	}

	log.Info().Msgf("public key file: %s", keyFile)

	return models.Encrypt{
		KeyFile:               keyFile,
		FileName:              name,
		Location:              location,
		EncryptionDestination: encryptionDestination,
	}, nil
}

// handles event from item arriving in the encrypt bucket
func HandleEncryptionRequest(ctx context.Context, name, location string, dataDeliveryStatusClient datadeliverystatus.Client) error {
	log.Info().
		Str("location", location).
		Str("file", name).
		Msgf("received encrypt request")

	encryptRequest, err := loadConfig(name, location)
	if err != nil {
		log.Err(err).Msgf("Creating encrypt request failed")
		return err
	}

	r, err := google.NewStorage(ctx)
	if err != nil {
		log.Err(err).Msg("Could not create GCP storage client")
		return err
	}
	encrypt := encryption.NewService(&r)

	if err := encrypt.EncryptFile(encryptRequest); err != nil {
		_, ddsErr := dataDeliveryStatusClient.Error("errored", name, err.Error())
		if ddsErr != nil {
			log.Err(ddsErr).Msgf("Updating data delivery status to 'errored' failed")
		}
		return err
	}
	_, ddsErr := dataDeliveryStatusClient.Update("encrypted", name)
	if ddsErr != nil {
		log.Err(ddsErr).Msgf("Updating data delivery status to 'encrypted' failed")
	}
	return nil
}
