package config

import (
	"fmt"
	"strconv"

	ini "gopkg.in/ini.v1"
)

// Save func
func (r Parameters) Save(file string) error {
	fmt.Println(file)
	cfg := ini.Empty()
	sect := cfg.Section("service")
	sect.Key("debug_port").SetValue(strconv.Itoa(int(r.Service.DebugPort)))
	sect.Key("grpc_port").SetValue(strconv.Itoa(int(r.Service.GrpcPort)))
	sect.Key("http_port").SetValue(strconv.Itoa(int(r.Service.HTTPPort)))
	sect.Key("http_addr").SetValue(r.Service.HTTPAddr)
	sect.Key("consul_port").SetValue(strconv.Itoa(int(r.Service.ConsulPort)))
	sect.Key("consul_addr").SetValue(r.Service.ConsulAddr)
	sect = cfg.Section("DB")
	sect.Key("database").SetValue(r.DB.DB)
	sect.Key("db_user").SetValue(r.DB.DbUser)
	sect.Key("db_password").SetValue(r.DB.DbPassword)
	cfg.SaveTo(file)
	return nil
}

// Read func
func (r *Parameters) Read(file string) error {
	if r.Service.DebugPort == 0 {
		r.Service.DebugPort = 9100
	}
	if r.Service.ConsulPort == 0 {
		r.Service.ConsulPort = 8500
	}
	if r.Service.HTTPPort == 0 {
		r.Service.HTTPPort = 9110
	}
	if r.Service.GrpcPort == 0 {
		r.Service.GrpcPort = 9120
	}
	cfg, _ := ini.LooseLoad(file)
	return cfg.MapTo(&r)
}

// Config interface
type Config interface {
	Read(file string) error
	Save(file string) error
}

// Parameters struct for store service params
type Parameters struct {
	Service Service `ini:"service,omitempty"`
	DB      DB      `ini:"DB,omitempty"`
}

// Service struct
type Service struct {
	DebugPort  uint16 `ini:"debug_port,omitempty"`
	GrpcPort   uint16 `ini:"grpc_port,omitempty"`
	HTTPAddr   string `ini:"http_addr,omitempty"`
	HTTPPort   uint16 `ini:"http_port,omitempty"`
	ConsulAddr string `ini:"consul_addr,omitempty"`
	ConsulPort uint16 `ini:"consul_port,omitempty"`
}

// DB struct
type DB struct {
	DB         string `ini:"database,omitempty"`
	DbUser     string `ini:"db_user,omitempty"`
	DbPassword string `ini:"db_password,omitempty"`
}
