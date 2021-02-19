
# run this from the project directory as scripts/deployEncryptFunction

gcloud functions deploy EncryptFunction --runtime go113 `
  --trigger-resource ons-blaise-dev-pds-27-encrypt `
  --region=europe-west2 --trigger-event google.storage.object.finalize `
  --set-env-vars ENCRYPT_OUTPUT=ons-blaise-dev-pds-27-encrypted,GPG_EXTENSION=false,PUBLIC_KEY='./serverless_function_source_code/pkg/encryption/keys/dev.gpg'
