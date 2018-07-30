package k8s


import (
	"time"
	"fmt"

	"github.com/hashwing/log"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/runtime"
)


// WatchConfigMaps watch configmaps change
func (c *Client)WatchConfigMaps(label,value string,cf func(del bool, cm map[string]string)){
	// Building list watcher
	configmapListWatcher :=cache.NewListWatchFromClient(c.cset.CoreV1().RESTClient(),"configmaps",apiv1.NamespaceAll,fields.Everything())
	log.Debug("Building configmaps queue")
	// Building queue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	indexer, informer := cache.NewIndexerInformer(configmapListWatcher ,&apiv1.ConfigMap{},0,cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	}, cache.Indexers{})


	// Let the workers stop when we are done
	defer runtime.HandleCrash()
	defer queue.ShutDown()

	stop := make(chan struct{})
	defer close(stop)
	// Starting controller
	log.Debug("Starting configmaps controller")
	go informer.Run(stop)

	log.Debug("Wait for all involved caches to be synced")
	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stop,informer.HasSynced){
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	// configMaps save all configMaps
	configMaps:=make(map[string]interface{})

	callfunc :=func(del bool,item interface{}){
		if del{
			key :=item.(string)
			log.Info("configmap %s was deleted",key)
			if configMaps[key]!=nil{
				cm:=configMaps[key].(*apiv1.ConfigMap)
				cf(true,cm.Data)
			}
			return
		}
		cm,ok:=item.(*apiv1.ConfigMap)
		if ok{
			if cm.GetLabels()[label]==value{
				log.Info("configmap %s was update or add",cm.GetName())
				cf(false,cm.Data)
			}
		}
	}
	ctl:=&Controller{
		queue:queue,
		indexer:indexer,
		informer:informer,
		clientset:c.cset,
		cf:callfunc,
	}
	log.Debug("watching configmaps.....")
	wait.Until(ctl.work,time.Second,stop)
	
}