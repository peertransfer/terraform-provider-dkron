FROM golang:1.12.0-alpine3.9 as base
ENV WORKDIR /go/src/github.com/peertransfer/terraform-provider-dkron
ENV GO111MODULE on
WORKDIR $WORKDIR

FROM base as dev
COPY --from=hashicorp/terraform:0.11.11 /bin/terraform /bin/terraform
