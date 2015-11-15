package nomadscheduler

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"net/http"

	"github.com/carlosdp/supply-chain/options"
	"github.com/carlosdp/supply-chain/scheduler"
)

func init() {
	scheduler.RegisterScheduler("nomad-scheduler", &nomadScheduler{})
}

type nomadScheduler struct {
}

type constraint struct {
	Hard    bool
	LTarget string
	RTarget string
	Operand string
	Weight  int
}

type network struct {
	Device        string
	CIDR          string
	IP            string
	MBits         int
	ReserverPorts []string
	DynamicPorts  []int
}

type resources struct {
	CPU      int
	MemoryMB int
	DiskMB   int
	IOPS     int
	Networks []network
	Meta     map[string]string
}

type task struct {
	Name        string
	Driver      string
	Config      map[string]string
	Constraints []constraint
	Resources   resources
}

type restartPolicy struct {
	Attempts int
	Interval time.Duration
	Delay    time.Duration
}

type taskGroup struct {
	Name          string
	Count         int
	Constraints   []constraint
	RestartPolicy restartPolicy
	Tasks         []task
}

type job struct {
	Region      string
	ID          string
	Name        string
	Type        string
	Priority    int
	AllAtOnce   bool
	Datacenters []string
	Constraints []constraint
	TaskGroups  []taskGroup
	Update      map[string]string
	Meta        map[string]string
}

type nomadRequest struct {
	Job job
}

func (ns *nomadScheduler) New() scheduler.Scheduler {
	return &nomadScheduler{}
}

func (ns *nomadScheduler) Schedule(image, name, id string, opts options.Options) (interface{}, error) {
	name = strings.Replace(name, "/", "-", -1)
	jobName := name + "-" + id

	host := opts.GetString("host")
	if host == "" {
		host = "localhost"
	}

	port := opts.GetString("port")
	if port == "" {
		port = "4646"
	}

	j := defaultJob(jobName)
	tg := defaultTaskGroup(jobName)
	t := defaultTask(jobName)
	t.Config = map[string]string{
		"image": image,
	}
	tg.Tasks = append(tg.Tasks, t)
	j.TaskGroups = append(j.TaskGroups, tg)

	buf, _ := json.Marshal(nomadRequest{j})

	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://"+host+":"+port+"/v1/jobs", bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New("could not provision job with Nomad: " + resp.Status)
	}

	return nil, nil
}

func (ns *nomadScheduler) Rollback(name, id string, opts options.Options, state options.Option) error {
	name = strings.Replace(name, "/", "-", -1)

	host := opts.GetString("host")
	if host == "" {
		host = "localhost"
	}

	port := opts.GetString("port")
	if port == "" {
		port = "4646"
	}

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", "http://"+host+":"+port+"/v1/job/"+name+"-"+id, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New("could not rollback with Nomad: " + resp.Status)
	}

	return nil
}

func defaultJob(name string) job {
	return job{
		Region:      "global",
		ID:          name,
		Name:        name,
		Type:        "service",
		Priority:    50,
		Datacenters: []string{"dc1"},
		TaskGroups:  []taskGroup{},
	}
}

func defaultTaskGroup(name string) taskGroup {
	return taskGroup{
		Name:  name,
		Count: 1,
		RestartPolicy: restartPolicy{
			Attempts: 10,
			Interval: time.Duration(10) * time.Minute,
			Delay:    time.Duration(30) * time.Second,
		},
		Tasks: []task{},
	}
}

func defaultTask(name string) task {
	return task{
		Name:   name,
		Driver: "docker",
		Resources: resources{
			CPU:      500,
			MemoryMB: 256,
			Networks: []network{},
		},
	}
}
