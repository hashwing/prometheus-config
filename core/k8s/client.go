package k8s


import (
	"strings"

	//appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/apimachinery/pkg/util/runtime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Client struct{
     cset *kubernetes.Clientset
}

// New create the clientset
func New(masterUrl,kubeconfig string)(*Client,error){
   // use the current context in kubeconfig
   config, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfig)
   if err != nil {
	  return nil,err
   }

   // create the clientset
   clientset, err := kubernetes.NewForConfig(config)
   if err!=nil{
		return nil,err
   }
    return &Client{cset:clientset},nil
}

// GetOwner get pod owner
func (c *Client)GetOwner(name string)string{
	return strings.Split(name,"-")[0]
}

// GetStatefulsetReplicas get statefulset replicas
func (c *Client)GetStatefulsetReplicas(name,namespace string)(int32,error){
	statefulset,err:=c.cset.AppsV1().StatefulSets(namespace).Get(name,metav1.GetOptions{})
	if err!=nil{
		return 0,err
	}

	return *(statefulset.Spec.Replicas),nil
}

// GetDeployReplicas get deploy
func (c *Client)GetDeployReplicas(name,namespace string)(int32,error){
	deploy,err:=c.cset.AppsV1().Deployments(namespace).Get(name,metav1.GetOptions{})
	if err!=nil{
		return 0,err
	}
	return *(deploy.Spec.Replicas),nil
}


type Controller struct {
	indexer  cache.Indexer
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
	clientset  *kubernetes.Clientset
	cf   func(delete bool,item interface{})
}

func (c *Controller)work(){
	for c.processNextItem(){

	}
}

func (c *Controller)processNextItem() bool {
	// Wait until there is a new item in the working queue
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	// Invoke the method containing the business logic
	err := c.syncToStdout(key.(string))
	// Handle the error if something went wrong during the execution of the business logic
	c.handleErr(err, key)
	return true
}


// syncToStdout is the business logic of the controller. In this controller it simply prints
// information about the pod to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (c *Controller) syncToStdout(key string) error {
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		return err
	}

	if !exists {
		// Below we will warm up our cache with a Pod, so that we will see a delete for one pod
		c.cf(true,key)
		
	} else {
		c.cf(false,obj)
	}
	return nil
}

// handleErr checks if an error happened and makes sure we will retry later.
func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.queue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.queue.NumRequeues(key) < 5 {

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
}