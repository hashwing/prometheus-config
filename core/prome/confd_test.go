package prome_test

import (
	"testing"
	"github.com/hashwing/prometheus-config/core/prome"
)

func Test_CreateConfig(t *testing.T){
	job:=&prome.Job{
		Service:true,
		Node:true,
		PushGateway:true,
		Pod:true,
		Endpoints:true,
		ApiServers:true,
		Cadvisor:true,
	}
	c:=prome.Config{
		ConfigPath:"./config_test.yml",
		RulesPath:"/etc/prometheus/rules/*.yml",
		ScrapeInterval:"1m",
		ScrapeTimeout:"10s",
		EvaluationInterval:"1m",
		RemoteR: "http://127.0.0.1:8080/read",
		RemoteW: "http://127.0.0.1:8080/write",
		AlertManager:true,
		LabelKey:"prometheus_label",
		LabelValue:"prometheus",
		Job:job,
		ShardsNum:1,
		ShardsSum:3,
	}
	err:=prome.CreateConfig(c)
	if err!=nil{
		t.Error(err)
	}
}