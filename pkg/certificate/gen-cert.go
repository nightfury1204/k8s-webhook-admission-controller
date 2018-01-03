package certificate

import (
	"crypto/x509"
	"fmt"
	"net"
	"path/filepath"

	"k8s.io/client-go/util/cert"
)

// initcertCmd represents the initcert command
func GenerateCertificate(writeDir string, ip string,dns string, genClientCert bool) error {
	dir := filepath.Join(writeDir, "pki")
	createDir(dir)

	cfg := cert.Config{
		CommonName: "ca",
	}

	caKey, err := cert.NewPrivateKey()
	if err != nil {
		return fmt.Errorf("Failed to generate private key. Reason: %v.", err)
	}
	caCert, err := cert.NewSelfSignedCACert(cfg, caKey)
	if err != nil {
		return fmt.Errorf("Failed to generate self-signed certificate. Reason: %v.", err)
	}
	err = WriteCertKey(filepath.Join(dir, "ca"), caCert, caKey)
	if err != nil {
		return fmt.Errorf("Failed to init ca. Reason: %v.", err)
	}
	fmt.Println("Wrote ca certificates in ", dir)

	altNames := cert.AltNames{}
	if len(ip)>0 {
		altNames.IPs = []net.IP{net.ParseIP(ip)}
	}
	if len(dns)>0 {
		altNames.DNSNames = []string{dns}
	}

	//create server.cert,server.key
	cfgForServer := cert.Config{
		CommonName: "server",
		AltNames: altNames,
		Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	serverKey, err := cert.NewPrivateKey()
	if err != nil {
		return fmt.Errorf("Failed to generate private key. Reason: %v.", err)
	}
	serverCert, err := cert.NewSignedCert(cfgForServer, serverKey, caCert, caKey)
	if err != nil {
		return fmt.Errorf("Failed to generate server certificate. Reason: %v.", err)
	}
	err = WriteCertKey(filepath.Join(dir, "server"), serverCert, serverKey)
	if err != nil {
		return fmt.Errorf("Failed to init server certificate pair. Reason: %v.", err)
	}
	fmt.Println("Wrote server certificates in ", dir)

	if !genClientCert {
		return nil
	}
	//create client.cert,client.key
	cfgForClient := cert.Config{
		CommonName: "client",
		Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	clientKey, err := cert.NewPrivateKey()
	if err != nil {
		return fmt.Errorf("Failed to generate private key. Reason: %v.", err)
	}
	clientCert, err := cert.NewSignedCert(cfgForClient, clientKey, caCert, caKey)
	if err != nil {
		return fmt.Errorf("Failed to generate client certificate. Reason: %v.", err)
	}
	err = WriteCertKey(filepath.Join(dir, "client"), clientCert, clientKey)
	if err != nil {
		return fmt.Errorf("Failed to init client certificate pair. Reason: %v.", err)
	}
	fmt.Println("Wrote client certificates in ", dir)

	return nil
}
