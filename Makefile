VERSION=0.2
build:
	docker build -t mconf/mconf-api:nginx-exporter-v${VERSION} .