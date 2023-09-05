package encryption

import (
	"testing"

	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/models"
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
				r: nil,
			},
			expectedError: "Google Storage/Encryption Service is not set",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{
				r: tt.fields.r,
			}

			err := s.EncryptFile(tt.args.encryptRequest)
			if ((err != nil) != tt.wantErr) && (err.Error() == tt.expectedError) {
				t.Errorf("service.EncryptFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
