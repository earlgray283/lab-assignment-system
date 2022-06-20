# lab-assignment-system

lab-assignment-system is a system which assigns students to laboratories in Shizuoka University.

## Build

1. This system depends on [datastore](https://cloud.google.com/datastore) and [firebase authentication](https://firebase.google.com/docs/auth), you should create GCP project and Firebase project and put secret files.

* `backend/credentials.json` ... a service account credentials which has permissions to access datastore
* `webapp/.env.local` ... a service account credentials which has permissions to access firebase authentication

```
VITE_API_KEY=
VITE_AUTH_DOMAIN=
VITE_PROJECT_ID=
VITE_STORAGE_BUCKET=
VITE_MESSAGING_SENDER_ID=
VITE_APP_ID=
VITE_MEASUREMENT_ID=
```

* `webapp/.env.development.local` ... backend url(localhost)

```
VITE_BACKEND_HOST=http://localhost:8080
```

2. Build with docker compose(TODO)

```shell
$ docker compose up -d --build
```
