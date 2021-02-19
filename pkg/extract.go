package pkg

import (
	"context"
	"encoding/json"
	"github.com/ONSDigital/blaise-mi-extract/pkg/extractor"
	"github.com/ONSDigital/blaise-mi-extract/pkg/storage/google"
	"github.com/ONSDigital/blaise-mi-extract/pkg/storage/mysql"
	"github.com/ONSDigital/blaise-mi-extract/pkg/util"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
)

var encryptDestination string
var extOnce sync.Once
var db extractor.DBRepository

// Sets up the database connection options and connects.
// Lazy call to avoid issues with multiple init() functions
func initialiseExtract() {

	util.Initialise()

	var found bool

	if encryptDestination, found = os.LookupEnv(util.ExtractOutput); !found {
		log.Fatal().Msg("The " + util.ExtractOutput + " variable has not been set")
		os.Exit(1)
	}

	database, found := os.LookupEnv(util.Database)
	if !found {
		database = util.DefaultDatabase
	}

	server := func(db *mysql.Storage) {
		var s string
		if s, found = os.LookupEnv(util.Server); !found {
			log.Fatal().Msg("The " + util.Server + " varible has not been set")
			os.Exit(1)
		}
		db.Server = s
	}

	user := func(db *mysql.Storage) {
		var user string
		if user, found = os.LookupEnv(util.User); !found {
			log.Fatal().Msg("The " + util.Server + " varible has not been set")
			os.Exit(1)
		}
		db.User = user
	}

	password := func(db *mysql.Storage) {
		var pwd string
		if pwd, found = os.LookupEnv(util.Password); !found {
			log.Fatal().Msg("The " + util.Password + " varible has not been set")
			os.Exit(1)
		}
		db.Password = pwd
	}

	db = mysql.NewStorage(database, server, user, password)

	if err := db.Connect(); err != nil {
		// errors have already been reported and we can't continue so stop
		os.Exit(1)
	}

}

type PayloadMessage struct {
	Action     string `json:"action"`
	Instrument string `json:"instrument_name"`
}

type PubSubMessage struct {
	Data []byte `json:"data"`
}

// handle extract request events from publish / subscribe  queue
func HandleExtractionRequest(ctx context.Context, m PubSubMessage) error {

	extOnce.Do(func() { initialiseExtract() })

	gcloudStorage := google.NewStorage(ctx)
	service := extractor.NewService(ctx, &gcloudStorage, db)

	payload := PayloadMessage{}

	if err := json.Unmarshal(m.Data, &payload); err != nil {
		log.Warn().Msgf("Cannot decode json message")
		return nil
	}

	log.Warn().Msgf("message  -> [%s]", m)

	// add additional actions as needed
	switch payload.Action {
	case "extract_mi":
		return extractMI(service, payload.Instrument)
	default:
		log.Warn().Msgf("message rejected, unknown action -> [%s]", payload.Action)
		return nil
	}
}

func extractMI(service extractor.Service, instrument string) error {
	log.Info().Msgf("received extract_mi request for %s", instrument)

	var err error

	destination := instrument + ".csv"
	if err = service.ExtractMiInstrument(instrument, encryptDestination, destination); err != nil {
		return err
	}

	log.Info().Msgf("extract_mi request for %s completed", instrument)

	return nil
}
