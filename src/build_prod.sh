#!/bin/bash
cd frontend
rm -rf build
mkdir build
cd ../../frontend
pnpm install --frozen-lockfile
pnpm run build
cp -r build/* ../src/frontend/build/
cd ../src
go build -o localstream ./cmd/server
