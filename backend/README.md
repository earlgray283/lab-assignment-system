# lab-assignment-system-backend

## Requirements

- Go(https://go.dev/dl/)
- gcloud CLI(https://cloud.google.com/sdk/docs/install?hl=ja)
- Docker(https://docs.docker.com/desktop/)
- direnv(https://github.com/direnv/direnv)

## Setup

### configure gcp project

[GCP 公式ドキュメント](https://cloud.google.com/resource-manager/docs/creating-managing-projects?hl=ja) を参考にして GCP プロジェクトを作成し，ローカルで project id をセットしてください．

```shell
$ gcloud auth login

$ gcloud projects list
PROJECT_ID                     NAME                    PROJECT_NUMBER
lab-assignment-system-project  lab-assignment-system   31415926535897

$ gcloud config set project lab-assignment-system-project
Updated property [core/project].

$ gcloud config get project
lab-assignment-system-project
```

### launch datastore emulator

lab-assignment-system-backend は cloud datastore(NoSQL DB) に依存しています．  
datastore はクラウドサービスですが，emulator を用いることでローカルで実行させることが可能です．  
ここで，デフォルトだと container が削除された時点でデータが削除されることに注意してください．データを永続化させたい場合は `--no-store-on-disk` コマンドを消してください．

```shell
$ make datastore-emulator/start
```

もしエクスポートされたデータがあれば以下の手順で import することが可能です．

1. `.datastore-exports/exports.overall_export_metadata` として置く
2. 以下のコマンドを実行する

```shell
$ mkdir .datastore-exports
$ gsutil -m cp -r \       
  "gs://lab-assignment-system-backup/exports/default_namespace" \
  "gs://lab-assignment-system-backup/exports/exports.overall_export_metadata" \
  .datastore-exports/
$ make datastore-emulator/import
```

うまくいけば import されるはずです．

emulator を stop & remove したい場合は以下のコマンドを実行してください．

```shell
$ make datastore-emulator/stop
```

## Run

### launch server

```shell
$ go run ./cmd/server
```

## Batch

### load users

ユーザの一覧を csv からロードして datastore に insert します．  
※ ヘッダはつけないでください

**csv format**

```csv
12345678,4.0,<admin|audience>
```

```shell
$ go run ./cmd/batch/load-users <path-to-users-csv> -year <year>
```

### load labs

研究室の一覧を csv からロードして datastore に insert します．  
※ ヘッダはつけないでください

**csv format**

```csv
szpp,SZPP研究室,3776
```

```shell
$ go run ./cmd/batch/load-labs <path-to-users-csv> -year <year>
```

### create-survey

アンケートの集計期間をセットします。

```shell
$ go run ./cmd/batch/create-survey -startAt 2023-07-15T00:00:00 -endAt 2023-07-22T15:00:00 -year 2023
```
