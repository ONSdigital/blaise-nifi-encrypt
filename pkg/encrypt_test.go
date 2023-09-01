package pkg

import (
	"reflect"
	"testing"

	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/models"
)

func Test_loadConfig(t *testing.T) {
	type args struct {
		name     string
		location string
	}
	tests := []struct {
		name          string
		args          args
		want          models.Encrypt
		expectedError string
		wantErr       bool
	}{
		{
			name: "ENCRYPTION_DESTINATION environment variable is not set",
			args: args{
				name:     "test.txt",
				location: "/home/user",
			},
			want:          models.Encrypt{},
			expectedError: "the ENCRYPTION_DESTINATION environment variable has not been set",
			wantErr:       true,
		},
		{
			name: "ENCRYPTION_DESTINATION environment variable is an empty string",
			args: args{
				name:     "test.txt",
				location: "/home/user",
			},
			want:          models.Encrypt{},
			expectedError: "the ENCRYPTION_DESTINATION environment variable is an empty string",
			wantErr:       true,
		},
		{
			name: "PUBLIC_KEY environment variable is not set",
			args: args{
				name:     "test.txt",
				location: "/home/user",
			},
			want:          models.Encrypt{},
			expectedError: "the PUBLIC_KEY environment variable has not been set",
			wantErr:       true,
		},
		{
			name: "PUBLIC_KEY environment variable is an empty string",
			args: args{
				name:     "test.txt",
				location: "/home/user",
			},
			want:          models.Encrypt{},
			expectedError: "the PUBLIC_KEY environment variable is an empty string",
			wantErr:       true,
		},
		{
			name: "ENCRYPTION_DESTINATION and PUBLIC_KEY environment variables are set",
			args: args{
				name:     "test.txt",
				location: "/home/user",
			},
			want: models.Encrypt{
				KeyFile:               "dummy",
				FileName:              "test.txt",
				Location:              "/home/user",
				EncryptionDestination: "/tmp/encrypted",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "ENCRYPTION_DESTINATION environment variable is an empty string" {
				t.Setenv("ENCRYPTION_DESTINATION", "")
			}
			if tt.name == "PUBLIC_KEY environment variable is an empty string" {
				t.Setenv("ENCRYPTION_DESTINATION", "/tmp/encrypted")
				t.Setenv("PUBLIC_KEY", "")
			}
			if tt.name == "ENCRYPTION_DESTINATION and PUBLIC_KEY environment variables are set" {
				t.Setenv("ENCRYPTION_DESTINATION", "/tmp/encrypted")
				t.Setenv("PUBLIC_KEY", "dummy")
			}

			got, err := loadConfig(tt.args.name, tt.args.location)
			if ((err != nil) != tt.wantErr) && (err.Error() == tt.expectedError) {
				t.Errorf("loadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
