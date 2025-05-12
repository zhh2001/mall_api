package main

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

func Register(address string, port int, name string, tags []string, id string, isCheck bool) error {
	cfg := api.DefaultConfig()
	cfg.Address = "10.120.21.77:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// 生成对应的检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           "http://10.120.221.149:8021/health",
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	registration.Address = address
	if isCheck {
		registration.Check = check
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	return nil
}

func AllServices() {
	cfg := api.DefaultConfig()
	cfg.Address = "10.120.21.77:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().Services()
	if err != nil {
		panic(err)
	}
	for key, value := range data {
		fmt.Println(key, value)
	}
}

func FilterServices() {
	cfg := api.DefaultConfig()
	cfg.Address = "10.120.21.77:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(`Service == "user-web"`)
	if err != nil {
		panic(err)
	}
	for key, value := range data {
		fmt.Println(key, value)
	}
}

func main() {
	_ = Register("10.120.221.149", 50052, "user-srv2", []string{"mall", "zhh"}, "user-srv2", false)
	AllServices()
	FilterServices()
}
