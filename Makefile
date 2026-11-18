DOCKER_IMAGE ?= dillmann/nginx-ignition
VERSION ?= 0.0.0
PR_ID ?= 0
SNAPSHOT_TAG_SUFFIX := $(if $(filter-out 0,$(PR_ID)),pr-$(PR_ID)-snapshot,$(VERSION)-snapshot)
LDFLAGS := -X 'github.com/lucasdillmann/nginx-ignition/internal/business/core/version.Number=$(VERSION)'
PNPM_INSTALL_FLAGS := $(if $(filter true,$(CI)),--frozen-lockfile,)
NSIS ?= makensis
NSIS_COMPRESS ?= lzma

.backend-prerequisites:
ifeq ($(CI),true)
	go mod download
	cd tools && go mod download
	cd tools/i18n && go mod download
else
	go mod tidy
	cd tools/i18n && go mod tidy
endif

.frontend-prerequisites:
	cd frontend/ && pnpm install $(PNPM_INSTALL_FLAGS)

.frontend-lint: .frontend-prerequisites
	cd frontend/ && pnpm run check

.backend-lint: .backend-prerequisites .backend-test-mocks
	go tool -modfile=tools/go.mod golangci-lint run ./internal/...

.frontend-build: .frontend-prerequisites .generate-i18n-files
	cd frontend/ && pnpm run build
	@find internal/api/domain/frontend/static -mindepth 1 -maxdepth 1 ! -name '.gitkeep' -exec rm -rf {} +
	cp -Rf frontend/build/. internal/api/domain/frontend/static/

.capture-prerequisites:
	cd tools/capture && pnpm install $(PNPM_INSTALL_FLAGS) && pnpm exec playwright install chromium

.docs-screenshots: .capture-prerequisites .frontend-build
	cd tools/capture && node ./src/capture.mjs screenshots

.docs-videos: .capture-prerequisites .frontend-build
	cd tools/capture && node ./src/capture.mjs videos

docs-assets: .docs-screenshots .docs-videos

.backend-build: .backend-prerequisites .generate-i18n-files .frontend-build
	$(MAKE) .backend-build-file OS=linux ARCH=amd64 DIR=linux
	$(MAKE) .backend-build-file OS=linux ARCH=arm64 DIR=linux
	$(MAKE) .backend-build-file OS=darwin ARCH=arm64 DIR=macos
	$(MAKE) .backend-build-file OS=windows ARCH=amd64 DIR=windows EXT=.exe
	$(MAKE) .backend-build-file OS=windows ARCH=arm64 DIR=windows EXT=.exe

.backend-build-file:
	@if [ "$(OS)" = "windows" ]; then $(MAKE) .windows-versioninfo VERSION=$(VERSION); fi
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o build/$(DIR)/$(ARCH)$(EXT) ./cmd

.windows-versioninfo:
	@rm -f cmd/resource_windows_amd64.syso cmd/resource_windows_arm64.syso
	@go tool -modfile=tools/go.mod goversioninfo -propagate-ver-strings=true -64 -arm=false \
		-file-version="$(VERSION)" -product-version="$(VERSION)" \
		-o cmd/resource_windows_amd64.syso tools/versioninfo.json
	@go tool -modfile=tools/go.mod goversioninfo -propagate-ver-strings=true -64 -arm=true \
		-file-version="$(VERSION)" -product-version="$(VERSION)" \
		-o cmd/resource_windows_arm64.syso tools/versioninfo.json

.generate-i18n-files:
	cd tools/i18n && go run .

.build-release-docker-image:
	docker buildx build \
		--progress=plain \
		--cache-from type=registry,ref=$(DOCKER_IMAGE):buildcache \
		--cache-to type=registry,ref=$(DOCKER_IMAGE):buildcache,mode=max \
		--tag $(DOCKER_IMAGE):$(VERSION) \
		--tag $(DOCKER_IMAGE):latest \
		--platform linux/amd64,linux/arm64 \
		--push .

.build-snapshot-docker-image:
	docker buildx build \
		--progress=plain \
		--cache-from type=registry,ref=$(DOCKER_IMAGE):buildcache \
		--cache-to type=registry,ref=$(DOCKER_IMAGE):buildcache,mode=max \
		--tag $(DOCKER_IMAGE):$(SNAPSHOT_TAG_SUFFIX) \
		--platform linux/amd64,linux/arm64 \
		--push .

.build-distribution-files: .build-distribution-non-windows .build-distribution-windows

