provider "dkron" {
  host = "http://192.168.58.235:9191"

  # version = "0.0.1"
}

resource "dkron_job" "job" {
    name = "delete-nat"
    owner = "wangzhihu"
    owner_email = "wangzhihu@lixiang.com"
    executor = "shell"
    executor_config = {
      command = "echo 12344"
    }
    timezone = "Asia/Shanghai"
    disabled = false
    #schedule = "@every 10s"
    schedule = "@at 2021-07-02T07:27:00Z"
    retries = 2
    concurrency = "forbid"
    tags = {
      role = "dkron:1"
    }
}
