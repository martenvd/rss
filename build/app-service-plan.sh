#!/usr/bin/env bash

## https://learn.microsoft.com/en-us/azure/app-service/tutorial-multi-container-app
az webapp create --resource-group test-rss-feed \
                 --plan test-rss-feed-plan \
                 --subscription "47c8daef-faf3-4422-a4f0-76da5082a559" \
                 --name test-rss-feed-marten \
                 --multicontainer-config-type compose \
                 --multicontainer-config-file docker-compose-app-service.yml

# az webapp config container set --resource-group test-rss-feed \
#                  --subscription "47c8daef-faf3-4422-a4f0-76da5082a559" \
#                  --name test-rss-feed-marten \
#                  --multicontainer-config-type compose \
#                  --multicontainer-config-file docker-compose-app-service.yml