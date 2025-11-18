package encryption

import (
	"testing"

	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/models"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/storage/google"
)

func Test_service_EncryptFile(mainTestCtx *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		encryptRequest models.Encrypt
	}

	tests := []struct {
		name          string
		fields        fields
		args          args
		expectedError string
		wantErr       bool
	}{
		{
			name: "Google storage repository is not set",
			args: args{
				encryptRequest: models.Encrypt{
					KeyFile:               "sandbox_key",
					FileName:              "LMS_Data",
					Location:              "NIFI_Staging",
					EncryptionDestination: "NIFI_Encrypt",
				},
			},
			fields: fields{
				repository: nil,
			},
			expectedError: "google storage/encryption service is not set",
			wantErr:       true,
		},
		{
			name: "Keyfile is an empty string",
			args: args{
				encryptRequest: models.Encrypt{
					KeyFile:               "",
					FileName:              "LMS_Data",
					Location:              "NIFI_Staging",
					EncryptionDestination: "NIFI_Encrypt",
				},
			},
			fields: fields{
				repository: &google.Storage{},
			},
			expectedError: "open : no such file or directory",
			wantErr:       true,
		},
	}
	for _, testCase := range tests {
		mainTestCtx.Run(testCase.name, func(testContext *testing.T) {
			service := service{
				repository: testCase.fields.repository,
			}

			err := service.EncryptFile(testCase.args.encryptRequest)
			if ((err != nil) && testCase.wantErr) && (err.Error() != testCase.expectedError) {
				testContext.Errorf("service.EncryptFile() error = %v, wantErr %v", err, testCase.wantErr)
				testContext.Errorf("service.EncryptFile() expected error: %v", testCase.expectedError)
				return
			}
		})
	}
}
