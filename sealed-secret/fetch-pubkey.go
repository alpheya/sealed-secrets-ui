package sealedsecret

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func (s SealedSecretService) makeHttpRequest(ctx context.Context) (string, error) {
	log.Info().Msg("making HTTP request to get public key")
	uri := fmt.Sprintf("http://%s.%s.svc.cluster.local:8080/v1/cert.pem", s.sealedSecretControllerName, s.sealedSecretControllerNamespace)

	httpClient := http.Client{Timeout: time.Duration(2) * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/x-pem-file")

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

func (s SealedSecretService) getPublicKey(ctx context.Context) (*rsa.PublicKey, error) {
	kubeSealCert, err := s.makeHttpRequest(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}
	certPEM := []byte(kubeSealCert)

	// Decode the PEM certificate
	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to decode PEM block containing certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Extract the public key from the certificate
	rsaPubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to extract public key from certificate")
	}

	return rsaPubKey, nil
}
