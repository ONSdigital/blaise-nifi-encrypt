package blaise_nifi_encrypt

import (
	"context"
	"testing"

	"cloud.google.com/go/functions/metadata"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/models"
)

func TestNiFiEncryptFunction(t *testing.T) {

	name := "LMS_DATA.txt"
	event := models.GCSEvent{
		Name: name,
	}
	meta := &metadata.Metadata{
		EventID: "EVENT_ID",
	}

	tests := []struct {
		name    string
		ctx     context.Context
		e       models.GCSEvent
		wantErr bool
	}{
		{
			name:    "Client id environment variable is not set up",
			ctx:     metadata.NewContext(context.Background(), meta),
			e:       event,
			wantErr: true,
		},
		{
			name:    "Client id environment variable is set with wrong value",
			ctx:     metadata.NewContext(context.Background(), meta),
			e:       event,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Client id environment variable is set with wrong value" {
				{
					t.Setenv("CLIENT_ID", "dummy")
				}
			}
			err := NiFiEncryptFunction(tt.ctx, tt.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("NiFiEncryptFunction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}
