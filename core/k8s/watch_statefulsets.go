package k8s


import (
	"time"
	"fmt"

	"github.com/hashwing/log"
	apiv1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/runtime"
)


// WatchStatefulsets watch statfulsets change
func (c *Client)WatchStatefulsets(label,value string,cf func(del bool,obj *appsv1.StatefulSet)){
	// Building list watcher
	statefulsetsListWatcher:=cache.NewListWatchFromClient(c.cset.AppsV1().RESTClient(),"statefulsets",apiv1.NamespaceAll,fields.Everything())
	log.Debug("Building queue")
	// Building queue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	indexer, informer := cache.NewIndexerInformer(statefulsetsListWatcher,&appsv1.StatefulSet{},0,cache.ResourceEventHandlerFuncs{
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
	log.Debug("Starting controller")
	go informer.Run(stop)

	log.Debug("Wait for all involved caches to be synced")
	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stop,informer.HasSynced){
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	// pods save all pods
	statfulsets:=make(map[string]*appsv1.StatefulSet)

	callfunc :=func(del bool,item interface{}){
		if del{
			key :=item.(string)
			log.Debug("Pod", key,"does not exist anymore")
			s:=statfulsets[key]
			if s!=nil{
				delete(statfulsets,key)
				cf(true,s)
			}
			return
		}
		statefulset,ok:=item.(*appsv1.StatefulSet)
		if ok{
			if statefulset.GetLabels()[label]==value{
				statfulsets[statefulset.GetName()]=statefulset
				cf(false,statefulset)	
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
	log.Debug("watching.....")
	wait.Until(ctl.work,time.Second,stop)
}