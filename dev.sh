#!/usr/bin/env bash

air &
templ generate --watch --proxy="localhost:8080" --open-browser=false &
npx tailwindcss -i ./assets/input.css -o ./assets/styles.css --minify --watch
