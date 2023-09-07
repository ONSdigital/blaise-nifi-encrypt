package pkg

import (
	"reflect"
	"testing"

	"github.com/ONSDigital/blaise-nifi-encrypt/pkg/models"
)

func Test_loadConfig(mainTestCtx *testing.T) {
	type args struct {
		name     string
		location string
	}

	testArgs := args{
		name:     "test.txt",
		location: "/home/user",
	}
	tests := []struct {
		name          string
		args          args
		want          models.Encrypt
		expectedError string
		wantErr       bool
	}{
		{
			name:          "ENCRYPTION_DESTINATION environment variable is not set",
			args:          testArgs,
			want:          models.Encrypt{},
			expectedError: "the ENCRYPTION_DESTINATION environment variable has not been set",
			wantErr:       true,
		},
		{
			name:          "ENCRYPTION_DESTINATION environment variable is an empty string",
			args:          testArgs,
			want:          models.Encrypt{},
			expectedError: "the ENCRYPTION_DESTINATION environment variable is an empty string",
			wantErr:       true,
		},
		{
			name:          "PUBLIC_KEY environment variable is not set",
			args:          testArgs,
			want:          models.Encrypt{},
			expectedError: "the PUBLIC_KEY environment variable has not been set",
			wantErr:       true,
		},
		{
			name:          "PUBLIC_KEY environment variable is an empty string",
			args:          testArgs,
			want:          models.Encrypt{},
			expectedError: "the PUBLIC_KEY environment variable is an empty string",
			wantErr:       true,
		},
		{
			name: "ENCRYPTION_DESTINATION and PUBLIC_KEY environment variables are set",
			args: testArgs,
			want: models.Encrypt{
				KeyFile:               "dummy",
				FileName:              "test.txt",
				Location:              "/home/user",
				EncryptionDestination: "/tmp/encrypted",
			},
			wantErr: false,
		},
	}
	for _, testCase := range tests {
		mainTestCtx.Run(testCase.name, func(testContext *testing.T) {
			if testCase.name == "ENCRYPTION_DESTINATION environment variable is an empty string" {
				testContext.Setenv("ENCRYPTION_DESTINATION", "")
			}
			if testCase.name == "PUBLIC_KEY environment variable is not set" {
				testContext.Setenv("ENCRYPTION_DESTINATION", "/tmp/encrypted")
			}
			if testCase.name == "PUBLIC_KEY environment variable is an empty string" {
				testContext.Setenv("ENCRYPTION_DESTINATION", "/tmp/encrypted")
				testContext.Setenv("PUBLIC_KEY", "")
			}
			if testCase.name == "ENCRYPTION_DESTINATION and PUBLIC_KEY environment variables are set" {
				testContext.Setenv("ENCRYPTION_DESTINATION", "/tmp/encrypted")
				testContext.Setenv("PUBLIC_KEY", "dummy")
			}

			got, err := loadConfig(testCase.args.name, testCase.args.location)
			if ((err != nil) && testCase.wantErr) && (err.Error() != testCase.expectedError) {
				testContext.Errorf("loadConfig() error = %v, wantErr %v", err, testCase.wantErr)
				testContext.Errorf("loadConfig() expected error: %v", testCase.expectedError)
				return
			}
			if !reflect.DeepEqual(got, testCase.want) {
				testContext.Errorf("loadConfig() = %v, want %v", got, testCase.want)
			}
		})
	}
}
