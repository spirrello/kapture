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
	"flag"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//externalKubeClient creates the external cluster config
func externalKubeClient(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil

}

//internalKubeClient creates the in-cluster config
func internalKubeClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func fetchPod(clientset *kubernetes.Clientset) {

	pods, err := clientset.CoreV1().Pods("kube-system").List(metav1.ListOptions{LabelSelector: "app=nfs-client-provisioner"})
	for _, pod := range pods.Items {
		fmt.Println(pod.Name, pod.Spec.NodeName)
	}

	if errors.IsNotFound(err) {
		fmt.Printf("Pod not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	}

}

func main() {

	//flag for external client
	var kubeconfig *string
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")

	flag.Parse()

	//if a kubeconfig is not provided we'll assume it's an in cluster deployment
	if *kubeconfig != "" {
		clientset, err := externalKubeClient(*kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		fetchPod(clientset)
	} else {
		clientset, err := internalKubeClient()
		if err != nil {
			panic(err.Error())
		}
		fetchPod(clientset)
	}

}
