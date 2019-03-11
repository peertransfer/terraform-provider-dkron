FROM hashicorp/terraform:0.11.11 as terraform

FROM golang:1.11.2-alpine3.7 as base
ENV GOPATH /go
ENV WORKDIR $GOPATH/src/github.com/peertransfer/terraform_plugin_dkron
WORKDIR $WORKDIR

FROM base as deps
COPY . $WORKDIR
RUN apk add --no-cache git curl && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
CMD ["dep", "ensure", "-v"]

FROM base as dev
COPY --from=terraform /bin/terraform /bin/terraform
