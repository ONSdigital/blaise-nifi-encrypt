package encryption

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/models"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

type Repository interface {
	GetReader(file, directory string) (io.ReadCloser, error)
	GetWriter(file, directory string) io.WriteCloser
}

type Service interface {
	EncryptFile(encryptRequest models.Encrypt) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s service) EncryptFile(encryptRequest models.Encrypt) error {
	storageReader, err := s.r.GetReader(encryptRequest.FileName, encryptRequest.Location)
	if err != nil {
		log.Err(err).Msgf("cannot create a reader")
		return err
	}
	defer storageReader.Close()

	fileName := encryptRequest.FileName
	storageWriter := s.r.GetWriter(fileName, encryptRequest.EncryptionDestination)
	defer storageWriter.Close()

	// Read public key
	recipient, err := readEntity(encryptRequest.KeyFile)
	if err != nil {
		log.Err(err).Msgf("cannot read public key")
		return err
	}
	// Check if public key signatures have expired
	for _, identity := range recipient.Identities {
		if identity.SelfSignature.KeyExpired(time.Now()) {
			err := fmt.Errorf("Key has expired")
			log.Err(err).Msgf("Cannot use public key for '%s'", identity.Name)
			return err
		}
	}

	if err := encrypt([]*openpgp.Entity{recipient}, nil, storageReader, storageWriter); err != nil {
		log.Err(err).Msgf("encrypt failed")
		return err
	}

	log.Info().Msgf("file %s encrypted and saved to %s/%s", encryptRequest.FileName,
		encryptRequest.EncryptionDestination, encryptRequest.FileName)

	return nil
}

func encrypt(recip []*openpgp.Entity, signer *openpgp.Entity, r io.Reader, w io.Writer) error {
	wc, err := openpgp.Encrypt(w, recip, signer, &openpgp.FileHints{IsBinary: true}, nil)
	if err != nil {
        log.Err(err).Msgf("failed to set up encryption: '%s'", err)
		return err
	}

	defer wc.Close()
	if _, err := io.Copy(wc, r); err != nil {
        log.Err(err).Msgf("failed to fetch content and encrypt: '%s'. Updating the Go version could fix tcp connection errors according to https://github.com/googleapis/google-cloud-go/issues/1253 and https://cloud.google.com/functions/docs/concepts/go-runtime", err)
		return err
	}

	return nil
}

func readEntity(name string) (*openpgp.Entity, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	block, err := armor.Decode(f)
	if err != nil {
		return nil, err
	}

	return openpgp.ReadEntity(packet.NewReader(block.Body))
}
