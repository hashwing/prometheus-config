package prome

import (
	"fmt"
	"os"
	"io/ioutil"
	"github.com/hashwing/log"
)

type Rule struct{
	Path	string
	Data	map[string]string
}

// CreateRules create rules
func (r *Rule)CreateRules()error{
	for k,v:=range r.Data{
		log.Info("add rule file:",k)
		err:=ioutil.WriteFile(fmt.Sprintf("%s/%s",r.Path,k),[]byte(v),0666)
		if err!=nil{
			return err
		}
	}
	return nil
}

// DeleteRules delete rules
func (r *Rule)DeleteRules()error{
	for k:=range r.Data{
		log.Info("delete rule file:",k)
		err:=os.Remove(fmt.Sprintf("%s/%s",r.Path,k))
		if err!=nil{
			return fmt.Errorf("delete %s error: %v",k,err)
		}
	}
	return nil
}