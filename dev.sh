#!/usr/bin/env bash

setupdb () {
  sleep 1
  echo "DEFINE NAMESPACE IF NOT EXISTS localhost;
  USE NAMESPACE localhost;
  DEFINE DATABASE IF NOT EXISTS web;
  USE DATABASE web;
  DEFINE USER OVERWRITE localhost ON NAMESPACE PASSWORD \"dev-database-pass\" ROLES EDITOR;" | surreal sql -e http://127.0.0.1:7999 -u root -p root
}

doair () {
  sleep 2
  air
}

export DEVPWD=$(realpath .)
export CGO_ENABLED=0

mkdir -p ./tmp/db

npx tailwindcss -i ./internal/input.css -o ./assets/styles.css --minify
(trap 'kill 0' SIGINT; 
templ generate --watch &
npx tailwindcss -i ./internal/input.css -o ./assets/styles.css --minify --watch &
surreal start --bind 127.0.0.1:7999 -A --log debug -u root -p root "file://${DEVPWD}/tmp/db/" &
setupdb &
doair
)
