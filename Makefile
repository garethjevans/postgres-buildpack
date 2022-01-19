SHELL = bash
.ONESHELL:
VERSION = 0.0.1
FORMAT = file
OS = linux

clean:
	rm -fr build

.PHONY: build
build: create-package build/postgres-buildpack build/package.toml
	pack buildpack package "garethjevans-postgres-buildpack-$(VERSION).cnb" --config ./build/package.toml --format "$(FORMAT)"

create-package:
	GO111MODULE=on go install github.com/paketo-buildpacks/libpak/cmd/create-package

.PHONY: build/postgres-buildpack
build/postgres-buildpack: create-package
	create-package --destination ./build/postgres-buildpack --version "0.0.1"

.PHONY: build/package.toml
build/package.toml:
	./scripts/create-package-config.sh ./build/package.toml ./postgres-buildpack "$(OS)"
	cat ./build/package.toml


test:
	go test ./...
