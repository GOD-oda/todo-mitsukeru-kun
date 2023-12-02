.PHONY: build
build:
	docker compose build --no-cache

.PHONY: deploy
deploy:
	git tag -d v0.0.1 | true
	git push -d origin v0.0.1 | true
	git tag -a -m "Description of this release" v0.0.1 | true
	git push origin v0.0.1