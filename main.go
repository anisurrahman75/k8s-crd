package main

import (
	"context"
	"flag"
	"fmt"
	myv1 "github.com/anisurrahman75/my-crd/pkg/apis/mycrd.dev/v1"
	sbclientset "github.com/anisurrahman75/my-crd/pkg/client/clientset/versioned"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	crdclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	_ "k8s.io/code-generator"
	_ "k8s.io/utils/pointer"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

func main() {
	log.Println("Configuring KubeConfig......")
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	customCRD := v1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "appscodes.mycrd.dev",
		},
		Spec: v1.CustomResourceDefinitionSpec{
			Group: "mycrd.dev",
			Versions: []v1.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Served:  true,
					Storage: true,
					Schema: &v1.CustomResourceValidation{
						OpenAPIV3Schema: &v1.JSONSchemaProps{
							Type: "object",

							Properties: map[string]v1.JSONSchemaProps{
								"spec": {
									Type: "object",
									Properties: map[string]v1.JSONSchemaProps{
										"name": {
											Type: "string",
										},
										"replicas": {
											Type: "integer",
										},
										"container": {
											Type: "object",
											Properties: map[string]v1.JSONSchemaProps{
												"image": {
													Type: "string",
												},
												"port": {
													Type: "integer",
												},
											},
										},
									},
								},
							},
						},
					},
					Subresources: &v1.CustomResourceSubresources{
						Scale: &v1.CustomResourceSubresourceScale{
							SpecReplicasPath:   ".spec.replicas",
							StatusReplicasPath: ".status.replicas",
							LabelSelectorPath:  nil,
						},
					},
				},
			},
			Scope: "Namespaced",
			Names: v1.CustomResourceDefinitionNames{
				Kind:     "AppsCode",
				Plural:   "appscodes",
				Singular: "appscode",
				ShortNames: []string{
					"ac",
				},
				Categories: []string{
					"all",
				},
			},
		},
	} // deleting existing CR
	DelCR(config, customCRD)

	// creating new one
	CreateCR(config, customCRD)

	time.Sleep(3 * time.Second)
	log.Println("CRD is Created!")
	// ..........Create ishtiaqvai objects
	log.Println("Press ctrl+c to create a ishtiaqvai resources")
	HandleUtils() // wait until ctrl+c
	createRes(config)

	// ..........Delete ishtiaqvai objects
	log.Println("Press ctrl+c to Delete  ishtiaqvai objects")
	HandleUtils() // wait until ctrl+c
	DeleRes(config)

	// ..........Delete AppsCode Custom Resources
	log.Println("Press ctrl+c to Delete AppsCode CR")
	HandleUtils() // wait until ctrl+c
	DelCR(config, customCRD)
}
func HandleUtils() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
func CreateCR(config *rest.Config, customCRD v1.CustomResourceDefinition) {
	crdClient, err := crdclientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	_, err = crdClient.ApiextensionsV1().CustomResourceDefinitions().Create(context.TODO(), &customCRD, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	log.Println("Creating Custom Resources Name: AppsCode")
	time.Sleep(2 * time.Second)
}
func DelCR(config *rest.Config, customCRD v1.CustomResourceDefinition) {
	crdClient, err := crdclientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	err = crdClient.ApiextensionsV1().CustomResourceDefinitions().Delete(context.TODO(), customCRD.Name, metav1.DeleteOptions{})
	if err != nil {
		//panic(err)
		fmt.Printf("Del error: %s", err.Error())
	}
	time.Sleep(2 * time.Second)
	time.Sleep(2 * time.Second)
}
func DeleRes(config *rest.Config) {
	resName := "ishtiaqvai"
	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	client, err := sbclientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	if err := client.MycrdV1().AppsCodes("default").Delete(context.TODO(), resName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}

	log.Println("ishtiaqvai resources  Deleted!!")
	time.Sleep(2 * time.Second)
}

// create demo resources
func createRes(config *rest.Config) {
	sbObj := &myv1.AppsCode{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ishtiaqvai",
		},
		Spec: myv1.AppsCodeSpec{
			Name:     "customappscode",
			Replicas: intptr(2),
			Container: myv1.ContainerSpec{
				Image: "havijavi/go-api-server",
				Port:  3000,
			},
		},
	}
	// subclientset= from client/clientset
	client, err := sbclientset.NewForConfig(config)
	_, err = client.MycrdV1().AppsCodes("default").Create(context.TODO(), sbObj, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	log.Println("ishtiaqvai resources  Created!!")
	time.Sleep(2 * time.Second)
}
func intptr(i int32) *int32 {
	return &i
}
