#!/bin/bash

gcloud pubsub topics create ons-blaise-dev-pds-27-extract-topic

gcloud functions deploy ExtractFunction --runtime go113 --trigger-topic ons-blaise-dev-pds-27-extract-topic --region europe-west2 \
  --set-env-vars ENCRYPT_LOCATION=ons-blaise-dev-pds-27-encrypt,DB_SERVER='ons-blaise-dev-pds-27:europe-west2:blaise-dev-068d804a',DB_USER='blaise',DB_PASSWORD='Xkjhb2vqVLZ4oo_D',DB_DATABASE='blaise',DB_SOCKET_DIR='/cloudsql'


