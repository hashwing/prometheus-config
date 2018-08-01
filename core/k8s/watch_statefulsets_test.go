package k8s_test

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	"github.com/hashwing/log"
	"github.com/hashwing/prometheus-config/core/k8s"
)


func Test_WatchStatefulsets(t *testing.T){
	c,err:=k8s.New("10.21.21.170:9999","/root/go/src/github.com/hashwing/prometheus-config/config")
	if err!=nil{
		t.Error(err)
		return
	}

	c.WatchStatefulsets("prometheus_monitor","kubernetes",func(del bool, s *appsv1.StatefulSet){
		log.Info("statefulset %s replicas is %v",s.Name,*(s.Spec.Replicas))
	})
}
