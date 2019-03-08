package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gopkg.in/resty.v1"
)

type dkronJob struct {
	Name       string `json:"name"`
	Schedule   string `json:"schedule"`
	Owner      string `json:"owner"`
	OwnerEmail string `json:"owner_email"`
	Disabled   bool   `json:"disabled"`
	Tags       struct {
	} `json:"tags"`
	Retries        int         `json:"retries"`
	Processors     interface{} `json:"processors"`
	Concurrency    string      `json:"concurrency"`
	Executor       string      `json:"executor"`
	ExecutorConfig struct {
		Command string `json:"command"`
	} `json:"executor_config"`
}

type dkronJobResponse struct {
	Name         string    `json:"name"`
	Timezone     string    `json:"timezone"`
	Schedule     string    `json:"schedule"`
	Owner        string    `json:"owner"`
	OwnerEmail   string    `json:"owner_email"`
	SuccessCount int       `json:"success_count"`
	ErrorCount   int       `json:"error_count"`
	LastSuccess  time.Time `json:"last_success"`
	LastError    time.Time `json:"last_error"`
	Disabled     bool      `json:"disabled"`
	Tags         struct {
	} `json:"tags"`
	Retries        int         `json:"retries"`
	DependentJobs  interface{} `json:"dependent_jobs"`
	ParentJob      string      `json:"parent_job"`
	Processors     interface{} `json:"processors"`
	Concurrency    string      `json:"concurrency"`
	Executor       string      `json:"executor"`
	ExecutorConfig struct {
		Command string `json:"command"`
	} `json:"executor_config"`
	Status string `json:"status"`
}

//Error handling
type AuthSuccess struct{}

func main() {
	job := new(dkronJob)
	jobresp := new(dkronJobResponse)
	job.Name = "hola"
	job.Schedule = "@every 10s"
	job.Owner = "paco"
	job.OwnerEmail = "paco@paco.com"
	job.Disabled = false
	job.Retries = 1
	job.Executor = "shell"
	job.ExecutorConfig.Command = "/bin/true"

	resp, err := resty.R().SetHeader("Content-Type", "application/json").SetBody(job).SetResult(&AuthSuccess{}).Post("http://dkron:8080/v1/jobs")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nResponse Body: %v", string(resp.Body()))
	json.Unmarshal(resp.Body(), &jobresp)
	fmt.Printf("\nResponse Body: %v", jobresp.Name)
	fmt.Printf("\nResponse Body: %v", resp.StatusCode())
}
