
# Blaise NiFi Encrypt 🔣

The encrypt function is triggered when a file arrives in a GCS bucket (configured as the function trigger). The file is encrypted using `golang.org/x/crypto/openpgp` with the configured public key, and the result is written to the bucket identified by the `ENCRYPTION_DESTINATION` environment variable.

Data is streamed directly between GCS and the encryption routines, so files of any size can be encrypted without memory or storage constraints.

## Environment variables

The following environment variables are required:

* `PUBLIC_KEY=<path to gpg public key file>` - path to the PGP public key used to encrypt files

* `ENCRYPTION_DESTINATION=<bucket>` - GCS bucket where encrypted files are written

* `DDS_URL=<url>` - base URL of the Data Delivery Status service

* `CLIENT_ID=<client-id>` - IAP client ID used to authenticate requests to the Data Delivery Status service

The following environment variables are optional:

* `LOG_FORMAT=Terminal|Json` - log output format (default: `Json`). Use `Terminal` for human-readable coloured output during local development.

* `DEBUG=true` - enables debug-level logging when set to `true`

## Running Tests

Run the test suite from the repository root:

```sh
go test -v ./...
```

To verify changes against a live environment:

* Point your sandbox to the feature branch you are working on
* Upload a sample file to the `nifi_staging` bucket in your GCP project (`ons-blaise-v2-dev-<sandbox-suffix>`)
* Confirm the file is encrypted and written to the `nifi` bucket in the same project

## Encryption key management

The Ingest team are responsible for generating the PGP encryption keys. When they expire, the Ingst team will provide new public keys which will need to be [committed to this repository](https://github.com/ONSdigital/blaise-nifi-encrypt/tree/main/pkg/encryption/keys).
