package dkron

import (
	"encoding/json"
	"fmt"
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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Optional: false,
			},
			"displayname": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"schedule": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"retries": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
			},
			"processors": {
				Type:        schema.TypeMap,
				Description: "Billing information of the Rds.",
				Required:    false,
				Optional:    true,
				Elem:        tagsSchema(),
			},
			"owner_email": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Required: true,
				Optional: false,
			},
			"concurrency": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"executor": {
				Type:     schema.TypeString,
				Required: true,
				Optional: false,
			},
			"executor_config": tagsSchema(),
			"tags":            tagsSchema(),
			"metadata":        tagsSchema(),
		},
	}
}
func tagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeMap,
		Description: "Tags, do not support modify",
		Optional:    true,
		ForceNew:    false,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func tagsComputedSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeMap,
		Description: "Tags",
		Computed:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

type Job struct {
	Name           string                 `json:"name"`
	Displayname    string                 `json:"displayname"`
	Timezone       string                 `json:"timezone"`
	Schedule       string                 `json:"schedule"`
	Owner          string                 `json:"owner"`
	OwnerEmail     string                 `json:"owner_email"`
	Disabled       bool                   `json:"disabled"`
	Tags           map[string]interface{} `json:"tags"`
	Retries        int                    `json:"retries"`
	Metadata       map[string]interface{} `json:"metadata"`
	Processors     map[string]interface{} `json:"processors"`
	Concurrency    string                 `json:"concurrency"`
	Executor       string                 `json:"executor"`
	ExecutorConfig map[string]interface{} `json:"executor_config"`
}

type JobResponse struct {
	Name           string                 `json:"name"`
	Displayname    string                 `json:"displayname"`
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
	Metadata       map[string]interface{} `json:"metadata"`
	Processors     map[string]interface{} `json:"processors"`
	Concurrency    string                 `json:"concurrency"`
	Executor       string                 `json:"executor"`
	ExecutorConfig map[string]interface{} `json:"executor_config"`
	Status         string                 `json:"status"`
}

func resourceJobCreate(d *schema.ResourceData, m interface{}) error {
	return createJobData(d, m)
}

func resourceJobRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(Config)
	dkronHost := config.Host
	job := fillRequest(d, meta)
	jobresp := new(JobResponse)
	id := d.Id()
	jobsEndpoint := fmt.Sprintf("%s/v1/jobs/%s", dkronHost, id)
	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(job).Get(jobsEndpoint)
	if err != nil {
		return err
	}
	log.Printf("response body:%s", string(resp.Body()))
	json.Unmarshal(resp.Body(), &jobresp)
	d.SetId(jobresp.Name)
	d.Set("displayname", jobresp.Displayname)
	d.Set("concurrency", jobresp.Concurrency)
	d.Set("disabled", jobresp.Disabled)
	d.Set("executor", jobresp.Executor)
	d.Set("executor_config", jobresp.ExecutorConfig)
	d.Set("metadata", jobresp.Metadata)
	d.Set("owner", jobresp.Owner)
	d.Set("owner_email", jobresp.OwnerEmail)
	d.Set("parent_job", jobresp.ParentJob)
	d.Set("processors", jobresp.Processors)
	d.Set("retries", jobresp.Retries)
	d.Set("schedule", jobresp.Schedule)
	d.Set("tags", jobresp.Tags)
	d.Set("timezone", jobresp.Timezone)
	return nil
}

func resourceJobUpdate(d *schema.ResourceData, m interface{}) error {
	return createJobData(d, m)
}

func resourceJobDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(Config)
	jobID := d.Id()
	dkronHost := config.Host
	endpoint := fmt.Sprintf("%s/v1/jobs/%s", dkronHost, jobID)

	_, err := resty.R().
		SetHeader("Content-Type", "application/json").
		Delete(endpoint)

	if err != nil {
		return err
	}

	return nil
}
func fillRequest(d *schema.ResourceData, meta interface{}) *Job {
	job := new(Job)

	job.Name = d.Get("name").(string)
	job.Displayname = d.Get("displayname").(string)
	job.Schedule = d.Get("schedule").(string)
	job.Owner = d.Get("owner").(string)
	job.OwnerEmail = d.Get("owner_email").(string)
	job.Disabled = d.Get("disabled").(bool)
	job.Retries = d.Get("retries").(int)
	job.Executor = d.Get("executor").(string)
	if v, ok := d.GetOk("executor_config"); ok {
		executor_config := v.(map[string]interface{})
		job.ExecutorConfig = executor_config
	}
	if v, ok := d.GetOk("metadata"); ok {
		job.Metadata = v.(map[string]interface{})
	}
	job.Concurrency = d.Get("concurrency").(string)

	if v, ok := d.GetOk("processors"); ok {
		job.Processors = v.(map[string]interface{})
	}
	job.Tags = d.Get("tags").(map[string]interface{})
	return job
}

func createJobData(d *schema.ResourceData, meta interface{}) error {
	config := meta.(Config)
	dkronHost := config.Host
	job := fillRequest(d, meta)
	jobresp := new(JobResponse)

	jobsEndpoint := fmt.Sprintf("%s/v1/jobs", dkronHost)
	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(job).
		Post(jobsEndpoint)

	if err != nil {
		return err
	}
	log.Printf("response body:%s", string(resp.Body()))
	json.Unmarshal(resp.Body(), &jobresp)
	d.SetId(jobresp.Name)
	return resourceJobRead(d, meta)
}
