package main

import (
	"encoding/json"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"

	"github.com/golang/glog"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "/home/cherwin/.kube/config")
	if err != nil {
		glog.Errorln(err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Errorln(err)
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactoryWithOptions(
		clientSet,
		0,
		kubeinformers.WithNamespace("db-ops"),
		kubeinformers.WithTweakListOptions(func(opts *v1.ListOptions) {
			opts.LabelSelector = "scripts=db-ops"
		}),
	)

	cmInformer := kubeInformerFactory.Core().V1().ConfigMaps().Informer()

	cmInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			b, err := json.MarshalIndent(obj, "", "  ")
			if err != nil {
				return
			}
			fmt.Printf("service added: %s \n", b)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			b, err := json.MarshalIndent(newObj, "", "  ")
			if err != nil {
				return
			}
			fmt.Printf("service updated: %s \n", b)
		},
	},)

	stop := make(chan struct{})
	defer close(stop)
	kubeInformerFactory.Start(stop)
	for {
		time.Sleep(time.Second)
	}
}