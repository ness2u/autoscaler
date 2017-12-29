package main

import (
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/tools/clientcmd"
)

type DeploymentState struct {
	Namespace  string
	Deployment string
	Replicas   int64
}

func SetDeploymentScale(scaleConfig ScaleConfig, n int64) {

	clientset := getClientset()

	deploymentsClient := clientset.AppsV1beta1().Deployments(scaleConfig.Namespace)
	result, err := deploymentsClient.Get(scaleConfig.Deployment, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("updating deployment: %s (%d --> %d replicas)\n", result.Name, *result.Spec.Replicas, n)

	var nnn = int32(n)
	result.Spec.Replicas = &nnn

	res, err := deploymentsClient.Update(result)
	fmt.Printf("deployment: %s (%d replicas)\n", res.Name, *res.Spec.Replicas)
}

func GetDeployment(scaleConfig ScaleConfig) DeploymentState {

	clientset := getClientset()

	deploymentsClient := clientset.AppsV1beta1().Deployments(scaleConfig.Namespace)
	result, err := deploymentsClient.Get(scaleConfig.Deployment, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("deployment: %s (%d replicas)\n", result.Name, *result.Spec.Replicas)

	return DeploymentState{scaleConfig.Namespace, scaleConfig.Deployment, int64(*result.Spec.Replicas)}
}

func getClientset() *kubernetes.Clientset {
	var kubeconfig = filepath.Join(homeDir(), ".kube", "config")
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
