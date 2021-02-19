package blaise_mi_extract

import (
	"context"
	"github.com/ONSDigital/blaise-mi-extract/pkg"
	"github.com/ONSDigital/blaise-mi-extract/pkg/util"

)

// handles event from item arriving in the encrypt  bucket
func EncryptFunction(ctx context.Context, e util.GCSEvent) error {
	return pkg.HandleEncryptionRequest(ctx, e.Name, e.Bucket)
}
