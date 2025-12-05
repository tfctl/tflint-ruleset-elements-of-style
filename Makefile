.PHONY: default build check clean install release test tflint
default: build

build:
	go build

clean-runs:
	gh run list --limit 100 --json databaseId,createdAt | \
		jq -r '.[] | select((.createdAt | fromdateiso8601) < (now - 3 * 24 * 60 * 60)) | .databaseId' | \
		xargs -n1 echo gh run delete

check:
	tools/check.sh --all

install: build
	mkdir -p ~/.tflint.d/plugins
	mv ./tflint-ruleset-elements-of-style ~/.tflint.d/plugins

release:
	@if [ -z "$(VERSION)" ]; then echo "Usage: make release VERSION=x.y.z"; exit 1; fi
	@if ! echo "$(VERSION)" | grep -qE '^[0-9]+\.[0-9]+\.[0-9]+$$'; then \
		echo "Error: VERSION must be a valid semantic version (e.g. 0.2.0) without leading 'v'. Got: $(VERSION)"; \
		exit 1; \
	fi

	sed -i 's/\(version = \|Version: \)"[0-9]\+\.[0-9]\+\.[0-9]\+"/\1"$(VERSION)"/' README.md main.go
	git add README.md main.go
	git diff --cached --quiet || git commit --no-verify --message "Release version $(VERSION)"
	git tag v$(VERSION) --message "Release version v$(VERSION)"
	@echo "Successfully bumped to $(VERSION) and created tag v$(VERSION)."
	@echo "Run 'git push origin main --tags' to publish."

test: build
	go test ./... --count 1 -v

tflint: install
	for d in rules/*/testdata; do tflint --chdir $$d; done


