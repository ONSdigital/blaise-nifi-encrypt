
# Overview

This package contains three seperate cloud functions that operate as described below.

## Extract Function

This function receives an event from Gcloud pub/pub with a payload describing the 
instrument the user wishes to extract csv response data for. The function  
matches response data to the MI fields in the database and writes a CSV file to a storage bucket
defined in the `EXTRACT_OUTPUT` environment variable. Note that this function is specific to the mi-extract functionality
and is therefore not reusable.
 
## Encrypt Function

The encrypt function is triggered when a file arrives in a bucket, which is defined in the function configuration. 
The file is encrypted using the build-in Golang PGP encryption functions with the stipulated public key and the 
result placed in the bucket identified by the `ENCRYPT_OUTPUT` environment variable. 

The Golang libraries allow for the streaming of data into and out of the encryption routines with the result being 
that any sized file can be encrypted without being constrained by memory 
or storage considerations.

## Zip Function

The zip function is triggered when a file arrives in a bucket, which is defined in the 
function configuration. 

The file is zipped and placed in the bucket identified by the `ZIP_OUTPUT` environment variable
 using the following file format:

`mi[1]_[2]_[3].zip` where:
1. Is the name of the instrument
2. Is a date in the format DDMMYYYY
3. Is the time is the format HHMMSS



# Configuration

### Google Functions Region Setting

Set the default functions region:

`gcloud config set functions/region europe-west2`

Otherwise, functions will be created somewhere far away in the ether...

### Environment Variables

The following environment variables are available (see the testing section for details on how to create buckets):

* `PUBLIC_KEY=<path to gpg public key file>` - required to encrypt the zip file

* `GPG_EXTENSION=true|false` - whether to add a .gpg extension to the encrypted file. Default is false.

* `EXTRACT_OUTPUT=<bucket>` - the GCloud bucket where the extracted data is places by the `extract` function.

* `ENCRYPT_OUTPUT=<bucket>` - the GCloud bucket where the file that has been encrypted is located. 
Placed there by the `encrypt_function`.

* `ZIP_OUTPUT=<bucket>` - the GCloud bucket where the file that has been zipped is placed. Placed
there by the `zip_function`.

* `GOOGLE_APPLICATION_CREDENTIALS=<file>` - google credentials file

* `LOG_FORMAT=Terminal|Json` - (json is the default) for logging messages. 
If you want pretty coloured output for local testing use `Terminal`

* `Debug=True|False|NotSet` - set debugging

* `DB_SERVER=<server>` - server address

* `DB_SOCKET_DIR` - the name of the Unix domain socket used by the GCloud SQL instance. Should be set to `/cloudsql` for 
production deployment, unset for testing. 

* `DB_DATABASE=<database>` - the name of the database, defaults to 'blaise'

* `DB_USER=<user>` - the database user

* `DB_PASSWORD=password` - the database password

