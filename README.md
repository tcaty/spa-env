# spa-env

Have you ever tried to use environment variables in single page applications at browser runtime? There are several ways to do it, however almost all of them are too slow, too unsafe or have redundant complexity. So, `spa-env` is fast, reliable and simple solution of this problem.

## How does it work?

### Environment variables in spa

First of all there is a little note about environment variables in spa. Imagine that we have `app.js` file with code below:
```
console.log(process.env.VITE_API_URL)
```
And there is `.env` file as well:
```
VITE_API_URL=https://api.com/
```
As you've noticed we'll use `vite` to build this application. After running `yarn build` there will be built static file, which looks like this:
```
console.log("https://api.com/")
```
So, the main idea is that variables that refer on environment variables are replaced by static values at buildtime.

### `spa-env` working steps 
1. this tool automatically finds `.env` file in `workdir` by filename
2. further it parses `key=placeholder` pairs where `key` has specified `prefix` from `.env` file
3. further `spa-env` parses pairs `placeholder=value` from actual environment
4. in the end it simply replaces `placeholder` by `value` in all files from `workdir` except `.env` file  

## Usage

Common `Dockerfile` for nextjs apps looks like this: 
```
# deps stage
...

# build stage
...

# runtime stage
...

# run app
CMD ["node", "server.js"]
```
All `spa-env` usage could be injected in runtime stage in `Dockerfile` of target spa. So, `Dockerfile` for nextjs app turns into this:
```
# deps stage
...

# build stage
...

# runtime stage
...

# download binary from official image
COPY --from=tcaty/spa-env /spa-env /spa-env

# run binary
ENTRYPOINT [ \
    "/spa-env", "replace", \
    "--workdir", "/app", \
    "--dotenv", ".env.production" \
    "--prefix", "NEXT_PUBLIC", \
    "--cmd", "node server.js", \
    "--verbose" \
]

```

## Examples

There are two available examples in `examples` folder:
* nextjs - simple nextjs app
* react - simple react app with vite
