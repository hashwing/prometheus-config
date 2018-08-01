package main

import (
	"os"
	"strconv"
	"strings"
	"github.com/namsral/flag"
	"github.com/hashwing/log"
	"github.com/hashwing/prometheus-config/core/k8s"
	"github.com/hashwing/prometheus-config/core/prome"	
)

func main(){
	var (
		masterURL     		= flag.String("master_url", "", "Kubernetes apiserver URL")
		kubeconfig    		= flag.String("kube_config", "", "Kubernetes kubeconfig path")
		podName    			= flag.String("pod_name", "", "Kubernetes pod name")
		labelKey      		= flag.String("label_key", "prometheus_shards", "Watch label key")
		labelValue    		= flag.String("label_value", "", "Watch label value")
		configPath    		= flag.String("prome_config_path", "/etc/prometheus/config/prometheus.yml", "Prometheus configfile path")
		scrapeInterval    	= flag.String("prome_scrape_inv", "1m", "Prometheus scrape interval")
		scrapeTimeout    	= flag.String("prome_scrape_timeout", "10s", "Prometheus scrape timeout")
		evaluationInterval 	= flag.String("prome_evaluation_inv", "1m", "Prometheus evalustion interval")
		rulesDir   			= flag.String("prome_rule_dir", "/etc/prometheus/rules", "Prometheus rules dir")
		remoteR  			= flag.String("prome_remote_read", "", "Prometheus remote read url")
		remoteW 			= flag.String("prome_remote_write", "", "Prometheus remote write url")
		alertmanager 		= flag.Bool("prome_alertmanager", false, "Prometheus alertmanager enadble")
		roles  				= flag.String("prome_role_jobs", "local,service,endpoints,pod,cadvisor", "Prometheus role job:local,service,endpoints,pod,cadvisor,pushgateway,apiservers,node")
		first				= flag.Bool("first_init",true,"First init config or not")
	)
	flag.Parse()

	logger,_:=log.NewBeegoLog("",false,true,false)
	log.SetHlogger(logger)

	// get pod number
	names:=strings.Split(*podName,"-")
	number,err:=strconv.Atoi(names[len(names)-1])
	if err!=nil{
		log.Error("get pod number error:",err)
		os.Exit(1)
	}

	// new a prometheus config struct
	job:=&prome.Job{}
	transformRoles(job,*roles)
	config:=prome.Config{
		ConfigPath:*configPath,
		RulesPath:*rulesDir+"/*",
		ScrapeInterval:*scrapeInterval,
		EvaluationInterval:*evaluationInterval,
		ScrapeTimeout:*scrapeTimeout,
		RemoteR: *remoteR,
		RemoteW:*remoteW,
		AlertManager:*alertmanager,
		LabelKey:*labelKey,
		LabelValue:*labelValue,
		Job:job,
		ShardsNum:number,
		ShardsSum:number,
	}

	// if init config write a default config file
	if *first{
		err=prome.CreateConfig(config)
		if err!=nil{
			log.Error("write template config file error:",err)
			os.Exit(1)
		}
		return
	}

	// new a k8s clientset
	c,err:=k8s.New(*masterURL,*kubeconfig)
	if err!=nil{
		log.Error("new a clientset error:",err)
		os.Exit(1)
	}

	// watch pods
	go c.WatchPods(*labelKey,*labelValue,func(n,ns string){
		r,err:=c.GetStatefulsetReplicas(c.GetOwner(n),ns)
		if err!=nil{
			log.Error("get statefulset replicas error:",err)
			return
		}
		if config.ShardsSum==int(r){
			return
		}
		config.ShardsSum=int(r)
		err=prome.CreateConfig(config)
		if err!=nil{
			log.Error("write template config file error:",err)
		}
		err=prome.ReloadEndpoint("http://127.0.0.1:9090/-/reload/")
		if err!=nil{
			log.Error(err)
		}
	})

	// watch configmaps
	c.WatchConfigMaps(*labelKey,*labelValue,func(del bool,data map[string]string){
		rule:=&prome.Rule{
			Path:*rulesDir,
			Data:data,
		}
		if del{
			err:=rule.DeleteRules()
			if err!=nil{
				log.Error("delete rules  error:",err)
			}
		}else{
			err:=rule.CreateRules()
			if err!=nil{
				log.Error("delete rules  error:",err)
			}
		}
		
		err:=prome.ReloadEndpoint("http://127.0.0.1:9090/-/reload/")
		if err!=nil{
			log.Error(err)
		}
		
	})
}

// transformRoles transform string to job struct
func transformRoles(job *prome.Job,roles string){
	jobArr:=strings.Split(roles,",")
	for _,jobStr := range jobArr{
		switch jobStr{
			case "local":
				job.Local=true
				break
			case "node":
				job.Node=true
				break
			case "service":
				job.Service=true
				break
			case "pod":
				job.Pod=true
				break
			case "endpoints":
				job.Endpoints=true
				break
			case "apiserver":
				job.ApiServers=true
				break
			case "cadvisor":
				job.Cadvisor=true
				break
			case "pushgateway":
				job.PushGateway=true
				break
		}

	}
}