package blaise_nifi_encrypt

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	"google.golang.org/api/idtoken"

	"github.com/ONSDigital/blaise-nifi-encrypt/pkg"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/datadeliverystatus"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/models"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/util"
)

func createDataDeliveryStatusClient(client *http.Client) (datadeliverystatus.Client, error) {
	ddsUrl, err := util.EnvironmentVariableExist("DDS_URL")
	if err != nil {
		return datadeliverystatus.Client{}, err
	}

	dataDeliveryStatusClient := datadeliverystatus.Client{
		Config: &datadeliverystatus.Config{
			BaseURL: ddsUrl,
		},
		HTTP: client,
	}

	return dataDeliveryStatusClient, nil
}

// handles event from item arriving in the nifi staging bucket
func NiFiEncryptFunction(ctx context.Context, e models.GCSEvent) error {
	util.ConfigureLogging()

	clientId, err := util.EnvironmentVariableExist("CLIENT_ID")
	if err != nil {
		return err
	}

	client, err := idtoken.NewClient(ctx, clientId)
	if err != nil {
		log.Err(err).Msgf("Could not get IAP token for DDS")
		return err
	}

	dataDeliveryStatusClient, err := createDataDeliveryStatusClient(client)
	if err != nil {
		log.Err(err).Msgf("Trying to create DDS HTTP client failed")
	}

	_, err = dataDeliveryStatusClient.Update(e.Name, "in_staging")
	if err != nil {
		log.Err(err).Msgf("Updating data delivery status to 'in_staging' failed")
	}

	return pkg.HandleEncryptionRequest(ctx, e.Name, e.Bucket, dataDeliveryStatusClient)
}
