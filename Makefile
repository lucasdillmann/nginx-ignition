DOCKER_IMAGE ?= dillmann/nginx-ignition
VERSION ?= 0.0.0
PR_ID ?= 0

prerequisites:
	go work sync
	cd frontend/ && npm i

check: prerequisites
	cd frontend/ && npm run check

format: prerequisites
	cd frontend/ && npx prettier --write .

build-frontend: prerequisites
	cd frontend/ && npm run build

build-backend-amd64: prerequisites
	GOARCH=amd64 CGO_ENABLED="0" GOOS="linux" go build -o build/linux/amd64 application/main.go

build-backend-arm64: prerequisites
	GOARCH=arm64 CGO_ENABLED="0" GOOS="linux" go build -o build/linux/arm64 application/main.go

build-release-docker-image:
	docker buildx build \
		--tag $(DOCKER_IMAGE):$(VERSION) \
		--tag $(DOCKER_IMAGE):latest \
		--platform linux/amd64,linux/arm64 \
		--build-arg NGINX_IGNITION_VERSION="$(VERSION)" \
		--push .

build-snapshot-docker-image:
	docker buildx build \
		--tag $(DOCKER_IMAGE):pr-$(PR_ID)-snapshot \
		--platform linux/amd64,linux/arm64 \
		--build-arg NGINX_IGNITION_VERSION="" \
		--push .

build-prerequisites: prerequisites build-frontend build-backend-amd64 build-backend-arm64

build-release: build-prerequisites build-release-docker-image

build-snapshot: build-prerequisites build-snapshot-docker-image
