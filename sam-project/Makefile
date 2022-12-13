build.dev:
	sam build --use-container

# TODO: update build script, use zip instead of image deployment
# TODO: after this done - use --container-env-var-file to pass env's 
build.deploy:
	sam build --use-container

invoke:
	make build.dev && sam local invoke

deploy: 
	make build.deploy && sam deploy	