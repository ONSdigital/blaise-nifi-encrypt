package encryption

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/models"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/openpgp"        //nolint:staticcheck
	"golang.org/x/crypto/openpgp/armor"  //nolint:staticcheck
	"golang.org/x/crypto/openpgp/packet" //nolint:staticcheck
)

type Repository interface {
	GetReader(file, directory string) (io.ReadCloser, error)
	GetWriter(file, directory string) io.WriteCloser
}

type Service interface {
	EncryptFile(encryptRequest models.Encrypt) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (service service) EncryptFile(encryptRequest models.Encrypt) error {
	if service.repository == nil {
		log.Error().Msgf("google storage/encryption service is not set")
		return fmt.Errorf("google storage/encryption service is not set")
	}

	// Read public key
	recipient, err := readEntity(encryptRequest.KeyFile)
	if err != nil {
		log.Err(err).Msgf("Failed to read public key")
		return err
	}
	// Check if public key signatures have expired
	for _, identity := range recipient.Identities {
		if identity.SelfSignature.KeyExpired(time.Now()) {
			err := fmt.Errorf("key has expired")
			log.Err(err).Msgf("Cannot use public key for '%s'", identity.Name)
			return err
		}
	}

	storageReader, err := service.repository.GetReader(encryptRequest.FileName, encryptRequest.Location)
	if err != nil {
		log.Err(err).Msgf("Storage Reader not created for passed file name")
		return err
	}
	defer func() {
		if err := storageReader.Close(); err != nil {
			log.Err(err).Msgf("Failed to close storageReader: %v", err)
		}
	}()

	fileName := encryptRequest.FileName
	storageWriter := service.repository.GetWriter(fileName, encryptRequest.EncryptionDestination)
	defer func() {
		if err := storageWriter.Close(); err != nil {
			log.Err(err).Msgf("Failed to close storageWriter: %v", err)
		}
	}()

	if err := encrypt([]*openpgp.Entity{recipient}, nil, storageReader, storageWriter); err != nil {
		log.Err(err).Msgf("Encrypt failed")
		return err
	}

	log.Info().Msgf("File %s encrypted and saved to %s/%s", encryptRequest.FileName,
		encryptRequest.EncryptionDestination, encryptRequest.FileName)

	return nil
}

func encrypt(recip []*openpgp.Entity, signer *openpgp.Entity, r io.Reader, w io.Writer) error {
	wc, err := openpgp.Encrypt(w, recip, signer, &openpgp.FileHints{IsBinary: true}, nil)
	if err != nil {
		log.Err(err).Msg("Failed to set up encryption")
		return err
	}

	defer func() {
		if err := wc.Close(); err != nil {
			log.Err(err).Msgf("Failed to close wc: %v", err)
		}
	}()
	if _, err := io.Copy(wc, r); err != nil {
		log.Err(err).Msgf("Failed to fetch content and encrypt. Updating the Go version could fix tcp connection errors according to https://github.com/googleapis/google-cloud-go/issues/1253 and https://cloud.google.com/functions/docs/concepts/go-runtime")
		return err
	}

	return nil
}

func readEntity(name string) (*openpgp.Entity, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("failed to close file: %v\n", err)
		}
	}()
	block, err := armor.Decode(f)
	if err != nil {
		return nil, err
	}

	return openpgp.ReadEntity(packet.NewReader(block.Body))
}
