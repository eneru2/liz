#!/bin/sh

PORT = $1
echo "Starting server on port $PORT..."

templ generate --watch --proxy="http://localhost:$PORT" --cmd="./.liz/bin/script.sh"
go run .