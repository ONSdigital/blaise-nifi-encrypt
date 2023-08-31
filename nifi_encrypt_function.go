package blaise_nifi_encrypt

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	"google.golang.org/api/idtoken"

	"github.com/ONSDigital/blaise-nifi-encrypt/pkg"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/datadeliverystatus"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/models"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/util"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func createDataDeliveryStatusClient(client *http.Client) (datadeliverystatus.Client, error) {
	var (
		ddsUrl string
		found  bool
	)

	if ddsUrl, found = os.LookupEnv("DDS_URL"); !found {
		log.Fatal().Msg("The DDS_URL variable has not been set")
		return datadeliverystatus.Client{}, (fmt.Errorf("the DDS_URL variable has not been set"))
	}

	dataDeliveryStatusClient := datadeliverystatus.Client{
		Config: &datadeliverystatus.Config{
			BaseURL: ddsUrl,
		},
		HTTP: client,
	}

	return dataDeliveryStatusClient, fmt.Errorf("")
}

// handles event from item arriving in the nifi staging bucket
func NiFiEncryptFunction(ctx context.Context, e models.GCSEvent) error {
	client, err := idtoken.NewClient(ctx, os.Getenv("CLIENT_ID"))
	if err != nil {
		log.Error().Msgf("Could not get IAP token for DDS: %s", err.Error())
		return err
	}
	util.ConfigureLogging()

	dataDeliveryStatusClient, err := createDataDeliveryStatusClient(client)
	if err.Error() != "" {
		log.Error().Msgf(err.Error())
	}

	_, err = dataDeliveryStatusClient.Update(e.Name, "in_staging")
	if err != nil {
		log.Error().Msgf("Updating data delivery status to 'in_staging' failed: %s", err.Error())
	}
	return pkg.HandleEncryptionRequest(ctx, e.Name, e.Bucket, dataDeliveryStatusClient)
}
