#!/bin/bash

gcloud functions deploy EncryptFunction --runtime go113 --trigger-resource ons-blaise-dev-pds-20-mi-encrypt \
  --trigger-event google.storage.object.finalize --set-env-vars ENCRYPTED_LOCATION='ons-blaise-dev-pds-20-mi-encrypted',PUBLIC_KEY='./serverless_function_source_code/pkg/encryption/keys/dev.gpg'

