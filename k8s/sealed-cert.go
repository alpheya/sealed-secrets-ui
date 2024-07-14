package k8s

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func MakeHttpRequest() (string, error) {
	log.Info().Msg("Making HTTP request to get public key")
	uri := "http://sealed-secrets-controller.kube-system.svc.cluster.local:8080/v1/cert.pem"

	httpClient := &http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Info().Msg("Received public key")

	return string(body), nil
}

func GetPublicKey() (*rsa.PublicKey, error) {
	// clientConfig := initClient()
	// f, err := kubeseal.OpenCert(context.Background(), clientConfig, metav1.NamespaceSystem, "sealed-secrets-controller", "http://sealed-secrets-controller.kube-system.svc.cluster.local:8080/v1/cert.pem")
	// if err != nil {
	// 	panic(err)
	// }

	// defer f.Close()

	kubeSealCert, err := MakeHttpRequest()
	if err != nil {
		return nil, err
	}

	certPEM := []byte(kubeSealCert)

	// certPEM, err := io.ReadAll(
	// if err != nil {
	// 	panic(err)
	// }

	// Decode the PEM certificate
	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		panic("failed to decode PEM block containing public key")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(err)
	}

	// Extract the public key from the certificate
	rsaPubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		panic("public key type is not RSA")
	}

	return rsaPubKey, nil
}
