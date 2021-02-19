package util

import (
	"time"
)

// common structures

const (
	ExtractOutput   = "EXTRACT_OUTPUT"
	EncryptOutput   = "ENCRYPT_OUTPUT"
	ZipOutput       = "ZIP_OUTPUT"
	UseGPGExtension = "GPG_EXTENSION"

	PublicKeyFile   = "PUBLIC_KEY"
	Server          = "DB_SERVER"
	User            = "DB_USER"
	Password        = "DB_PASSWORD"
	Database        = "DB_DATABASE"
	DefaultDatabase = "blase"
)

// google storage events data
type GCSEvent struct {
	Kind                    string                 `json:"kind"`
	ID                      string                 `json:"id"`
	SelfLink                string                 `json:"selfLink"`
	Name                    string                 `json:"name"`
	Bucket                  string                 `json:"bucket"`
	Generation              string                 `json:"generation"`
	Metageneration          string                 `json:"metageneration"`
	ContentType             string                 `json:"contentType"`
	TimeCreated             time.Time              `json:"timeCreated"`
	Updated                 time.Time              `json:"updated"`
	TemporaryHold           bool                   `json:"temporaryHold"`
	EventBasedHold          bool                   `json:"eventBasedHold"`
	RetentionExpirationTime time.Time              `json:"retentionExpirationTime"`
	StorageClass            string                 `json:"storageClass"`
	TimeStorageClassUpdated time.Time              `json:"timeStorageClassUpdated"`
	Size                    string                 `json:"size"`
	MD5Hash                 string                 `json:"md5Hash"`
	MediaLink               string                 `json:"mediaLink"`
	ContentEncoding         string                 `json:"contentEncoding"`
	ContentDisposition      string                 `json:"contentDisposition"`
	CacheControl            string                 `json:"cacheControl"`
	Metadata                map[string]interface{} `json:"metadata"`
	CRC32C                  string                 `json:"crc32c"`
	ComponentCount          int                    `json:"componentCount"`
	Etag                    string                 `json:"etag"`
	CustomerEncryption      struct {
		EncryptionAlgorithm string `json:"encryptionAlgorithm"`
		KeySha256           string `json:"keySha256"`
	}
	KMSKeyName    string `json:"kmsKeyName"`
	ResourceState string `json:"resourceState"`
}

type Zip struct {
	FileName     string
	FromLocation string
	ToLocation   string
	DeleteFile   bool
}

type Encrypt struct {
	KeyFile              string
	FileName             string
	Location             string
	EncryptedDestination string
	DeleteFile           bool
	UseGPGExtension      bool
}
