export VAULT_TOKEN := root
export CONTAINER_NAME := vault-dev
export VAULT_ADDR := http://0.0.0.0:8200

GOARCH := amd64
COMMIT := $(shell git rev-parse HEAD)
LDFLAGS := "-X main.COMMIT=${COMMIT}"
BINARY := traverse

darwin:
	GOOS=darwin GOARCH=${GOARCH} go build -ldflags ${LDFLAGS} -o ${BINARY}

test:
	go test -v -count=1 ./...

spinup: create-cluster sleep write-secrets
	
# spin up a new vault cluster in docker for testing purpsoses. running in the background in dev mode.
create-cluster:
	docker run -d --cap-add=IPC_LOCK -e VAULT_DEV_ROOT_TOKEN_ID=${VAULT_TOKEN} --name ${CONTAINER_NAME} -p 8200:8200 vault

write-secrets: 
	vault kv put secret/hello foo=world
	vault kv put secret/aws/credentials AWS_ACCESS_KEY_ID=AKIA-123 AWS_SECRET_ACCESS_KEY=nahnahnah
	vault kv put secret/team/team1/test/PD_TOKEN value=abcdef123456
	vault kv put secret/team/team1/test/TWITTER_TOKEN value=pffffffffffft
	vault kv put secret/team/team1/test/k8s/SERVICE_TOKEN value=2
	vault kv put secret/team/team1/test/k8s/CA_DATA value=test_ca_data
	vault kv put secret/team/team1/test/k8s/API_SERVER_URL value=localhost:8000
	vault kv put secret/team/notrad/production/k8s/SERVICE_TOKEN value=nothingherebutitsprod
	vault kv put secret/team/notrad/production/CA_DATA value=ls0f123mnbsdhjsdkljsdf

# add in a sleep to make sure all vault spinup activities are finished before writing new secrets
sleep:
	sleep 5

# clean throws errors if docker is not running or the container is not running/does not exist
clean:
	rm -f traverse
	docker kill vault-dev
	docker rm vault-dev

.PHONY: test spinup create-cluster write-secrets sleep
