package k8s

import (
	"io"

	"k8s.io/client-go/tools/clientcmd"
)

func initClient(r io.Reader) clientcmd.ClientConfig {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	return clientcmd.NewInteractiveDeferredLoadingClientConfig(loadingRules, nil, r)
}
