VERSION:=0.0.4

.PHONY:=build
build:
	docker build -t sdorra/host-path-provisioner:${VERSION} .

.PHONY:=deploy
deploy:
	docker push sdorra/host-path-provisioner:${VERSION}
