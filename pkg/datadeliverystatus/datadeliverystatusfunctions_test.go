package datadeliverystatus

import (
	"net/http"
	"testing"
)

func TestClient_Update(t *testing.T) {
	type Client struct {
		Config *Config
		HTTP   *http.Client
	}

	tests := []struct {
		name     string
		state    string
		fileName string
		want     string
		wantErr  bool
	}{
		{
			name:     "DDS url is not set up",
			state:    "in_staging",
			fileName: "dummy",
			want:     "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Getenv("DDS_URL")
			client := &Client{
				Config: {BaseURL string: ""},
				HTTP:   tt.fields.HTTP,
			}
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
