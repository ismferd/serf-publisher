APPLICATION      := serf-publisher
LINUX            := build/${APPLICATION}-linux-amd64
DARWIN           := build/${APPLICATION}-darwin-amd64
ARM              := build/${APPLICATION}-arm						
DOCKER_USER      ?= ""
DOCKER_PASS      ?= ""
BIN_DIR          := $(GOPATH)/bin
GOMETALINTER     := $(BIN_DIR)/gometalinter
COVER            := $(BIN_DIR)/gocov-xml
JUNITREPORT      := $(BIN_DIR)/go-junit-report
TRAVIS_COMMIT    ?= latest


.PHONY: $(DARWIN)
$(DARWIN):
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -o ${DARWIN} *.go

.PHONY: linux
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ${LINUX} *.go

.PHONY: arm
arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -a -installsuffix cgo -o ${ARM} *.go
	
.PHONY: release
release:
	echo "${DOCKER_PASS}" | docker login -u "${DOCKER_USER}" --password-stdin
	docker build -t "${DOCKER_IMAGE}" "."
	docker tag "${DOCKER_USER}/""${DOCKER_IMAGE}" "${DOCKER_IMAGE}:0.1.0"
	docker push "${DOCKER_IMAGE}"

e2e:
	./e2e_test.sh
