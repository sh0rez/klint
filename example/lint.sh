#!/usr/bin/env bash
set -euo pipefail

go run .. -R rules/namespaces/namespaces.klint.go \
          -R rules/no-latest/noLatest.go \
          -R rules/resource-requests/resourceRequests.go \
          deployment.yml
