package main

import (
	"fmt"
	crdv1beta1 "github.com/cnych/controller-demo/pkg/apis/stable/v1beta1"
	informer "github.com/cnych/controller-demo/pkg/client/informers/externalversions/stable/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
	"time"
)

type Controller struct {
	informer informer.CronTabInformer
	queue    workqueue.RateLimitingInterface
}

func NewController(informer informer.CronTabInformer) *Controller {
	c := &Controller{
		informer: informer,
		queue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(),
			"crontab-controller"),
	}
	klog.Info("setting up crontab controller")
	//注册事件监听函数
	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAdd,
		UpdateFunc: c.onUpdate,
		DeleteFunc: c.onDelete,
	})
	return c
}

func (c *Controller) Run(stopCh <-chan struct{}) error { // stopCh 用来停止管道
	//函数结束以后关闭协程
	defer runtime.HandleCrash()

	//停止控制器需要关闭队列
	defer c.queue.ShutDown()

	//启动同样控制器框架
	klog.Infof("starting crontab controller")

	//等待所有相关的(indexer)缓存同步完成 然后在开始(workqueue)处理队列中的数据
	if !cache.WaitForCacheSync(stopCh, c.informer.Informer().HasSynced) {
		return fmt.Errorf("timed out waiting for caches to sync")
	}

	klog.Info("Informer caches to sync complated")

	//启动work处理队列中的数据
	go wait.Until(c.runWork, time.Second, stopCh)

	//这里stopCh读取管道 如果没有传递值，将会一直阻塞
	//上面协程就不会因为函数结束而结束 如果传递值了代表函数处理完成，打印退出日志
	<-stopCh
	klog.Info("stoppting crontab controller")
	return nil
}

func (c *Controller) runWork() {
	//处理队列中扽数据
	for c.processNextItem() {
	}

}

// 实现业务逻辑
func (c *Controller) processNextItem() bool {
	//从workqueue取出一个key
	obj, shutdown := c.queue.Get()
	//quit如果为true 队列已经关闭
	if shutdown {
		return false
	}

	//闭包
	err := func(obj interface{}) error {
		var ok bool
		var key string
		if key, ok = obj.(string); !ok {
			c.queue.Forget(obj)
			return fmt.Errorf("expected string in workqueue but get %#v\n", obj)
		}
		//业务逻辑处理
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("sync error: %v\n", err)
		}
		c.queue.Forget(obj)
		klog.Infof("Successfully synced %s\n", key)
		return nil
	}(obj)

	//函数退出的时候告诉队列已经处理了该key
	defer c.queue.Done(obj)
	//根据key去处理我们的业务逻辑
	//syncToStdout函数把从workqueue获取的key 打印出来 模拟业务逻辑

	if err != nil {
		runtime.HandleError(err)
	}

	return true
}

// key -> crontab -> indexer
func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	//获取crontab
	crontab, err := c.informer.Lister().CronTabs(namespace).Get(name)
	if err != nil {
		//对象已经被删除了
		if errors.IsNotFound(err) {
			klog.Warningf("Crontab deleting: %s/%s...", namespace, name)
			return nil
		}
		return err
	}
	klog.Infof("crontab try to process: %#v...", crontab)
	return nil
}

func (c *Controller) onAdd(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
	}
	c.queue.AddRateLimited(key)

}

func (c *Controller) onUpdate(old, new interface{}) {
	oldObj := old.(*crdv1beta1.CronTab)
	newObj := new.(*crdv1beta1.CronTab)
	//比较两个资源对象的资源版本是否一致
	if oldObj.ResourceVersion == newObj.ResourceVersion {
		return
	} else {
		//如果资源发生变动就把新的变动给onAdd
		c.onAdd(newObj)
	}
}

func (c *Controller) onDelete(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	} else {
		c.onAdd(key)
	}
}
