# lab-assignment-system

lab-assignment-system is a web application which assigns students to laboratories in Shizuoka University.

## Screenshots

**top page**

<img src=https://user-images.githubusercontent.com/43411965/184492625-ab8d031c-a586-4de3-ae9e-a1fd45281e9f.png width=600>

**laboratory list**

<img src=https://user-images.githubusercontent.com/43411965/184492690-f4cb5786-5e32-41dd-acb3-dd6d8f729f90.png width=600>


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
