package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	//"k8s.io/apimachinery/pkg/api/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

const (
	namespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
)

func readSecretFromK8s(secretName string, namespace string, isOutsideOfPod bool, usernameField string, passwordField string) {
	var config *rest.Config
	var err error

	if isOutsideOfPod {
		if namespace == "" {
			panic("Namespace must be specified")
		}

		var configPath string

		if runtime.GOOS == "windows" {
			configPath = os.Getenv("USERPROFILE")
		} else {
			configPath = "~"
		}

		config, err = clientcmd.BuildConfigFromFlags("", configPath+"/.kube/config")

		if err != nil {
			panic(err.Error())
		}
	} else {
		ns, err := os.ReadFile(namespaceFile)
		if err != nil {
			panic(err.Error())
		}

		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}

		namespace = string(ns)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	value, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("username =", string(value.Data[usernameField]))
	fmt.Println("password =", string(value.Data[passwordField]))
}

func readSecretFromAWS(secretName string, region string, isOutsideOfPod bool, usernameField string, passwordField string) {

	var cfg aws.Config
	var err error

	if isOutsideOfPod {
		cfg, err = config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(region),
			config.WithSharedConfigProfile("default"),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	}

	if err != nil {
		log.Fatal(err)
	}

	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatal(err.Error())
	}

	var secretString string = *result.SecretString

	var sec map[string]interface{}

	if err := json.Unmarshal([]byte(secretString), &sec); err != nil {
		panic(err)
	}

	fmt.Println("username =", sec[usernameField])
	fmt.Println("password =", sec[passwordField])
}

func main() {
	secretName := flag.String("secret", "", "name of secret")
	namespace := flag.String("namespace", "default", "name of namespace where secret is located")
	region := flag.String("region", "", "AWS region where secret is located")
	isOutsideOfPod := flag.Bool("outside", false, "Is app running outside of pod?")
	usernameField := flag.String("username", "username", "name of username field in secret")
	passwordField := flag.String("password", "password", "name of password field in secret")
	isAWSSecret := flag.Bool("aws", false, "Is secret stored in AWS?")
	isK8sSecret := flag.Bool("k8s", false, "Is secret stored in k8s?")
	readInterval := flag.Int("interval", 10, "Read interval (in seconds). Default: 10 seconds")

	flag.Parse()

	duration := time.Duration(*readInterval)

	if !*isAWSSecret && !*isK8sSecret {
		panic("specify .")
	}

	if *secretName == "" {
		panic("secret name cannot be empty.")
	}

	if *isK8sSecret {
		for {
			readSecretFromK8s(*secretName, *namespace, *isOutsideOfPod, *usernameField, *passwordField)
			time.Sleep(duration * time.Second)
		}
	}

	if *isAWSSecret {
		if *region == "" {
			panic("region cannot be empty.")
		}
		for {
			readSecretFromAWS(*secretName, *region, *isOutsideOfPod, *usernameField, *passwordField)
			time.Sleep(duration * time.Second)
		}
	}
}
