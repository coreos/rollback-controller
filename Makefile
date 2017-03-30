.PHONY: build
build: bin/rollback-controller

.PHONY: FORCE

bin/%: FORCE
	@go build -i ./cmd/$*
	@go build -o ./bin/$* ./cmd/$*

.PHONY: vendor
vendor:
	@glide up --strip-vendor
	@glide-vc --no-tests --only-code

.PHONY: clean
clean:
	rm -rf bin
