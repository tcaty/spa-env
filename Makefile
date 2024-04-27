SHELL := /bin/bash

.PHONY: restore
restore:
	rm -rf data/app
	cp -r data/orig/ data/app/

.PHONY: dev
dev:
	source data/.env && \
	go run main.go replace \
		--workdir data/app \
		--env-file .env.production
