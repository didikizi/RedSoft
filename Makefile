export GO111MODULE := on
export GOOS=linux
export GOARCH=amd64

TARGET := testTask

################################################################################
### DataBase
###

test-db:
	sudo docker compose -f docker-compose.db.yml up --build --remove-orphans -d

################################################################################
### Builds
###

build:
	go build -o bin/$(TARGET) ./cmd/$(TARGET)/

################################################################################
### Linters
###

lint: tidy
	gofumpt -w .
	gci write . --skip-generated -s standard -s default
	make linters

linters: golangci-lint

golangci-lint: build
	test $(TARGET) != golang-template-project && grep -r --exclude-dir='.git' 'golang-template-project' | grep -v 'should be absent in normal project' && echo "all golang-template-project should be absent in normal project" && exit 1 || :
	find -type f -name "*.go" | grep -v '.*\.pb\.go' | grep -v '\/[a-z_]*.go' && echo "Files should be named in snake case" && exit 1 || echo "All files named in snake case"
	golangci-lint version
	golangci-lint run

################################################################################
### Golang helpers
###

tidy:
	go mod tidy

clean:
	rm -rf bin/$(TARGET)

download:
	go mod download

strip: build
	strip bin/$(TARGET)

modup:
	go get -u ./...
	go mod tidy

run:
	go mod tidy
	go run ./cmd/testTask/main.go
	
docker-compact-prod:
	docker compose -f docker-compose.prod.yml up --build --force-recreate --renew-anon-volumes --remove-orphans -d 
