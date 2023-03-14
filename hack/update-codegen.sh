#!/bin/bash

set -x

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

vendor/k8s.io/code-generator/generate-groups.sh all \
  github.com/anisurrahman75/my-crd/pkg/client \
  github.com/anisurrahman75/my-crd/pkg/apis \
  mycrd.dev:v1 \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt