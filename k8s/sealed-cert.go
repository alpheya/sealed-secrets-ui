package k8s

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/bitnami-labs/sealed-secrets/pkg/kubeseal"
)

func GetPublicKey() (*rsa.PublicKey, error) {
	clientConfig := initClient(os.Stdout)
	f, err := kubeseal.OpenCert(context.Background(), clientConfig, metav1.NamespaceSystem, "sealed-secrets-controller", "")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	// _, err = io.Copy(os.Stdout, f)
	// if err != nil {
	// 	panic(err)
	// }
	certPEM, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	// Decode the PEM certificate
	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		panic("failed to decode PEM block containing public key")
	}

	// Parse the certificate
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
