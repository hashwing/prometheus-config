package k8s_test

import (
	"testing"

	"github.com/hashwing/log"
	"github.com/hashwing/prometheus-config/core/k8s"
)


func Test_WatchPods(t *testing.T){
	c,err:=k8s.New("10.21.21.170:9999","/root/go/src/github.com/hashwing/prometheus-config/config")
	if err!=nil{
		t.Error(err)
		return
	}

	c.WatchPods("app","grafana",func(n,ns string){
		log.Info(n,ns)
		r,err:=c.GetDeployReplicas(c.GetOwner(n),ns)
		if err!=nil{
			t.Error(err)
		}
		log.Info(r)
	})
}
