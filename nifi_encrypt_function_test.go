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

func Test_NiFiEncryptFunction(mainTestCtx *testing.T) {

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
	for _, testCase := range tests {
		mainTestCtx.Run(testCase.name, func(testContext *testing.T) {
			if testCase.name == "CLIENT_ID environment variable is set with wrong value" {
				testContext.Setenv("CLIENT_ID", "dummy")
			}
			if testCase.name == "CLIENT_ID environment variable is an empty string" {
				testContext.Setenv("CLIENT_ID", "")
			}

			err := NiFiEncryptFunction(testCase.ctx, testCase.e)
			if ((err != nil) != testCase.wantErr) && (err.Error() == testCase.expectedError) {
				testContext.Errorf("NiFiEncryptFunction() error = %v, wantErr %v", err, testCase.wantErr)
			}
		})
	}

}

func Test_createDataDeliveryStatusClient(mainTestCtx *testing.T) {
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
	for _, testCase := range tests {
		mainTestCtx.Run(testCase.name, func(testContext *testing.T) {
			if testCase.name == "DDS_URL environment variable is an empty string" {
				testContext.Setenv("DDS_URL", "")
			}

			if testCase.name == "DDS_URL environment variable is set up" {
				testContext.Setenv("DDS_URL", "www.dds-blaise.com")
			}

			got, err := createDataDeliveryStatusClient(testCase.client)
			if ((err != nil) && testCase.wantErr) && (err.Error() != testCase.expectedError) {
				testContext.Errorf("createDataDeliveryStatusClient() error = %v, wantErr %v", err, testCase.wantErr)
				testContext.Errorf("createDataDeliveryStatusClient() expected error: %v", testCase.expectedError)
				return
			}
			if !reflect.DeepEqual(got, testCase.want) {
				testContext.Errorf("createDataDeliveryStatusClient() = %v, want %v", got, testCase.want)
			}
		})
	}
}
