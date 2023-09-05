package blaise_nifi_encrypt

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"cloud.google.com/go/functions/metadata"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/datadeliverystatus"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/models"
)

func Test_NiFiEncryptFunction(t *testing.T) {

	name := "LMS_DATA.txt"
	event := models.GCSEvent{
		Name: name,
	}
	meta := &metadata.Metadata{
		EventID: "EVENT_ID",
	}
	tests := []struct {
		name          string
		ctx           context.Context
		e             models.GCSEvent
		expectedError string
		wantErr       bool
	}{
		{
			name:          "CLIENT_ID environment variable is not set up",
			ctx:           metadata.NewContext(context.Background(), meta),
			e:             event,
			expectedError: "the CLIENT_ID environment variable has not been set",
			wantErr:       true,
		},
		{
			name:    "CLIENT_ID environment variable is set with wrong value",
			ctx:     metadata.NewContext(context.Background(), meta),
			e:       event,
			wantErr: true,
		},
		{
			name:          "CLIENT_ID environment variable is an empty string",
			ctx:           metadata.NewContext(context.Background(), meta),
			e:             event,
			expectedError: "the CLIENT_ID environment variable is an empty string",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "CLIENT_ID environment variable is set with wrong value" {
				t.Setenv("CLIENT_ID", "dummy")
			}
			if tt.name == "CLIENT_ID environment variable is an empty string" {
				t.Setenv("CLIENT_ID", "")
			}

			err := NiFiEncryptFunction(tt.ctx, tt.e)
			if ((err != nil) != tt.wantErr) && (err.Error() == tt.expectedError) {
				t.Errorf("NiFiEncryptFunction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func Test_createDataDeliveryStatusClient(t *testing.T) {
	tests := []struct {
		name          string
		client        *http.Client
		want          datadeliverystatus.Client
		wantErr       bool
		expectedError string
	}{
		{
			name:          "DDS_URL environment variable isn't set up",
			client:        &http.Client{},
			want:          datadeliverystatus.Client{},
			wantErr:       true,
			expectedError: "the DDS_URL environment variable has not been set",
		},
		{
			name:          "DDS_URL environment variable is an empty string",
			client:        &http.Client{},
			want:          datadeliverystatus.Client{},
			wantErr:       true,
			expectedError: "the DDS_URL environment variable is an empty string",
		},
		{
			name:   "DDS_URL environment variable is set up",
			client: &http.Client{},
			want: datadeliverystatus.Client{
				Config: &datadeliverystatus.Config{
					BaseURL: "www.dds-blaise.com",
				},
				HTTP: &http.Client{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "DDS_URL environment variable is an empty string" {
				t.Setenv("DDS_URL", "")
			}

			if tt.name == "DDS_URL environment variable is set up" {
				t.Setenv("DDS_URL", "www.dds-blaise.com")
			}

			got, err := createDataDeliveryStatusClient(tt.client)
			if ((err != nil) && tt.wantErr) && (err.Error() != tt.expectedError) {
				t.Errorf("createDataDeliveryStatusClient() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("createDataDeliveryStatusClient() expected error: %v", tt.expectedError)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createDataDeliveryStatusClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
