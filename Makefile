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

.PHONY: nextjs
nextjs:
	./scripts/use-env.sh ${NEXTJS}/.env && \
	go run main.go replace \
		--workdir ${NEXTJS}/.next \
		--dotenv .env.production \
		--cmd "while true; do echo 1; sleep 1; done" \
		--form shell \
		--verbose

.PHONY: react
react:
	./scripts/use-env.sh ${REACT}/.env && \
	go run main.go replace \
		--workdir ${REACT}/dist \
		--dotenv .env.production \
		--cmd "echo react" \
		--verbose