all: centaurid-upgrade-env

centaurid-upgrade-env:
	docker build --platform linux/amd64 --no-cache --build-arg entrypoint=cosmovisor-entrypoint.sh --tag centauri/centaurid-upgrade-env .

.PHONY: all centaurid-upgrade-env
