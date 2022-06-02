GOLANG_CROSS_VERSION  ?= v1.18.2
PACKAGE_NAME          := github.com/gobuffalo/pop

.PHONY: release
release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	docker run \
		--rm \
		--privileged \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/depbot \
		-v `pwd`/sysroot:/sysroot \
		-w /depbot \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --rm-dist