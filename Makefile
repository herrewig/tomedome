export TOMEDOME_DB_FILEPATH=$(PWD)/internal/assets/mock_data.json
export RUN=docker-compose run --rm api

# Workflow stuff
test: test-js test-go

integration-test:
	$(RUN) gotest -v ./...

test-js:
	cd js; npm test

test-go:
	$(RUN) gotest -short -v ./...

run-server:
# DISABLE_VOLUME is a hack to prevent the volume from clobering the compose watch
# you can view the watch config in the docker-compose.yml file
	@DISABLE_VOLUME="/dev/null:/ignore" docker-compose up --watch

clean:
	rm -rf js/node_modules
	docker-compose down
	docker image rm api

###############
# Build stuff #
###############
build-dev:
	cd js; npm install
	docker-compose build

build-db:
	LOGLEVEL=error docker-compose run --rm api go run cmd/main.go --backend stratz --dump > go/internal/assets/data.json

build-image: build-db
	cd go ; docker build -t us-east4-docker.pkg.dev/tomedome/tomedome/api:latest .

#################
# Publish stuff #
#################
bounce-api:
	gcloud run deploy tomedome-api --image us-east4-docker.pkg.dev/tomedome/tomedome/api:latest --region us-east4

publish-image: build-image
	gcloud auth print-access-token | docker login -u oauth2accesstoken --password-stdin https://us-east4-docker.pkg.dev
	docker push us-east4-docker.pkg.dev/tomedome/tomedome/api:latest
	docker logout https://us-east4-docker.pkg.dev

publish-static:
	gsutil -h Cache-Control:"Cache-Control:private, max-age=5, no-transform" cp -r static/* gs://tomedome-static-site/
	gsutil -h Cache-Control:"Cache-Control:private, max-age=5, no-transform" rsync -x 'node_modules.*' -x '.*\.test\..*' js/src/ gs://tomedome-static-site/js/
