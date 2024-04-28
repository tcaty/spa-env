SHELL := /bin/bash

IMAGE=tcaty/spa-env

.PHONY: restore
restore:
	rm -rf data/app
	cp -r data/orig/ data/app/

.PHONY: replace
replace:
	source data/.env && \
	go run main.go replace \
		--workdir data/app \
		--dotenv .env.production 

.PHONY: cmd
cmd:
	source data/.env && \
	go run main.go replace \
		--workdir data/app \
		--dotenv .env.production \
		--cmd "while true; do echo 1; sleep 1; done"
		
.PHONY: deploy 
deploy:
	docker build -t ${IMAGE} . && docker push ${IMAGE}
	