[![CircleCI](https://circleci.com/gh/peertransfer/terraform-provider-dkron.svg?style=svg&circle-token=ad4b655899e45d9726cacc4d85f2e02a86147b40)](https://circleci.com/gh/peertransfer/terraform-provider-dkron)


# Terraform provider Dkron

## Install the plugin

Download the binary from the release and copy it in your project `~/.terraform.d/plugins` and run `terraform init`

## Interface

```
resource "dkron_job" "my-job" {
    name = "hola_from_tf"
    owner = "omar"
    owner_email = "a@a.com"
    dkron_host = "http://dkron:8080"
    executor = "shell"
    command = "date"
    disabled = false
}
```

## Steps to start developing

First of all you should download all dependencies:
```shell
$ make build_deps
```

And then you can start developing your Terraform plugin by:
```shell
$ make dev
```

When you're done, you can test your plugin by running:

```shell
$ make compile_plugin
```

```shell
$ make init
```

```shell
$ make plan
```

```shell
$ make apply
```

# TODO

- [ ] Use Dkron client instead to hardcode it into the provider
- [ ] Decouple provider from resource
- [ ] Add unit testing
- [ ] Implement destroy
- [ ] Complete API endpoints

