package dkron

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"gopkg.in/resty.v1"
)

func resourceJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceJobCreate,
		Read:   resourceJobRead,
		Update: resourceJobUpdate,
		Delete: resourceJobDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Optional: false,
			},
			"schedule": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"timezone": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"owner_email": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"disabled": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				Optional: false,
			},
			"dkron_host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Optional: false,
			},
			"concurrency": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"executor": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Optional: false,
			},
			"command": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Optional: false,
			},
		},
	}
}

type Job struct {
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

type JobResponse struct {
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

func resourceJobCreate(d *schema.ResourceData, m interface{}) error {
	job := new(Job)
	jobresp := new(JobResponse)

	job.Name = d.Get("name").(string)
	job.Schedule = "@every 10s"
	job.Owner = d.Get("owner").(string)
	job.OwnerEmail = d.Get("owner_email").(string)
	job.Disabled = d.Get("disabled").(bool)
	job.Retries = 1
	job.Executor = "shell"
	job.ExecutorConfig.Command = d.Get("command").(string)
	dkronHost := d.Get("dkron_host").(string)

	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(job).
		Post(dkronHost)

	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(resp.Body(), &jobresp)

	d.SetId(jobresp.Name)

	return resourceJobRead(d, m)
}

func resourceJobRead(d *schema.ResourceData, m interface{}) error {

	return nil
}

func resourceJobUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceJobRead(d, m)
}

func resourceJobDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
