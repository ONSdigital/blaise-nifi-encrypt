
# Blaise NiFi Encrypt

The encrypt function is triggered when a file arrives in a bucket, which is defined in the function configuration.
The file is encrypted using the build-in Golang PGP encryption functions with the stipulated public key and the
result placed in the bucket identified by the `ENCRYPTION_DESTINATION` environment variable.

The Golang libraries allow for the streaming of data into and out of the encryption routines with the result being
that any sized file can be encrypted without being constrained by memory
or storage considerations.

## Configuration

### Google Functions region setting

Set the default functions region:

`gcloud config set functions/region europe-west2`

Otherwise, functions will be created somewhere far away in the ether...

### Environment variables

The following environment variables are available (see the testing section for details on how to create buckets):

* `PUBLIC_KEY=<path to gpg public key file>` - required to encrypt the zip file

* `ENCRYPTION_DESTINATION=<bucket>` - the GCloud bucket where the file that has been encrypted is located.
Placed there by the `nifi_encrypt_function`.

* `GOOGLE_APPLICATION_CREDENTIALS=<file>` - google credentials file

* `LOG_FORMAT=Terminal|Json` - (json is the default) for logging messages.
If you want pretty coloured output for local testing use `Terminal`

* `Debug=True|False|NotSet` - set debugging

## Manual deploy

```sh
gcloud functions deploy NiFiEncryptFunction \
  --runtime go113 \
  --region=europe-west2 \
  --trigger-event=google.storage.object.finalize \
  --trigger-resource="${NIFI_STAGING_BUCKET}" \
  --set-env-vars PUBLIC_KEY="./serverless_function_source_code/pkg/encryption/keys/${ENV}-key.gpg,ENCRYPTION_DESTINATION=${NIFI_BUCKET}"
```

## Running Test Cases

To run the test cases, you need to be in the directory where TESTFILE_test.go file exists.
Run the following command to run the test cases and see the results
```sh
go test -v
```
To make sure that the changes you have made are working and have not broken anything:

* Point your sandbox to the feature Branch that you are working on
* Put some sample file into "nifi_staging" Bucket in your GCP Project ons-blaise-v2-dev-<sandbox-suffix>
* Make sure everything is working expectedly, concourse/sandbox is not broken and it encrypts the file and places it in "nifi" Bucket for your GCP Project ons-blaise-v2-dev-<sandbox-suffix>

## Encryption key management

The CIA team are responsible for generating the PGP encryption keys. When they expire the CIA team will send us new public keys, they will need to [overwrite the existing public keys in this repository](https://github.com/ONSdigital/blaise-nifi-encrypt/tree/main/pkg/encryption/keys).
