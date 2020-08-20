package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var clientset *kubernetes.Clientset
var podName string
var podPort string
var namespace string

func FetchIPsFromCluster(podName, namespace string) []string {
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	ips := []string{}

	if len(pods.Items) >= 1 {
		for _, pod := range pods.Items {
			if pod.Name[:len(podName)] == podName {
				ips = append(ips, pod.Status.PodIP)
			}
		}
	}

	return ips
}

func exporterHandler(w http.ResponseWriter, r *http.Request) {
	resp := ""

	ips := FetchIPsFromCluster(podName, namespace)
	fmt.Println("Fetching metrics for ips:", ips)

	for _, ip := range ips {
		u := "http://" + ip + ":" + podPort + "/metrics"
		r, err := http.Get(u)
		if err != nil {
			fmt.Println(err)
		} else {
			body := r.Body
			bytes, err := ioutil.ReadAll(body)
			if err != nil {
				fmt.Println(err)
			} else {
				resp += string(bytes)
			}
		}
	}

	w.Write([]byte(resp))
}

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	podName = os.Getenv("NGINX_EXPORTER_POD_NAME")
	namespace = os.Getenv("NGINX_EXPORTER_NAMESPACE")
	podPort = os.Getenv("NGINX_EXPORTER_POD_PORT")

	if podName == "" {
		podName = "nginx-controller"
		fmt.Println("No pod name defined, using ", podName)
	}
	if namespace == "" {
		namespace = "default"
		fmt.Println("No pod name defined, using ", namespace)
	}
	if podPort == "" {
		podPort = "10254"
		fmt.Println("No pod port defined, using ", podPort)
	}

	http.HandleFunc("/metrics", exporterHandler)

	port := os.Getenv("NGINX_EXPORTER_PORT")
	if port == "" {
		port = "8888"
	}
	fmt.Printf("Fetching from pod '%s' port '%s' on namespace '%s'\n", podName, podPort, namespace)

	fmt.Printf("Listening at port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}
