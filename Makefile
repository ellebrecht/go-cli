#
# geeny-cli-go makefile
#

ifeq ($(AWS_PROFILE),)
  AWS_PROFILE := "default"
endif

# deploy bins to s3
s3deploy: 
	aws s3 sync $(CIRCLE_ARTIFACTS)/bin/ s3://developers.geeny.io/downloads/cli/ --delete --profile "$(AWS_PROFILE)"

build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/geeny -v -ldflags "-s -w -X version.timestamp=$(date -u '+%d.%m.%Y_%H:%M:%S') -X version.version=$(git describe --tags)" geeny

smoketest: build-linux
	docker-compose -f docker-compose.test.yml up -d downloads_server
	docker-compose -f docker-compose.test.yml run wait_for_server
	AWS_ECR_REPOSITORY=$(AWS_ECR_REPOSITORY) docker-compose -f docker-compose.test.yml run smoketest

# chmod and zip up all the bins
PATHS = linux/amd64 linux/x86 linux/arm7 linux/arm64 windows/amd64 windows/x86 osx/amd64
chmod-zip:
	for path in $(PATHS); do\
	  cd "$(CIRCLE_ARTIFACTS)/bin/$$path" && chmod 0555 geeny* && md5sum geeny* > geeny.md5 && zip geeny.zip geeny*;\
	done;

coverage:
	echo "mode: atomic" > coverage.txt;
	for d in `go list geeny/... | grep -v vendor`; do\
      go test -v -coverprofile=profile.out -covermode=atomic "$$d";\
        if [ -f profile.out ]; then\
          tail -n +2 profile.out >> coverage.txt;\
          rm profile.out;\
        fi;\
	done;\
	go tool cover -func=coverage.txt -o coverage.func.txt
	go tool cover -html=coverage.txt -o coverage.html
