#!/bin/bash

gcloud functions deploy ZipFunction --runtime go113 \
  --trigger-resource ons-blaise-dev-pds-20-mi-encrypted \
  --region=europe-west2 --trigger-event google.storage.object.finalize \
  --set-env-vars ZIP_LOCATION=ons-blaise-dev-pds-20-mi-zip


