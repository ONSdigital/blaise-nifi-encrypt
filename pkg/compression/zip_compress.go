package compression

import (
	"archive/zip"
	"github.com/ONSDigital/blaise-mi-extract/pkg/util"
	"github.com/rs/zerolog/log"
	"io"
	"path/filepath"
	"strings"
	"time"
)

type Repository interface {
	CreateFile(location, destinationFile string) (io.Writer, error)
	DeleteFile(file, directory string) error
	GetReader(file, directory string) (io.ReadCloser, error)
	GetWriter(file, directory string) io.WriteCloser
}

type Service interface {
	ZipFile(c util.Zip) (string, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s service) ZipFile(c util.Zip) (string, error) {

	currentTime := time.Now()
	t := strings.TrimSuffix(c.FileName, filepath.Ext(c.FileName)) // strip off .gpg suffix
	u := strings.TrimSuffix(t, filepath.Ext(t))                   // strip off .csv suffix
	zipName := "mi_" + u + "_" + currentTime.Format("02012006") + "_" + currentTime.Format("150405") + ".zip"

	storageReader, err := s.r.GetReader(c.FileName, c.FromLocation)
	if err != nil {
		log.Err(err).Msgf("cannot create a reader")
		return "", err
	}
	defer func() { _ = storageReader.Close() }()

	storageWriter := s.r.GetWriter(zipName, c.ToLocation)
	defer func() { _ = storageWriter.Close() }()

	zipWriter := zip.NewWriter(storageWriter)
	defer func() { _ = zipWriter.Close() }()

	// add filename to compress
	zipFile, err := zipWriter.Create(c.FileName)
	if err != nil {
		log.Err(err).Msgf("error adding file to compress: %s in directory %s", zipName+".compress", c.ToLocation)
		return "", err
	}

	_, err = io.Copy(zipFile, storageReader)

	if err != nil {
		log.Err(err).Msgf("error creating compress file: %s in directory %s", c.FileName+".compress", c.ToLocation)
		return "", err
	}

	log.Debug().Msgf("saved %s/%s", c.ToLocation, c.FileName)

	if c.DeleteFile {
		if err := s.r.DeleteFile(c.FileName, c.FromLocation); err != nil {
			return "", err
		}
	}

	return zipName, nil
}
