build:
	cd ux && yarn build
	./assets > assets.go
	go fmt ./assets.go

docker:
	docker build -t filefrog/cf-todo:latest .
	docker push filefrog/cf-todo:latest
