# React + TypeScript + Vite example

This is simple example. Please, pay attention to these files:
* `.env.development` - environment variables for development mode.
* `.env.production` - environment variables for production mode. This file was generated automatically.
* `Dockerfile` - the place to use `spa-env`.
* `docker-compose.yml` - the place to configure environment variables. Here you can see `env` usage, it is for easy development. You should use keyword `environment`.

## Running

To run this example just execute command below:
```
docker-compose up -d
```
And checkout `http://0.0.0.0:3000` in your browser.

## Usage

### Run `spa-env`

You could find detailed description of `spa-env` usage in `Dockerfile`. 

> [!WARNING]
> Vite doesn't copy `.env.production` file to dist automatically. So it must be copied manually, see `Dockerfile` for more details.

### Environment variables

There are two environment variables in this example:
* `VITE_API_URL`
* `VITE_SECRET_TOKEN`