.build-distribution-non-windows:
	$(MAKE) .build-distribution-zip ARCH=amd64 OS=linux SERVICE_FILE_EXT=service
	$(MAKE) .build-distribution-zip ARCH=arm64 OS=linux SERVICE_FILE_EXT=service
	$(MAKE) .build-distribution-zip ARCH=arm64 OS=macos SERVICE_FILE_EXT=plist
	$(MAKE) .build-distribution-packages ARCH=amd64 OS=linux
	$(MAKE) .build-distribution-packages ARCH=arm64 OS=linux

.build-distribution-windows:
	$(MAKE) .build-distribution-zip ARCH=amd64 OS=windows BIN_EXT=.exe
	$(MAKE) .build-distribution-zip ARCH=arm64 OS=windows BIN_EXT=.exe
	$(MAKE) .build-windows-installer ARCH=amd64
	$(MAKE) .build-windows-installer ARCH=arm64

.build-distribution-zip:
	rm -Rf build/nginx-ignition.$(OS)-$(ARCH).zip
	mkdir -p build/zip
	cp LICENSE build/zip/LICENSE.txt
	cp dist/$(OS)/instructions.md build/zip/instructions.md
	cp dist/$(OS)/nginx-ignition.properties build/zip/
	[ -z "$(SERVICE_FILE_EXT)" ] || cp dist/$(OS)/nginx-ignition.$(SERVICE_FILE_EXT) build/zip/
	cp build/$(OS)/$(ARCH)$(BIN_EXT) build/zip/nginx-ignition$(BIN_EXT)
	cd build/zip && zip -q -r ../nginx-ignition-$(VERSION).$(OS)-$(ARCH).zip .
	rm -Rf build/zip

.build-windows-installer:
	$(NSIS) -V2 \
		-DARCH=$(ARCH) \
		-DVERSION=$(VERSION) \
		-DNSIS_COMPRESS=$(NSIS_COMPRESS) \
		-DBINARY=$(CURDIR)/build/windows/$(ARCH).exe \
		-DOUTFILE=$(CURDIR)/build/nginx-ignition-$(VERSION).windows-$(ARCH)-setup.exe \
		-DASSETSDIR=$(CURDIR)/dist/windows/assets \
		-DLANGDIR=$(CURDIR)/dist/windows/lang \
		-DSCRIPTDIR=$(CURDIR)/dist/windows/scripts \
		-DTEMPLATEDIR=$(CURDIR)/dist/windows/templates \
		-DROOT=$(CURDIR) \
		dist/windows/installer.nsi

.build-windows-installers:
	$(MAKE) .build-windows-installer ARCH=amd64
	$(MAKE) .build-windows-installer ARCH=arm64

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
	go tool -modfile=tools/go.mod fieldalignment -fix ./internal/...
	go tool -modfile=tools/go.mod golangci-lint run --fix ./internal/...

clean:
	@find internal -type f -name "*.mock.go" -delete
	@find internal/api/domain/frontend/static -mindepth 1 -maxdepth 1 ! -name '.gitkeep' -exec rm -rf {} +
	@rm -f cmd/resource_windows_*.syso

.backend-test-mocks: .backend-prerequisites
	@echo "Generating mock files..."
	@find internal -type f -name "*.go" \
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
		go tool -modfile=tools/go.mod go.uber.org/mock/mockgen \
			-source "$$file" \
			-package "$$package_name" \
			-destination "$$mock_file" \
			-mock_names "$${mock_names_flag#,}" \
			-self_package "$$(cd $$dir && go list)" || true; \
	done

.backend-test: .backend-test-mocks .generate-i18n-files
	go test -coverprofile=coverage.out -covermode=atomic ./internal/...

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
	go get -u ./internal/...
	go get -u github.com/netbirdio/netbird@latest
	go get -u tailscale.com@latest
	go mod tidy
	cd tools && go get -u ./... && go mod tidy
	cd tools/i18n && go get -u ./... && go mod tidy
	cd frontend && pnpm update
	cd tools/capture && pnpm update

lint: .frontend-lint .backend-lint

format: .frontend-format .backend-format

test: .backend-prerequisites .backend-test

build-release: .frontend-build .backend-build .build-release-docker-image .build-distribution-files

build-snapshot:
	$(MAKE) .frontend-build .backend-build VERSION=0.0.0
	$(MAKE) .build-windows-installers VERSION=0.0.0 NSIS_COMPRESS=zlib
	$(MAKE) .build-snapshot-docker-image VERSION=$(VERSION)
