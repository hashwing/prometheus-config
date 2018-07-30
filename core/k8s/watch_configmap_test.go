package k8s_test

import (
	"testing"

	"github.com/hashwing/log"
	"github.com/hashwing/prometheus-config/core/k8s"
	"github.com/hashwing/prometheus-config/core/prome"
)


func Test_WatchConfigMaps(t *testing.T){
	c,err:=k8s.New("10.21.21.170:9999","/root/go/src/github.com/hashwing/prometheus-config/config")
	if err!=nil{
		t.Error(err)
		return
	}

	c.WatchConfigMaps("app","grafana",func(del bool,data map[string]string){
		for k:=range data{
			log.Info(k)
		}
		rule:=prome.Rule{
			Path:".",
			Data:data,
		}
		err:=rule.CreateRules()
		if err!=nil{
			log.Error(err)
			t.Error(err)
		}
	})
}
