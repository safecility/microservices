@ECHO OFF

:choice
set /P c=deploy to PROD[Y/N]?
if /I "%c%" EQU "Y" goto :deploy
if /I "%c%" EQU "N" goto :exit
goto :choice

:deploy

echo "Deploying"

call go mod vendor
echo "Mod Vendor"
call gcloud config set project safecility-prod
call gcloud run deploy pipeline-usage-bigquery-device --source ./ --region "europe-west1"
call gcloud config set project safecility-test

:exit
echo "exiting"
