DOCKER_IMAGE ?= dillmann/nginx-ignition
PRIVATE_DOCKER_IMAGE ?= docker.lan.dillmann.dev/nginx-ignition
VERSION ?= 0.0.0
PR_ID ?= 0
BUILDKIT_CACHE ?= build/cache/docker
SNAPSHOT_TAG_SUFFIX := $(if $(filter-out 0,$(PR_ID)),pr-$(PR_ID)-snapshot,$(VERSION)-snapshot)
LDFLAGS := -X 'dillmann.com.br/nginx-ignition/core/common/version.Number=$(VERSION)'

.backend-prerequisites:
	go work sync

.frontend-prerequisites:
	cd frontend/ && pnpm install

.frontend-lint: .frontend-prerequisites
	cd frontend/ && pnpm run check

.backend-lint: .backend-prerequisites .backend-test-mocks
	go tool golangci-lint run \
		./api/... \
		./application/... \
		./certificate/commons/... \
		./certificate/custom/... \
		./certificate/external/... \
		./certificate/letsencrypt/... \
		./certificate/selfsigned/... \
		./core/... \
		./database/... \
		./integration/docker/... \
		./integration/truenas/... \
		./vpn/netbird/... \
		./vpn/tailscale/...

.frontend-build: .frontend-prerequisites .generate-i18n-files
	cd frontend/ && pnpm run build

.backend-build: .backend-prerequisites .generate-i18n-files
	$(MAKE) .backend-build-file OS=linux ARCH=amd64 DIR=linux
	$(MAKE) .backend-build-file OS=linux ARCH=arm64 DIR=linux
	$(MAKE) .backend-build-file OS=darwin ARCH=arm64 DIR=macos
	$(MAKE) .backend-build-file OS=windows ARCH=amd64 DIR=windows EXT=.exe
	$(MAKE) .backend-build-file OS=windows ARCH=arm64 DIR=windows EXT=.exe

.backend-build-file:
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o build/$(DIR)/$(ARCH)$(EXT) application/main.go

.generate-i18n-files:
	go run ./tools/i18n/

.build-release-docker-image:
	mkdir -p $(BUILDKIT_CACHE)
	docker buildx build \
		--cache-from type=local,src=$(BUILDKIT_CACHE) \
		--cache-to type=local,dest=$(BUILDKIT_CACHE),mode=max \
		--tag $(DOCKER_IMAGE):$(VERSION) \
		--tag $(DOCKER_IMAGE):latest \
		--tag $(PRIVATE_DOCKER_IMAGE):$(VERSION) \
		--tag $(PRIVATE_DOCKER_IMAGE):latest \
		--platform linux/amd64,linux/arm64 \
		--push .

.build-snapshot-docker-image:
	mkdir -p $(BUILDKIT_CACHE)
	docker buildx build \
		--cache-from type=local,src=$(BUILDKIT_CACHE) \
		--cache-to type=local,dest=$(BUILDKIT_CACHE),mode=max \
		--tag $(DOCKER_IMAGE):$(SNAPSHOT_TAG_SUFFIX) \
		--tag $(PRIVATE_DOCKER_IMAGE):$(SNAPSHOT_TAG_SUFFIX) \
		--platform linux/amd64,linux/arm64 \
		--push .

.build-distribution-files:
	$(MAKE) .build-distribution-zip ARCH=amd64 OS=linux SERVICE_FILE_EXT=service
	$(MAKE) .build-distribution-zip ARCH=arm64 OS=linux SERVICE_FILE_EXT=service
	$(MAKE) .build-distribution-zip ARCH=arm64 OS=macos SERVICE_FILE_EXT=plist
	$(MAKE) .build-distribution-zip ARCH=amd64 OS=windows BIN_EXT=.exe
	$(MAKE) .build-distribution-zip ARCH=arm64 OS=windows BIN_EXT=.exe
	$(MAKE) .build-distribution-packages ARCH=amd64 OS=linux
	$(MAKE) .build-distribution-packages ARCH=arm64 OS=linux

.build-distribution-zip:
	rm -Rf build/nginx-ignition.$(OS)-$(ARCH).zip
	mkdir -p build/zip
	cp -Rf frontend/build build/zip/frontend
	cp -Rf database/common/migrations/scripts build/zip/migrations
	cp dist/$(OS)/instructions.md build/zip/instructions.md
	cp dist/$(OS)/nginx-ignition.properties build/zip/
	[ -z "$(SERVICE_FILE_EXT)" ] || cp dist/$(OS)/nginx-ignition.$(SERVICE_FILE_EXT) build/zip/
	cp build/$(OS)/$(ARCH)$(BIN_EXT) build/zip/nginx-ignition$(BIN_EXT)
	cd build/zip && zip -q -r ../nginx-ignition-$(VERSION).$(OS)-$(ARCH).zip .
	rm -Rf build/zip

.build-distribution-packages:
	export VERSION=$(VERSION); \
	export OS=$(OS); \
	export ARCH=$(ARCH); \
	export PACKAGE_ARCH=$(ARCH); \
	envsubst < dist/linux/nfpm.yaml > build/nfpm.yaml
	nfpm package --config build/nfpm.yaml --packager deb --target build/nginx-ignition-$(VERSION).$(ARCH).deb
	nfpm package --config build/nfpm.yaml --packager rpm --target build/nginx-ignition-$(VERSION).$(ARCH).rpm
	nfpm package --config build/nfpm.yaml --packager apk --target build/nginx-ignition-$(VERSION).$(ARCH).apk
	nfpm package --config build/nfpm.yaml --packager archlinux --target build/nginx-ignition-$(VERSION).$(ARCH).pkg.tar.zst
	nfpm package --config build/nfpm.yaml --packager ipk --target build/nginx-ignition-$(VERSION).$(ARCH).ipk
	rm -Rf build/nfpm.yaml

