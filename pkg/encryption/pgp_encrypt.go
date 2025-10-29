package encryption

import (
	"bufio"
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
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (service service) EncryptFile(encryptRequest models.Encrypt) error {
	if service.repository == nil {
		log.Error().Msgf("Google Storage/Encryption Service is not set")
		return fmt.Errorf("Google Storage/Encryption Service is not set")
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
	defer storageReader.Close()

	fileName := encryptRequest.FileName
	storageWriter := service.repository.GetWriter(fileName, encryptRequest.EncryptionDestination)
	defer storageWriter.Close()

	// Add buffering for better performance with large files
	// 1MB buffer size for efficient I/O
	bufferedReader := bufio.NewReaderSize(storageReader, 1024*1024)
	bufferedWriter := bufio.NewWriterSize(storageWriter, 1024*1024)
	defer bufferedWriter.Flush()

	start := time.Now()
	if err := encrypt([]*openpgp.Entity{recipient}, nil, bufferedReader, bufferedWriter); err != nil {
		log.Err(err).Msgf("Encrypt failed")
		return err
	}
	duration := time.Since(start)
	log.Info().
		Dur("duration", duration).
		Float64("duration_minutes", duration.Minutes()).
		Msgf("Encryption completed for file %s in %.2f minutes", fileName, duration.Minutes())

	return nil
}

func encrypt(recip []*openpgp.Entity, signer *openpgp.Entity, r io.Reader, w io.Writer) error {
	// Configure packet config for better memory efficiency with large files
	config := &packet.Config{
		DefaultCipher: packet.CipherAES256,
		// Use stronger compression for better network efficiency
		DefaultCompressionAlgo: packet.CompressionZIP,
		CompressionConfig: &packet.CompressionConfig{
			Level: 6, // Balanced compression level
		},
	}
	
	wc, err := openpgp.Encrypt(w, recip, signer, &openpgp.FileHints{IsBinary: true}, config)
	if err != nil {
		log.Err(err).Msg("Failed to set up encryption")
		return err
	}

	defer wc.Close()
	log.Info().Msg("Starting to copy and encrypt file content")
	
	// Track progress for large files
	pr := &progressReader{
		reader:       r,
		bytesRead:    0,
		lastLogBytes: 0,
		logInterval:  100 * 1024 * 1024, // Log every 100MB
	}
	
	if _, err := io.Copy(wc, pr); err != nil {
		log.Err(err).Str("error_details", err.Error()).Msgf("Failed to fetch content and encrypt. Updating the Go version could fix tcp connection errors according to https://github.com/googleapis/google-cloud-go/issues/1253 and https://cloud.google.com/functions/docs/concepts/go-runtime")
		return err
	}
	
	log.Info().Int64("total_bytes", pr.bytesRead).Msgf("Total bytes read and encrypted: %d MB", pr.bytesRead/(1024*1024))

	return nil
}

// progressReader wraps an io.Reader to track and log progress
type progressReader struct {
	reader       io.Reader
	bytesRead    int64
	lastLogBytes int64
	logInterval  int64
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.bytesRead += int64(n)
	
	// Log progress at intervals
	if pr.bytesRead-pr.lastLogBytes >= pr.logInterval {
		log.Info().Int64("bytes_processed", pr.bytesRead).Msgf("Encryption progress: %d MB processed", pr.bytesRead/(1024*1024))
		pr.lastLogBytes = pr.bytesRead
	}
	
	return n, err
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
