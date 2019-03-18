package main

import (
	"flag"
	"log"

	"github.com/kubernetes-sigs/sig-storage-lib-external-provisioner/controller"
	"github.com/sdorra/host-path-provisioner/pkg/storage"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	directory := flag.String("directory", "/volumes", "base path for volumes")
	flag.Parse()
	flag.Set("logtostderr", "true")

	log.Printf("start with %s as volume directory", *directory)

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	serverVersion, err := clientset.Discovery().ServerVersion()
	if err != nil {
		log.Fatalf("Error getting server version: %v", err)
	}

	provisioner := storage.NewHostPathProvisioner(*directory)

	pc := controller.NewProvisionController(clientset, storage.NAME, provisioner, serverVersion.GitVersion)
	pc.Run(wait.NeverStop)
}
