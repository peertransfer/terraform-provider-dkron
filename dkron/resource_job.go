package dkron

import (
	"encoding/json"
	"fmt"
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
			"retries": &schema.Schema{
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
			},
			"slack": &schema.Schema{
				Type:     schema.TypeMap,
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
			"tags": &schema.Schema{
				Type:     schema.TypeMap,
				Required: false,
				Optional: true,
			},
		},
	}
}

type Processors struct {
	Slack map[string]interface{} `json:"slack"`
}

type Job struct {
	Name           string                 `json:"name"`
	Schedule       string                 `json:"schedule"`
	Owner          string                 `json:"owner"`
	OwnerEmail     string                 `json:"owner_email"`
	Disabled       bool                   `json:"disabled"`
	Tags           map[string]interface{} `json:"tags"`
	Retries        int                    `json:"retries"`
	Processors     Processors             `json:"processors"`
	Concurrency    string                 `json:"concurrency"`
	Executor       string                 `json:"executor"`
	ExecutorConfig struct {
		Command string `json:"command"`
	} `json:"executor_config"`
}

type JobResponse struct {
	Name           string                 `json:"name"`
	Timezone       string                 `json:"timezone"`
	Schedule       string                 `json:"schedule"`
	Owner          string                 `json:"owner"`
	OwnerEmail     string                 `json:"owner_email"`
	SuccessCount   int                    `json:"success_count"`
	ErrorCount     int                    `json:"error_count"`
	LastSuccess    time.Time              `json:"last_success"`
	LastError      time.Time              `json:"last_error"`
	Disabled       bool                   `json:"disabled"`
	Tags           map[string]interface{} `json:"tags"`
	Retries        int                    `json:"retries"`
	DependentJobs  interface{}            `json:"dependent_jobs"`
	ParentJob      string                 `json:"parent_job"`
	Processors     Processors             `json:"processors"`
	Concurrency    string                 `json:"concurrency"`
	Executor       string                 `json:"executor"`
	ExecutorConfig struct {
		Command string `json:"command"`
	} `json:"executor_config"`
	Status string `json:"status"`
}

func resourceJobCreate(d *schema.ResourceData, m interface{}) error {
	return postJobData(d, m)
}

func resourceJobRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceJobUpdate(d *schema.ResourceData, m interface{}) error {
	return postJobData(d, m)
}

func resourceJobDelete(d *schema.ResourceData, m interface{}) error {
	jobID := d.Id()
	dkronHost := d.Get("dkron_host").(string)
	endpoint := fmt.Sprintf("%s/v1/jobs/%s", dkronHost, jobID)

	_, err := resty.R().
		SetHeader("Content-Type", "application/json").
		Delete(endpoint)

	if err != nil {
		return err
	}

	return nil
}

func postJobData(d *schema.ResourceData, m interface{}) error {
	job := new(Job)
	jobresp := new(JobResponse)

	job.Name = d.Get("name").(string)
	job.Schedule = d.Get("schedule").(string)
	job.Owner = d.Get("owner").(string)
	job.OwnerEmail = d.Get("owner_email").(string)
	job.Disabled = d.Get("disabled").(bool)
	job.Retries = d.Get("retries").(int)
	job.Executor = d.Get("executor").(string)
	job.ExecutorConfig.Command = d.Get("command").(string)
	job.Concurrency = d.Get("concurrency").(string)
	job.Processors.Slack = d.Get("slack").(map[string]interface{})
	dkronHost := d.Get("dkron_host").(string)
	job.Tags = d.Get("tags").(map[string]interface{})

	jobsEndpoint := fmt.Sprintf("%s/v1/jobs", dkronHost)
	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(job).
		Post(jobsEndpoint)

	if err != nil {
		return err
	}

	json.Unmarshal(resp.Body(), &jobresp)

	d.SetId(jobresp.Name)

	return resourceJobRead(d, m)
}