.frontend-format: .frontend-prerequisites
	cd frontend/ && pnpm exec prettier --write .

.backend-format: .backend-prerequisites .backend-test-mocks
	go tool fieldalignment -fix \
		./api/... \
		./application/... \
		./certificate/commons/... \
		./certificate/custom/... \
		./certificate/external/... \
		./certificate/letsencrypt/... \
		./certificate/selfsigned/... \
		./core/... \
		./database/... \
		./integration/docker/... \
		./integration/truenas/... \
		./vpn/netbird/... \
		./vpn/tailscale/...
	go tool golangci-lint run --fix \
		./api/... \
		./application/... \
		./certificate/commons/... \
		./certificate/custom/... \
		./certificate/external/... \
		./certificate/letsencrypt/... \
		./certificate/selfsigned/... \
		./core/... \
		./database/... \
		./integration/docker/... \
		./integration/truenas/... \
		./vpn/netbird/... \
		./vpn/tailscale/...

clean:
	@find api application certificate core database i18n integration vpn -type f -name "*.mock.go" -delete

.backend-test-mocks: .backend-prerequisites
	@echo "Generating mock files..."
	@find api application certificate core database i18n integration vpn -type f -name "*.go" \
		-not -name "*_test.go" \
		-exec sh -c 'grep -q "^type [a-zA-Z0-9_]* interface" "$$1" && echo "$$1"' _ {} \; | \
	while read -r file; do \
		dir=$$(dirname "$$file"); \
		base=$$(basename "$$file" .go); \
		mock_file="$$dir/$${base}.mock.go"; \
		if [ -f "$$mock_file" ]; then continue; fi; \
		package_name=$$(basename "$$dir"); \
		interfaces=$$(grep -oE "^type [a-zA-Z0-9_]+ interface" "$$file" | awk '{print $$2}'); \
		mock_names_flag=""; \
		for i in $$interfaces; do \
			mock_names_flag="$$mock_names_flag,$$i=Mocked$$i"; \
		done; \
		go tool go.uber.org/mock/mockgen \
			-source "$$file" \
			-package "$$package_name" \
			-destination "$$mock_file" \
			-mock_names "$${mock_names_flag#,}" \
			-self_package "$$(cd $$dir && go list)" || true; \
	done

.backend-test: .backend-test-mocks .generate-i18n-files
	go test -coverprofile=coverage.out -covermode=atomic \
		./api/... \
		./application/... \
		./certificate/commons/... \
		./certificate/custom/... \
		./certificate/external/... \
		./certificate/letsencrypt/... \
		./certificate/selfsigned/... \
		./core/... \
		./database/... \
		./integration/docker/... \
		./integration/truenas/... \
		./vpn/netbird/... \
		./vpn/tailscale/...

.update-nginx-docker-image:
	@LATEST=$$(curl -fsSL https://raw.githubusercontent.com/docker-library/official-images/master/library/nginx | \
		awk '/^Tags:/{tags=$$0} /Directory: mainline\/alpine/{print tags; exit}' | \
		grep -oE '[0-9]+\.[0-9]+\.[0-9]+-alpine' | head -1); \
	CURRENT=$$(grep -E '^FROM nginx:' Dockerfile | grep -oE '[0-9]+\.[0-9]+\.[0-9]+-alpine' || true); \
	if [ "$$CURRENT" = "$$LATEST" ]; then \
		echo "nginx Docker image is up to date ($$LATEST)"; \
	elif [ -z "$$CURRENT" ]; then \
		sed -i'' -e "s|^FROM nginx:.* AS builder|FROM nginx:$$LATEST AS builder|" Dockerfile; \
		echo "Pinned nginx Docker image to $$LATEST"; \
	elif [ "$$(printf '%s\n' "$$CURRENT" "$$LATEST" | sort -V | tail -1)" = "$$LATEST" ] && [ "$$CURRENT" != "$$LATEST" ]; then \
		sed -i'' -e "s|^FROM nginx:.* AS builder|FROM nginx:$$LATEST AS builder|" Dockerfile; \
		echo "Updated nginx Docker image: $$CURRENT -> $$LATEST"; \
	else \
		echo "nginx Docker image ($$CURRENT) is newer than upstream ($$LATEST), skipping"; \
	fi

update-dependencies: .backend-prerequisites .frontend-prerequisites .update-nginx-docker-image
	cd api && go get -u ./...
	cd application && go get -u ./...
	cd certificate/commons && go get -u ./...
	cd certificate/custom && go get -u ./...
	cd certificate/external && go get -u ./...
	cd certificate/letsencrypt && go get -u ./...
	cd certificate/selfsigned && go get -u ./...
	cd core && go get -u ./...
	cd database && go get -u ./...
	cd integration/docker && go get -u ./...
	cd integration/truenas && go get -u ./...
	cd tools && go get -u ./...
	cd i18n && go get -u ./...
	cd vpn/netbird && go get -u ./...
	cd vpn/tailscale && go get -u ./...
	go work sync
	cd frontend && pnpm update

lint: .frontend-lint .backend-lint

format: .frontend-format .backend-format

test: .backend-prerequisites .backend-test

build-release: .frontend-build .backend-build .build-release-docker-image .build-distribution-files

build-snapshot:
	$(MAKE) .frontend-build .backend-build VERSION=0.0.0
	$(MAKE) .build-snapshot-docker-image VERSION=$(VERSION)
