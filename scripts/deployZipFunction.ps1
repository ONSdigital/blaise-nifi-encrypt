
# run this from the project directory as scripts/deployZipFunction

gcloud functions deploy ZipFunction --runtime go113 `
  --trigger-resource ons-blaise-dev-pds-27-zip `
  --region=europe-west2 --trigger-event google.storage.object.finalize `
  --set-env-vars ZIP_OUTPUT=ons-blaise-dev-pds-27-encrypt
