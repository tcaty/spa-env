SHELL := /bin/bash

NEXTJS=examples/nextjs
REACT=examples/react

.PHONY: prepare
prepare:
	yarn --cwd ${NEXTJS} install && yarn --cwd ${NEXTJS} build && \
	rm -rf ${NEXTJS}/.next.backup && \
	cp -r ${NEXTJS}/.next ${NEXTJS}/.next.backup && \
	yarn --cwd ${REACT} install && yarn --cwd ${REACT} build && \
	rm -rf ${REACT}/dist.backup && \
	cp ${REACT}/.env.production ${REACT}/dist/ && \
	cp -r ${REACT}/dist ${REACT}/dist.backup 
	
.PHONY: restore
restore:
	rm -rf ${NEXTJS}/.next && \
	cp -r ${NEXTJS}/.next.backup ${NEXTJS}/.next && \
	rm -rf ${REACT}/dist && \
	cp -r ${REACT}/dist.backup ${REACT}/dist 

.PHONY: replace-nextjs
replace-nextjs:
	export $(shell grep -v '^#' ${NEXTJS}/.env | xargs -d '\n') && \
	go run main.go replace \
		--workdir ${NEXTJS}/.next \
		--dotenv .env.production \
		--key-prefix NEXT_PUBLIC \
		--placeholder-prefix PLACEHOLDER \
		--cmd "while true; do echo 1; sleep 1; done" \
		--cmd-form shell \
		--log-level DEBUG

.PHONY: react
replace-react:
	export $(shell grep -v '^#' ${REACT}/.env | xargs -d '\n') && \
	go run main.go replace \
		--workdir ${REACT}/dist \
		--dotenv .env.production \
		--placeholder-prefix PLACEHOLDER \
		--cmd "echo react" \
		--log-level DEBUG
	
.PHONY: generate-nextjs
generate-nextjs:
	go run main.go generate \
		--workdir ${NEXTJS} \
		--dotenv-dev .env.development \
		--dotenv-prod .env.production \
		--key-prefix NEXT_PUBLIC \
		--placeholder-prefix PLACEHOLDER \
		--enable-comments \
		--log-level DEBUG
	
.PHONY: generate-react
generate-react:
	go run main.go generate \
		--workdir ${REACT} \
		--dotenv-dev .env.development \
		--dotenv-prod .env.production \
		--key-prefix VITE \
		--placeholder-prefix PLACEHOLDER_ \
		--log-level DEBUG