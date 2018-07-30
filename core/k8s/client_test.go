package k8s_test

import (
	"testing"
	"github.com/hashwing/prometheus-config/core/k8s"
)


func Test_GetStatefulsetReplicas(t *testing.T){
	c,err:=k8s.New("10.21.21.170:9999","./config")
	if err!=nil{
		t.Error(err)
		return
	}
	r,err:=c.GetStatefulsetReplicas("test","prometheus")
	if err!=nil{
		t.Error(err)
		return
	}
	t.Log(r)
}
