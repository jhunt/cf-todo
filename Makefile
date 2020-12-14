build:
	cd ux && yarn build
	./assets > assets.go
	go fmt ./assets.go
