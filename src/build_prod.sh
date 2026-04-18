#!/bin/bash
cd frontend
pnpm install --frozen-lockfile
pnpm run build
cd ..
go build -o localstream ./cmd/server
