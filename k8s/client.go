package k8s

import (
	"k8s.io/client-go/tools/clientcmd"
)

func initClient() clientcmd.ClientConfig {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, nil)
}
