package consulnotifier

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/carlosdp/harbor/notifier"
	"github.com/carlosdp/harbor/options"
)

type consulNotifier struct {
}

type service struct {
	ID      string
	Name    string
	Tags    []string
	Address string
	Port    string
	Check   map[string]string
}

func init() {
	notifier.RegisterNotifier("consul-notifier", &consulNotifier{})
}

func (c *consulNotifier) New() notifier.Notifier {
	return &consulNotifier{}
}

func (c *consulNotifier) Notify(name, id string, opts options.Options) (interface{}, error) {
	host := opts.GetString("host")
	if host == "" {
		host = "localhost"
	}

	port := opts.GetString("port")
	if port == "" {
		port = "8500"
	}

	serviceName := opts.GetString("service")
	if serviceName == "" {
		serviceName = name
	}

	servicePort := opts.GetString("service_port")
	if servicePort == "" {
		return nil, errors.New("service_port must be defined")
	}

	s := service{
		ID:      serviceName + id,
		Name:    serviceName,
		Address: "127.0.0.1",
		Port:    servicePort,
	}

	buf, _ := json.Marshal(s)

	client := &http.Client{}

	req, err := http.NewRequest("PUT", "http://"+host+":"+port+"/v1/agent/service/register", bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	_, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *consulNotifier) Rollback(name, id string, opts options.Options, state options.Option) error {
	return nil
}
