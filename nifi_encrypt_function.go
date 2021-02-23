package blaise_nifi_encrypt

import (
	"context"

	"github.com/ONSDigital/blaise-nifi-encrypt/pkg"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/models"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/util"
)

// handles event from item arriving in the nifi staging bucket
func NiFiEncryptFunction(ctx context.Context, e models.GCSEvent) error {
	util.ConfigureLogging()
	return pkg.HandleEncryptionRequest(ctx, e.Name, e.Bucket)
}
