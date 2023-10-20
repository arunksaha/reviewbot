EXE = reviewbot
TAR_PREFIX = $(EXE)-arun-saha
TARGZ = ./$(TAR_PREFIX).tar.gz

all: build docker-build

build:
	go build -o $(EXE)

docker-build:
	docker build -t $(EXE) .

clean:
	rm $(EXE) $(TARGZ) coverage.out -rf

check:
	go vet
	go fmt
	staticcheck

test:
	go test -v

test-coverage:
	go test -v -coverprofile=coverage.out
	go tool cover -func=coverage.out

# For viewing test coverage visually per line, try:
#		go tool cover -html=coverage.out 

# https://staticcheck.io/docs/getting-started/
install:
	go install honnef.co/go/tools/cmd/staticcheck@latest

archive:
	git archive --format=tar.gz -o $(TARGZ) --prefix=$(TAR_PREFIX)/ main

