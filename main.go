package main

import (
	"encoding/json"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"time"

	"github.com/golang/glog"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	restConfig, err := clientcmd.BuildConfigFromFlags("", "/home/cherwin/.kube/config")
	if err != nil {
		glog.Errorln(err)
	}
	client, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		glog.Errorln(err)
	}

	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(
		client,
		0,
		"db-ops",
		func(opts *v1.ListOptions) {
			opts.LabelSelector = "scripts=db-ops"
		},
	)

	gvr := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "configmaps",
	}

	informer := factory.ForResource(gvr)

	sharedInformer := informer.Informer()

	sharedInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
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
	})

	stop := make(chan struct{})
	defer close(stop)
	sharedInformer.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}
