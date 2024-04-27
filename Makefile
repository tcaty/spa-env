SHELL := /bin/bash

.PHONY: restore
restore:
	rm -rf data/app
	cp -r data/orig/ data/app/

.PHONY: replace
replace:
	source data/.env && \
	go run main.go replace \
		--workdir data/app \
		--env-file .env.production 

.PHONY: entrypoint
entrypoint:
	source data/.env && \
	go run main.go replace \
		--workdir data/app \
		--env-file .env.production \
		sleep 10
