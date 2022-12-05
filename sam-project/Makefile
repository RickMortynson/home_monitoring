build.dev:
	sam build --use-container

build.deploy:
	sam build --use-container --container-env-var GO_ENV=production

invoke:
	make build.dev && sam local invoke

deploy: 
	make build.deploy && sam deploy