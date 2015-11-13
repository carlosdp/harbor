package dockerscheduler

import (
	"os"
	"strings"

	"github.com/carlosdp/harbor/options"
	"github.com/carlosdp/harbor/scheduler"
)

func init() {
	scheduler.RegisterScheduler("nomad-scheduler", &dockerScheduler{})
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

type networks struct {
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
	Networks networks
	Meta     map[string]string
}

type task struct {
	Name        string
	Driver      string
	Config      map[string]string
	Constraints []constraint
	Resources   resources
}

type taskGroup struct {
	Name        string
	Count       int
	Constraints []constraint
	Tasks       []task
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

func (ds *nomadScheduler) New() scheduler.Scheduler {
	return &nomadScheduler{}
}

func (ds *nomadScheduler) Schedule(image, name, id string, ops options.Options) (interface{}, error) {
	name = strings.Replace(name, "/", "-", -1)
	cID, err := createContainer(image, name+"-"+id)
	return cID, err
}

func (ds *nomadScheduler) Rollback(name, id string, ops options.Options, state options.Option) error {
	cID := state.GetString()
}
