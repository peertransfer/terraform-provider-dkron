provider "dkron" {
  host = "http://dkron:8080"

  # version = "0.0.1"
}

resource "dkron_job" "my-job" {
    name = "lolaso"
    owner = "omar"
    owner_email = "a@a.com"
    dkron_host = "http://dkron:8080"
    executor = "shell"
    command = "date"
    disabled = false
    schedule = "@every 10s"
    retries = 2
    concurrency = "forbid"
}
