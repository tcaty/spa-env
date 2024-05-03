# -- build stage --
FROM node:21.5.0-alpine as build

WORKDIR /usr/src/react

COPY package*.json .
RUN yarn install --silent

COPY ./ .
RUN yarn build


# -- runtime stage --
FROM node:21.5.0-alpine as production

WORKDIR /usr/src/react

# -- copy static files
COPY --from=build /usr/src/react/dist dist/
# -- copy spa-env binary from official image
COPY --from=tcaty/spa-env /spa-env /spa-env
# -- WARNING: vite doesn't copy .env file to dist by default
# -- you should copy it manually
COPY .env.production ./dist

RUN yarn add serve --silent

ENV HOST="0.0.0.0" \
    PORT="3000"

ENTRYPOINT [ \
    "/spa-env", "replace", \
    "--workdir", "/usr/src/react/dist", \
    "--dotenv", ".env.production", \
    "--cmd", "npx serve -s dist -l tcp://$HOST:$PORT", \
    "--form", "shell" \
]
