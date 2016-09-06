build: update-build-number
	go build -o tmp/build/sentry .
update-build-number:
	@echo updating version.go
	@perl -pi -e 's{__BUILD_NUM__ = (\d+)}{$$n=$$1+1; "__BUILD_NUM__ = $$n"}e' version.go
