# lab-assignment-system-frontend

## Development

- `.env.localhost`

```shell
VITE_BACKEND_HOST=http://localhost:8080
```

```shell
$ yarn
$ yarn dev
```

## Deploy

```shell
$ yarn firebase login
```

### test 環境

- `.env.test`

```shell
VITE_BACKEND_HOST=
```

```shell
$ yarn deploy-test
```

### production 環境

- `.env.production`

```shell
VITE_BACKEND_HOST=
```

```shell
$ yarn deploy-prod
```
