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

func TestClient_Update(t *testing.T) {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := tt.client

			got, err := client.Update(tt.args.state, tt.args.fileName)

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Error(t *testing.T) {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := tt.client

			got, err := client.Error(tt.args.state, tt.args.fileName, tt.args.errorMessage)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Error() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
