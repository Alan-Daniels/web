#!/usr/bin/env bash

(trap 'kill 0' SIGINT; 
air &
templ generate --watch &
npx tailwindcss -i ./internal/input.css -o ./assets/styles.css --minify --watch
)
