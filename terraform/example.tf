provider "dkron" {
  host = "http://dkron:8080"

  # version = "0.0.1"
}

resource "dkron_job" "my-job" {
    name = "hola_from_tf"
    owner = "omar"
    owner_email = "a@a.com"
    dkron_host = "http://dkron:8080/v1/jobs"
    executor = "shell"
    command = "date"
    disabled = false
}
