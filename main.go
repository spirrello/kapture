/*
Copyright 2016 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		//_, err = clientset.CoreV1().Pods("default").Get("example-xxxxx", metav1.GetOptions{})
		pods, err = clientset.CoreV1().Pods("kube-system").List(metav1.ListOptions{FieldSelector: "metadata.labels=nfs-client-provisioner"})
		for _, pod := range pods.Items {
			fmt.Println(pod.Name, pod.Status)
		}
		//NEED TO FIX THIS VAR.
		// var labelOptions = kubernetes.CoreV1().Pods.GetOptions
		// appLabel := []byte(`{"app": "nfs-client-provisioner"}`)
		// err := json.Unmarshal(appLabel, &labelOptions)
		// _, err = clientset.CoreV1().Pods("kube-system").Get("", metav1.GetOptions{labelOptions, "0"})

		if errors.IsNotFound(err) {
			fmt.Printf("Pod not found\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found pod\n")
		}

		time.Sleep(10 * time.Second)
	}
}
