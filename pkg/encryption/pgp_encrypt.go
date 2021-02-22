package encryption

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ONSDigital/blaise-mi-extract/pkg/util"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

type Repository interface {
	DeleteFile(file, directory string) error
	GetReader(file, directory string) (io.ReadCloser, error)
	GetWriter(file, directory string) io.WriteCloser
}

type Service interface {
	EncryptFile(encryptRequest util.Encrypt) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s service) DeleteFile(file, directory string) error {
	return s.r.DeleteFile(file, directory)
}

func (s service) EncryptFile(encryptRequest util.Encrypt) error {

	storageReader, err := s.r.GetReader(encryptRequest.FileName, encryptRequest.Location)
	if err != nil {
		log.Err(err).Msgf("cannot create a reader")
		return err
	}
	defer func() { _ = storageReader.Close() }()

	fileName := encryptRequest.FileName
	if encryptRequest.UseGPGExtension {
		fileName = fileName + ".gpg"
	}
	storageWriter := s.r.GetWriter(fileName, encryptRequest.EncryptedDestination)
	defer func() { _ = storageWriter.Close() }()

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
		encryptRequest.EncryptedDestination, encryptRequest.FileName+".gpg")

	if encryptRequest.DeleteFile {
		if err := s.r.DeleteFile(encryptRequest.FileName, encryptRequest.Location); err != nil {
			return err
		}
	}

	return nil
}

func encrypt(recip []*openpgp.Entity, signer *openpgp.Entity, r io.Reader, w io.Writer) error {
	wc, err := openpgp.Encrypt(w, recip, signer, &openpgp.FileHints{IsBinary: true}, nil)
	if err != nil {
		return err
	}

	defer func() { _ = wc.Close() }()
	if _, err := io.Copy(wc, r); err != nil {
		return err
	}

	return nil
}

func readEntity(name string) (*openpgp.Entity, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	block, err := armor.Decode(f)
	if err != nil {
		return nil, err
	}

	return openpgp.ReadEntity(packet.NewReader(block.Body))
}
