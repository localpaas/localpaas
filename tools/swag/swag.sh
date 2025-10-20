#!/bin/bash
set -eo pipefail

DEVTOOLS_IMAGE=devtools

# Gen swagger.json
docker run --entrypoint "/go/bin/swag" --rm --volume "${PWD}":/app --volume "${HOME}/go/pkg/mod":/go/pkg/mod ${DEVTOOLS_IMAGE} init \
  -g app/interface/api/server/gin.go -o docs/openapi --outputTypes json \
  --parseDependencyLevel 3 --requiredByDefault

# Convert swagger.json to OpenAPI v3 format
docker run --user $(id -u) --rm --volume "${PWD}:/app" openapitools/openapi-generator-cli generate \
  --skip-validate-spec -i /localpaas_app/docs/openapi/swagger.json -g openapi -o /localpaas_app/tmp/swago && \
  cp tmp/swago/openapi.json docs/openapi/swagger.json && rm -rf tmp/swago

# Support additional features of OpenAPI v3
docker run --rm --volume "${PWD}":/app ${DEVTOOLS_IMAGE} python3 tools/swag/support-openapi-v3.py
