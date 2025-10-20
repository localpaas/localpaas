#!/bin/bash
set -eo pipefail

TEST_RESULT_DIR="${TEST_RESULTS:-./test-results}"
mkdir -p ${TEST_RESULT_DIR}

PKG_LIST=$(go list ./...\
  | sed -e "s/github.com\/localpaas\/localpaas/./g"\
  | grep -v\
      -e /mock\
      -e ^./cmd\
      -e ^./config\
      -e ^./db\
      -e ^./docs\
      -e ^./scripts\
      -e ^./tests\
      -e ^./test-results\
      -e ^./tools\
      -e ^./localpaas_app/interface/api/docs\
  | tr '\n' ',')

echo "---------------------------------------------------------------"
echo "Test:"
echo "---------------------------------------------------------------"
go test -race\
  -coverpkg=${PKG_LIST}\
  -coverprofile ${TEST_RESULT_DIR}/.testCoverage.txt\
  ./...

echo "---------------------------------------------------------------"
echo "Result:"
echo "---------------------------------------------------------------"
go tool cover -func ${TEST_RESULT_DIR}/.testCoverage.txt
go tool cover -html=${TEST_RESULT_DIR}/.testCoverage.txt -o ${TEST_RESULT_DIR}/coverage.html

echo "---------------------------------------------------------------"
echo "DONE."
echo "---------------------------------------------------------------"
