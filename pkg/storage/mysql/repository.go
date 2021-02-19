package mysql

import (
	"fmt"
	"github.com/ONSDigital/blaise-mi-extract/pkg/extractor"
	"github.com/rs/zerolog/log"
	"os"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

type Storage struct {
	DB       sqlbuilder.Database
	Server   string
	Database string
	User     string
	Password string
}

func NewStorage(database string, options ...func(*Storage)) *Storage {
	s := Storage{}

	s.Database = database

	for _, option := range options {
		if option != nil {
			option(&s)
		}
	}

	return &s
}

// connect to the database. Options (database, user etc.) have been set in NewStorage
func (s *Storage) Connect() error {

	var dbURI string
	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")

	if !isSet {
		dbURI = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", s.User, s.Password, s.Server,
			s.Database)
	} else {
		dbURI = fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", s.User, s.Password, socketDir,
			s.Server, s.Database)
	}

	log.Info().Str("Connection string", dbURI).Msg("Connecting to DB")

	settings, err := mysql.ParseURL(dbURI)

	sess, err := mysql.Open(settings)

	if err != nil {
		log.Error().
			Err(err).
			Str("databaseName", s.Database).
			Msg("Cannot connect to database")
		return err
	}

	log.Debug().
		Str("databaseName", s.Database).
		Msg("Connected to database")

	s.DB = sess

	return nil
}

func (s Storage) Close() {
	if s.DB != nil {
		_ = s.DB.Close()
	}
}

func (s Storage) GetMISpecs(instrument string) ([]extractor.MISpec, error) {

	var specs []extractor.MISpec

	q := s.DB.Select("header_name", "response_key").
		From("instrument", "mi_spec", "mi_values").
		Where("instrument.mi_spec = mi_spec.id").
		And("mi_spec.id = mi_values.spec_id").
		And("instrument.name = ?", instrument).
		And("instrument.phase = ?", "live")

	if err := q.All(&specs); err != nil {
		log.Warn().Msgf("no instruments found or no mi specs for %s or database error", instrument)
		return specs, err
	}
	return specs, nil
}

func (s Storage) LoadResponseData(name string) ([]extractor.ResponseData, error) {

	var responses []extractor.ResponseData

	q := s.DB.Select("response_data").
		From("case_response cr", "instrument i", "blaise.case c").
		Where("c.instrument_id = i.id").
		And("cr.case_id = c.id").
		And("i.name = ?", name)

	if err := q.All(&responses); err != nil {
		log.Warn().Msgf("no responses found for %s or database error", name)
		return responses, err
	}

	return responses, nil
}
