package datadeliverystatus

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type fakeDDSService func(*http.Request) (*http.Response, error)

func (s fakeDDSService) RoundTrip(req *http.Request) (*http.Response, error) {
	return s(req)
}

func createCustomRoundTripper(expectedValue string) http.RoundTripper {
	return fakeDDSService(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: io.NopCloser(strings.NewReader(expectedValue)),
		}, nil
	})
}

func TestClient_Update(mainTestCtx *testing.T) {
	config := &Config{
		BaseURL: "www.dds-blaise.com",
	}

	type args struct {
		state    string
		fileName string
	}
	tests := []struct {
		name    string
		client  Client
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Set 'state' to 'in_staging'",
			client: Client{
				Config: config,
				HTTP: &http.Client{
					Transport: createCustomRoundTripper(`{"state":"in_staging"}`),
				},
			},
			args: args{
				state:    "in_staging",
				fileName: "test.txt",
			},
			want:    `{"state":"in_staging"}`,
			wantErr: false,
		},
	}
	for _, testCase := range tests {
		mainTestCtx.Run(testCase.name, func(testContext *testing.T) {
			client := testCase.client

			got, err := client.Update(testCase.args.state, testCase.args.fileName)

			if (err != nil) != testCase.wantErr {
				testContext.Errorf("Client.Update() error = %v, wantErr %v", err, testCase.wantErr)
				return
			}
			if got != testCase.want {
				testContext.Errorf("Client.Update() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestClient_Error(mainTestCtx *testing.T) {
	config := &Config{
		BaseURL: "www.dds-blaise.com",
	}

	type args struct {
		state        string
		fileName     string
		errorMessage string
	}
	tests := []struct {
		name    string
		client  Client
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Set 'state' to 'errored' and 'error_info' to 'Oops, something went wrong'",
			client: Client{
				Config: config,
				HTTP: &http.Client{
					Transport: createCustomRoundTripper(`{"state":"errored", "error_info": "Oops, something went wrong"}`),
				},
			},
			args: args{
				state:        "errored",
				fileName:     "test.txt",
				errorMessage: "Oops, something went wrong",
			},
			want:    `{"state":"errored", "error_info": "Oops, something went wrong"}`,
			wantErr: false,
		},
	}
	for _, testCase := range tests {
		mainTestCtx.Run(testCase.name, func(testContext *testing.T) {
			client := testCase.client

			got, err := client.Error(testCase.args.state, testCase.args.fileName, testCase.args.errorMessage)
			if (err != nil) != testCase.wantErr {
				testContext.Errorf("Client.Error() error = %v, wantErr %v", err, testCase.wantErr)
				return
			}
			if got != testCase.want {
				testContext.Errorf("Client.Error() = %v, want %v", got, testCase.want)
			}
		})
	}
}
