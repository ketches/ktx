.PHONY:
release:
	@if [ -z "${VERSION}" ]; then \
		echo "VERSION is not set"; \
		exit 1; \
	fi
	@if git rev-parse "refs/tags/${VERSION}" >/dev/null 2>&1; then \
        echo "Git tag ${VERSION} already exists, please use a new version."; \
        exit 1; \
    fi
	@sed -E -i '' 's/(const VERSION = ")[^"]+(")/\1${VERSION}\2/' cmd/version.go
	@git add cmd/version.go
	@git commit -m "Release ${VERSION}"
	@git push
	@git tag -a "${VERSION}" -m "release ${VERSION}"
	@git push origin "${VERSION}"
