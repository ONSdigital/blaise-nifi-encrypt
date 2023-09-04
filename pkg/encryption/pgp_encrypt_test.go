package encryption

import (
	"context"
	"testing"

	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/models"
	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/storage/google"
)

func Test_service_EncryptFile(t *testing.T) {
	type fields struct {
		r Repository
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
		// TODO: Add test cases.
		{
			name: "Google storage repository is not set",
			args: args{
				encryptRequest: models.Encrypt{
					KeyFile:               "sandbox_key",
					FileName:              "LMS_Data",
					Location:              "NIFI_Stanging",
					EncryptionDestination: "NIFI_Encrypt",
				},
			},
			fields: fields{
				r: nil,
			},
			expectedError: "Google Storage/Encryption Service is not set",
			wantErr:       true,
		},
		{
			name: "Encryption key is not set in argument object",
			args: args{
				encryptRequest: models.Encrypt{
					KeyFile:               "",
					FileName:              "LMS_Data",
					Location:              "NIFI_Stanging",
					EncryptionDestination: "NIFI_Encrypt",
				},
			},
			fields: fields{
				r: google.NewStorage(context.Background()),
			},
			expectedError: "Encryption/Public key problem",
			wantErr:       true,
		},
		{
			name: "Storage reader is not created due to wrong file name",
			args: args{
				encryptRequest: models.Encrypt{
					KeyFile:               "./keys/dev-key.gpg",
					FileName:              "fake_file_name",
					Location:              "NIFI_Stanging",
					EncryptionDestination: "NIFI_Encrypt",
				},
			},
			fields: fields{
				r: google.NewStorage(context.Background()),
			},
			expectedError: "Storage Reader not created for passed file name",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				r: tt.fields.r,
			}
			if err := s.EncryptFile(tt.args.encryptRequest); (err != nil) != tt.wantErr {
				t.Errorf("service.EncryptFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
