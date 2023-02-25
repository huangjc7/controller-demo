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
	//clientset, err := kubernetes.NewForConfig(config)
	//if err != nil {
	//	klog.Fatalf("Error init kubernetes clientset %s\n", err.Error())
	//}
	//实例化一个cronTab的clientset 需要用代码生成器根据CronTab生成的client进行声明
	crontabClientset, err := versioned.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error init kubernetes crontab clientset %s\n", err.Error())
	}
	crontabInformerFactory := externalversions.NewSharedInformerFactory(crontabClientset, time.Second*30)
	crontabInformer := crontabInformerFactory.Stable().V1beta1().CronTabs()
	crontabInformer.Informer()

}
