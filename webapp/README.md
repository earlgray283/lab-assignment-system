# lab-assignment-system-frontend

## Requirements

- Node.js, bun
    - recommend: volta(https://docs.volta.sh/guide/getting-started)

## Development

- `.env.localhost`

```shell
VITE_BACKEND_HOST=http://localhost:8080
```

```shell
$ bun
$ bun dev
```

## Deploy

```shell
$ bun firebase login
```

### test 環境

- `.env.test`

```shell
VITE_BACKEND_HOST=
```

```shell
$ bun deploy-test
```

### production 環境

- `.env.production`

```shell
VITE_BACKEND_HOST=
```

```shell
$ bun deploy-prod
```
