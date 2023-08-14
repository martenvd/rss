#!/usr/bin/env bash

## https://learn.microsoft.com/en-us/azure/app-service/tutorial-multi-container-app
az webapp create --resource-group '<your-resource-group-name>' \
                 --plan '<your-plan-name>' \
                 --subscription "*************" \
                 --name '<your-app-service-name>' \
                 --multicontainer-config-type compose \
                 --multicontainer-config-file docker-compose-app-service.yml

az webapp config appsettings set --resource-group '<your-resource-group-name>' --name <app-name> --settings WEBSITES_ENABLE_APP_SERVICE_STORAGE=TRUE

# az webapp config container set --resource-group '<your-resource-group-name>' \
#                  --subscription "*************" \
#                  --name '<your-app-service-name>' \
#                  --multicontainer-config-type compose \
#                  --multicontainer-config-file docker-compose-app-service.yml