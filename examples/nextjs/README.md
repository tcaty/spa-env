# NextJS example

This is simple example. Please, pay attention to these files:
* `.env.development` - environment variables for development mode. There are `key=value` entries.
* `.env.production` - environment variable for production mode. There are `key=placeholder` entries.
* `Dockerfile` - the place to use `spa-env`.
* `docker-compose.yml` - the place to configure environment variables `placeholder=value`.

## Running

To run this example just execute command below:
```
docker-compose up -d
```
And checkout `http://0.0.0.0:3000` in your browser. Also checkout container logs to see nextjs server side logs.

## Usage

### Run `spa-env`

You could find detailed description of `spa-env` usage in `Dockerfile`. 

> [!NOTE]
> NextJS copies `.env.production` file to dist folder automatically.

### Environment variables

There are three environment variables in this example:
* `POSTGRES_CONN_STRING` - server side variable
* `NEXT_PUBLIC_API_URL` - client side variable
* `NEXT_PUBLIC_SECRET_TOKEN` - client side variable

As it's known, server side code has access to environment variables during node runtime without any problems. So, we must just skip server side variables. To do this, there is a specified flag `--prefifx=NEXT_PUBLIC` in `Dockerfile`, variables without this prefix will be skipped. 
