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

# copy static files
COPY --from=build /usr/src/react/dist dist/
# copy spa-env binary from official image
COPY --from=tcaty/spa-env /spa-env /spa-env
# WARNING: 
# vite doesn't copy .env file to dist by default
# you should copy it manually
COPY .env.production ./dist

RUN yarn add serve --silent

ENV HOSTNAME="0.0.0.0" \
    PORT="3000"
    
EXPOSE 3000

ENTRYPOINT [ \
    # execute binary
    "/spa-env", "replace", \
    # specify workdir where built static is placed
    "--workdir", "/usr/src/react/dist", \
    # specify name of .env file for production mode
    "--dotenv", ".env.production", \
    # specify command which will be run after successfull replacement
    "--cmd", "npx serve -s dist -l tcp://$HOSTNAME:$PORT", \
    # specify in which form command will be run
    # as you can see, command above uses environment variables $HOSTNAME and $PORT
    # therefore, it has to be run in shell form (by default command runs in exec form)
    "--form", "shell" \
]
