package main

import (
	"context"
	"io"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/bitnami-labs/sealed-secrets/pkg/kubeseal"
	"k8s.io/client-go/tools/clientcmd"
)

func initClient(r io.Reader) clientcmd.ClientConfig {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	return clientcmd.NewInteractiveDeferredLoadingClientConfig(loadingRules, nil, r)
}

func main() {
	clientConfig := initClient(os.Stdout)
	f, err := kubeseal.OpenCert(context.Background(), clientConfig, metav1.NamespaceSystem, "sealed-secrets-controller", "")
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(os.Stdout, f)
	if err != nil {
		panic(err)
	}

	f.Close()
}
