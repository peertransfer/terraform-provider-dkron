VERSION = $(shell cat .semver)

dev:
	docker-compose up --build -d dev
	docker-compose up --build -d dkron

shell:
	docker-compose exec dev sh

workdir = /go/src/github.com/peertransfer/terraform-provider-dkron
plugin_dir = $(workdir)/terraform/plugins
version := $(shell cat .semver)

compile_plugin:
	docker-compose exec dev sh -c  "mkdir -p $(plugin_dir) && go build -mod=vendor -o $(plugin_dir)/terraform-provider-dkron_v$(version)"

init:
	docker-compose exec dev sh -c "cd terraform && terraform init -plugin-dir $(plugin_dir)"

plan:
	docker-compose exec dev sh -c "cd terraform && terraform plan"

apply:
	docker-compose exec dev sh -c "cd terraform && terraform apply"

destroy:
	docker-compose exec dev sh -c "cd terraform && terraform destroy"

build:
	CGO_ENABLED=0 go build -mod=vendor -o ./terraform-provider-dkron_v$(VERSION)

release: build
	go get github.com/tcnksm/ghr
	ghr -t ${GITHUB_TOKEN} -u ${CIRCLE_PROJECT_USERNAME} -r ${CIRCLE_PROJECT_REPONAME} -c ${CIRCLE_SHA1} v$(VERSION) ./terraform-provider-dkron_v$(VERSION)
