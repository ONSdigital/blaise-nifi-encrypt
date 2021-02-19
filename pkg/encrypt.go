package pkg

import (
	"context"
	"github.com/ONSDigital/blaise-mi-extract/pkg/encryption"
	"github.com/ONSDigital/blaise-mi-extract/pkg/storage/google"
	"github.com/ONSDigital/blaise-mi-extract/pkg/util"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"sync"
)

var encryptedDestination string
var gpg = false
var keyFile string
var encryptOnce sync.Once

func initialiseEncrypt() {
	util.Initialise()

	var found bool

	if encryptedDestination, found = os.LookupEnv(util.EncryptOutput); !found {
		log.Fatal().Msg("The " + util.EncryptOutput + " variable has not been set")
		os.Exit(1)
	}

	log.Info().Msgf("encrypted destination: %s", encryptedDestination)

	if keyFile, found = os.LookupEnv(util.PublicKeyFile); !found {
		log.Fatal().Msg("The " + util.PublicKeyFile + " variable has not been set")
		os.Exit(1)
	}

	log.Info().Msgf("public key file: %s", keyFile)

	useGpg, found := os.LookupEnv(util.UseGPGExtension)
	if !found {
		gpg = false
	} else {
		if strings.EqualFold(useGpg, "TRUE") {
			gpg = true
		}
	}
}

// handles event from item arriving in the encrypt bucket
func HandleEncryptionRequest(ctx context.Context, name, location string) error {

	encryptOnce.Do(func() {
		initialiseEncrypt()
	})

	log.Info().
		Str("location", location).
		Str("file", name).
		Msgf("received encrypt request")

	r := google.NewStorage(ctx)
	encrypt := encryption.NewService(&r)

	encryptRequest := util.Encrypt{
		KeyFile:              keyFile,
		FileName:             name,
		Location:             location,
		EncryptedDestination: encryptedDestination,
		DeleteFile:           true,
		UseGPGExtension:      gpg,
	}

	if err := encrypt.EncryptFile(encryptRequest); err != nil {
		log.Warn().Msg("encrypt failed")
		return err
	}

	return nil
}
