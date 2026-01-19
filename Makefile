DOCKER_IMAGE ?= dillmann/nginx-ignition
VERSION ?= 0.0.0
PR_ID ?= 0
SNAPSHOT_TAG_SUFFIX := $(if $(filter-out ,$(PR_ID)),$(if $(filter-out 0,$(PR_ID)),pr-$(PR_ID)-snapshot,snapshot),snapshot)
LDFLAGS := -X 'dillmann.com.br/nginx-ignition/core/common/version.Number=$(VERSION)'

.backend-prerequisites:
	go work sync

.frontend-prerequisites:
	cd frontend/ && npm ci

.frontend-lint: .frontend-prerequisites
	cd frontend/ && npm run check

.backend-lint: .backend-prerequisites
	go tool golangci-lint run \
		./api/... \
		./application/... \
		./certificate/commons/... \
		./certificate/custom/... \
		./certificate/letsencrypt/... \
		./certificate/selfsigned/... \
		./core/... \
		./database/... \
		./integration/docker/... \
		./integration/truenas/... \
		./vpn/tailscale/...

.frontend-build: .frontend-prerequisites .generate-i18n-files
	cd frontend/ && npm run build

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
	docker buildx build \
		--tag $(DOCKER_IMAGE):$(VERSION) \
		--tag $(DOCKER_IMAGE):latest \
		--platform linux/amd64,linux/arm64 \
		--push .

.build-snapshot-docker-image:
	docker buildx build \
		--tag $(DOCKER_IMAGE):$(SNAPSHOT_TAG_SUFFIX) \
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
	cd frontend/ && npx prettier --write .

.backend-format: .backend-prerequisites
	go tool fieldalignment -fix \
		./api/... \
		./application/... \
		./certificate/commons/... \
		./certificate/custom/... \
		./certificate/letsencrypt/... \
		./certificate/selfsigned/... \
		./core/... \
		./database/... \
		./integration/docker/... \
		./integration/truenas/... \
		./vpn/tailscale/...
	go tool golangci-lint run --fix \
		./api/... \
		./application/... \
		./certificate/commons/... \
		./certificate/custom/... \
		./certificate/letsencrypt/... \
		./certificate/selfsigned/... \
		./core/... \
		./database/... \
		./integration/docker/... \
		./integration/truenas/... \
		./vpn/tailscale/...

.backend-test-mocks: .backend-prerequisites
	@echo "Generating mock files..."
	@find api application certificate core database i18n integration vpn -type f -name "*.mock.go" -delete;
	@find api application certificate core database i18n integration vpn -type f -name "*.go" \
		-not -name "*_test.go" \
		-exec sh -c 'grep -q "^type [a-zA-Z0-9_]* interface" "$$1" && echo "$$1"' _ {} \; | \
	while read -r file; do \
		dir=$$(dirname "$$file"); \
		base=$$(basename "$$file" .go); \
		mock_file="$$dir/$${base}.mock.go"; \
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
	go test \
		./api/... \
		./application/... \
		./certificate/commons/... \
		./certificate/custom/... \
		./certificate/letsencrypt/... \
		./certificate/selfsigned/... \
		./core/... \
		./database/... \
		./integration/docker/... \
		./integration/truenas/... \
		./vpn/tailscale/...

update-dependencies: .backend-prerequisites .frontend-prerequisites
	cd api && go get -u all
	cd application && go get -u all
	cd certificate/commons && go get -u all
	cd certificate/custom && go get -u all
	cd certificate/letsencrypt && go get -u all
	cd certificate/selfsigned && go get -u all
	cd core && go get -u all
	cd database && go get -u all
	cd integration/docker && go get -u all
	cd integration/truenas && go get -u all
	cd tools && go get -u all
	cd vpn/tailscale && go get -u all
	go work sync
	cd frontend && npm update

lint: .frontend-lint .backend-lint

format: .frontend-format .backend-format

test: .backend-prerequisites .backend-test

build-release: .frontend-build .backend-build .build-release-docker-image .build-distribution-files

build-snapshot: .frontend-build .backend-build .build-snapshot-docker-image .build-distribution-files
