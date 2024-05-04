# TODO: refactor

SHELL := /bin/bash

.PHONY: nextjs
nextjs:
	source examples/nextjs/.env.development && \
	go run main.go replace \
		--workdir examples/nextjs/.next \
		--dotenv .env.production \
		--cmd "while true; do echo 1; sleep 1; done" \
		--form shell \
		--verbose
