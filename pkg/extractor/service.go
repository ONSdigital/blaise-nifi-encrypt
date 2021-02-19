package extractor

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
)

type Service interface {
	ExtractMiInstrument(instrument, destination, destinationFile string) error
}

type FileRepository interface {
	CreateFile(location, destinationFile string) (io.Writer, error)
	DeleteFile(file, directory string) error
	CloseFile()
}

type DBRepository interface {
	GetMISpecs(instrument string) ([]MISpec, error)
	LoadResponseData(name string) ([]ResponseData, error)
	Close()
	Connect() error
}

type ResponseData struct {
	ResponseData string `db:"response_data"`
}

type MISpec struct {
	Header      string `db:"header_name"`
	ResponseKey string `db:"response_key"`
}

type service struct {
	fileRepository FileRepository
	dbRepository   DBRepository
	ctx            context.Context
}

// create a new service instance
func NewService(ctx context.Context, fileRepo FileRepository, dbRepo DBRepository) Service {
	return &service{ctx: ctx, fileRepository: fileRepo, dbRepository: dbRepo}
}

// extract data from the database and save as a csv
func (s service) ExtractMiInstrument(instrument, destination, destinationFile string) error {
	var miSpec []MISpec
	var responseData []ResponseData
	var err error

	if miSpec, err = s.dbRepository.GetMISpecs(instrument); err != nil {
		return err // error already shown
	}

	responseData, err = s.dbRepository.LoadResponseData(instrument)
	if err != nil {
		return err // error already shown
	}

	var c io.Writer
	c, err = s.fileRepository.CreateFile(destination, destinationFile)
	if err != nil {
		log.Err(err).Msgf("cannot create CSV file")
		return err
	}

	csvFile := csv.NewWriter(c)
	// write header line
	keys := make([]string, len(miSpec))
	for i := 0; i < len(miSpec); i++ {
		keys[i] = miSpec[i].Header
	}

	err = csvFile.Write(keys)
	if err != nil {
		log.Err(err).Msgf("cannot write CSV header")
		return err
	}

	// write rows
	for i := 0; i < len(responseData); i++ {
		var js = responseData[i].ResponseData

		m := map[string]string{}
		err = json.Unmarshal([]byte(js), &m)
		if err != nil {
			log.Err(err).Msg("invalid json string in response_data")
			return nil
		}

		var r []string
		for i := 0; i < len(miSpec); i++ {
			var v = miSpec[i].ResponseKey

			if val, ok := m[v]; ok {
				r = append(r, val)
			} else {
				r = append(r, "")
			}
		}

		err = csvFile.Write(r)
		if err != nil {
			log.Err(err).Msgf("cannot write CSV row")
			return err
		}

		r = nil
	}

	csvFile.Flush()
	s.fileRepository.CloseFile()

	return nil
}
