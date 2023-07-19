.PHONY: datastore-emulator/start
datastore-emulator/start:
	@docker run -d \
		-e CLOUDSDK_CONFIG="/config/gcloud" \
		-e DATASTORE_PROJECT_ID=${GCP_PROJECT_ID} \
		-p ${DATASTORE_EMULATOR_CONTAINER_INTERNAL_PORT} \
		-v "${HOME}/.config/gcloud":"/config/gcloud" \
		-v "$(shell pwd)/.datastore-exports:/.datastore-exports" \
		--name ${DATASTORE_EMULATOR_CONTAINER_NAME} \
		google/cloud-sdk:emulators \
		gcloud beta emulators datastore start \
			--use-firestore-in-datastore-mode \
			--host-port="0.0.0.0:${DATASTORE_EMULATOR_CONTAINER_INTERNAL_PORT}" \
			--no-store-on-disk

.PHONY: datastore-emulator/stop
datastore-emulator/stop:
	@docker stop ${DATASTORE_EMULATOR_CONTAINER_NAME}
	@docker rm ${DATASTORE_EMULATOR_CONTAINER_NAME}

.PHONY: datastore-emulator/import
datastore-emulator/import:
	curl -X POST "http://localhost:${DATASTORE_EMULATOR_CONTAINER_PORT}/v1/projects/${GCP_PROJECT_ID}:import" \
		-H 'Content-Type: application/json' \
		-d '{"input_url":"/.datastore-exports/exports.overall_export_metadata"}'


.PHONY: datastore-emulator/show-port
datastore-emulator/show-port:
	@echo ${DATASTORE_EMULATOR_CONTAINER_PORT}

.PHONY: test
test:
	@docker run -d \
		-e CLOUDSDK_CONFIG="/config/gcloud" \
		-e DATASTORE_PROJECT_ID="lab-assignment-system-test" \
		-p 32932 \
		-v "${HOME}/.config/gcloud":"/config/gcloud" \
		--name lab-assignment-system_datastore-emulator_test \
		google/cloud-sdk:emulators \
		gcloud beta emulators datastore start \
			--use-firestore-in-datastore-mode \
			--host-port="0.0.0.0:32932" \
			--no-store-on-disk
	@DATASTORE_EMULATOR_CONTAINER_PORT=$(docker port lab-assignment-system_datastore-emulator_test  | grep "32932/tcp" | awk -F ' -> ' '{print $$2}' | awk -F ':' '{print $$2}') \
		DATASTORE_EMULATOR_HOST="0.0.0.0:${DATASTORE_EMULATOR_CONTAINER_PORT}" \
		go test -v ./... \
		; docker rm -f lab-assignment-system_datastore-emulator_test