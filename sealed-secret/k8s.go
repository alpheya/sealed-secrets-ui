package sealedsecret

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/rs/zerolog/log"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/util/homedir"
)

func getLocalClient() (*kubernetes.Clientset, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	return kubernetes.NewForConfig(config)
}

func getClusterClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get in-cluster config: %w", err)
	}

	return kubernetes.NewForConfig(config)
}

func decodeSecret(secretData map[string][]byte) map[string]string {
	data := make(map[string]string)
	for key, value := range secretData {
		data[key] = string(value)
	}

	return data
}

func getSecretData(ctx context.Context, namespace, secretName string) (map[string]string, error) {
	var clientset *kubernetes.Clientset
	var err error

	clientset, err = getClusterClient()
	if err != nil {
		clientset, err = getLocalClient()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	secret, err := clientset.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		log.Warn().Msg("Secret not found")
		return nil, nil
	}

	data := decodeSecret(secret.Data)

	return data, nil
}
