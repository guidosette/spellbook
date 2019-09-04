#!/usr/bin/env bash

# run the datastore emulator on different process
osascript -e 'tell application "Terminal" to do script "gcloud beta emulators datastore start --no-store-on-disk"'

sleep 10
$(gcloud beta emulators datastore env-init)

dev_appserver.py --application=spellbook --env_var DATASTORE_PROJECT_ID=dev-app app.yaml --clear_search_indexes=true --default_gcs_bucket mage-middleware.appspot.com

