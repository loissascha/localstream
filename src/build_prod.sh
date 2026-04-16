#!/bin/bash
cd frontend
bun install
bun run build
cd ..
go build -o localstream ./cmd/server
