#!/bin/sh
set -eu

source_dir="${1:?source directory is required}"
work_dir="${2:?work directory is required}"

test -f "$source_dir/go.mod"
mkdir -p "$work_dir"
find "$work_dir" -mindepth 1 -maxdepth 1 -exec rm -rf -- {} +
cp -a "$source_dir/." "$work_dir/"

for module in /modules/*; do
  if [ -f "$module/Makefile" ]; then
    make -C "$module" all TOOLS_DIR=/usr/local/bin \
      PROTOC=/usr/local/bin/protoc OAPI_CODEGEN=/usr/local/bin/oapi-codegen
  fi
done

cd "$work_dir"
go mod edit -replace github.com/gorundebug/servicelib=/servicelib
go mod edit -replace github.com/gorundebug/inventory_service_api=/modules/inventory_service_api
go mod edit -replace github.com/gorundebug/model_go=/modules/model_go
go mod edit -replace github.com/gorundebug/order_service_api=/modules/order_service_api
make -f make.generated.mk gen-proto TOOLS_DIR=/usr/local/bin PROTOC=/usr/local/bin/protoc