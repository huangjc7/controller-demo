package main

import (
	"fmt"
	"github.com/cnych/controller-demo/pkg/client/clientset/versioned"
	"github.com/cnych/controller-demo/pkg/client/informers/externalversions"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"log"
	"os"
	"time"
)

var (
	kubeConfigPath string
)

func main() {
	//生成clientset
	kubeConfigPath = fmt.Sprintf("%s%s", os.Getenv("HOME"), "/.kube/config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	//versioed是根据代码生成器生成出来的NewForConfig
	crontabClientset, err := versioned.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error init kubernetes crontab clientset %s\n", err.Error())
	}
	crontabInformerFactory := externalversions.NewSharedInformerFactory(crontabClientset, time.Second*30)

	stopCh := make(<-chan struct{})

	//实例化 CronTab控制器
	controller := NewController(crontabInformerFactory.Stable().V1beta1().CronTabs())

	// 启动informer 执行ListAndWatch操作
	//在informer启动之前需要先注册一下informer 在启动才可以
	go crontabInformerFactory.Start(stopCh)

	//启动控制器的控制循环
	if err := controller.Run(stopCh); err != nil {
		klog.Fatalf("crontab controller start fatal %s\n", err)
	}

}
