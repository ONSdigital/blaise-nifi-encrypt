package blaise_mi_extract

import (
	"context"
	"github.com/ONSDigital/blaise-mi-extract/pkg"
	"github.com/ONSDigital/blaise-mi-extract/pkg/util"
)

func ZipFunction(ctx context.Context, e util.GCSEvent) error {
	return pkg.ZipCompress(ctx, e.Name, e.Bucket)
}
