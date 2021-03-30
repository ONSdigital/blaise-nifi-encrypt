package blaise_nifi_encrypt

import (
	"context"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/ONSDigital/blaise-nifi-encrypt/pkg"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/datadeliverystatus"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/models"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/util"
)

// handles event from item arriving in the nifi staging bucket
func NiFiEncryptFunction(ctx context.Context, e models.GCSEvent) error {
	util.ConfigureLogging()
	dataDeliveryStatusClient := datadeliverystatus.Client{
		Config: &datadeliverystatus.Config{
			BaseURL: os.Getenv("DDS_URL"),
		},
		HTTP: &http.Client{},
	}
	_, err := dataDeliveryStatusClient.Update(e.Name, "in_staging")
	if err != nil {
		log.Error().Msgf("Updating data delivery status to 'in_staging' failed: %s", err.Error())
	}
	return pkg.HandleEncryptionRequest(ctx, e.Name, e.Bucket, dataDeliveryStatusClient)
}
