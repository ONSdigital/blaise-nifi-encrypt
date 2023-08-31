package blaise_nifi_encrypt

import (
	"context"
	"reflect"
	"testing"

	"cloud.google.com/go/functions/metadata"
	http_mocks "github.com/ONSDigital/blaise-nifi-encrypt/nifi_encrypt_function_mocks"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/datadeliverystatus"
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

func Test_createDataDeliveryStatusClient(t *testing.T) {
	tests := []struct {
		name    string
		client  *HTTPClient
		want    datadeliverystatus.Client
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "DDS URL is not set up",
			client:  &http_mocks.NewHTTPClient(),
			want:    datadeliverystatus.Client{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createDataDeliveryStatusClient(tt.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("createDataDeliveryStatusClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createDataDeliveryStatusClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func NewHTTPClient() {
	panic("unimplemented")
}
