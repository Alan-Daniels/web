#!/usr/bin/env nix-shell
#!nix-shell -i bash -p bash tailwindcss

setupdb () {
  sleep 1
  echo "DEFINE NAMESPACE IF NOT EXISTS localhost;
  USE NAMESPACE localhost;
  DEFINE DATABASE IF NOT EXISTS web;
  USE DATABASE web;
  DEFINE USER OVERWRITE localhost ON NAMESPACE PASSWORD \"dev-database-pass\" ROLES EDITOR;" | surreal sql -e http://127.0.0.1:8000 -u root -p root
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
templ generate -source-map-visualisations -watch &
# tailwindcss -i ./internal/input.css -o ./assets/styles.css --minify --watch &
surreal start --bind 127.0.0.1:8000 -A --log debug -u root -p root "file://${DEVPWD}/tmp/db/" &
setupdb &
doair
)
