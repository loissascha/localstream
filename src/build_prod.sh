#!/bin/bash
cd frontend
bun install --frozen-lockfile
bun run build
cd ..
go build -o localstream ./cmd/server
