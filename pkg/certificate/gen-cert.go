package certificate

import (
	"crypto/x509"
	"fmt"
	"net"
	"path/filepath"

	"k8s.io/client-go/util/cert"
	"crypto/rsa"
)

// initcertCmd represents the initcert command
func GenerateCertificate(writeDir string, ip string,dns string,genCaCert bool,genServerCert bool, genClientCert bool) error {
	dir := filepath.Join(writeDir, "pki")
	createDir(dir)

	if genCaCert {
		err :=generateCaCertificates(dir)
		if err!=nil {
			return err
		}
	}

	if genServerCert {
		err := generateServerCertificates(dir,ip,dns)
		if err!=nil {
			return err
		}
	}

	if genClientCert {
		err := generateClientCertificates(dir)
		if err!=nil {
			return err
		}
	}

	return nil
}

func generateCaCertificates(dir string) error {
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
	return nil
}

func generateServerCertificates(dir,ip,dns string) error {
	caCert,caKey,err := getCaPair()
	if err !=nil {
		return err
	}

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
	return nil
}

func generateClientCertificates(dir string) error {
	caCert,caKey,err := getCaPair()
	if err !=nil {
		return err
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

func getCaPair() (*x509.Certificate,*rsa.PrivateKey, error) {
	caCertBytes,err := ReadCert("pki/ca.cert")
	if err !=nil {
		return nil,nil,err
	}
	caCert,err := cert.ParseCertsPEM(caCertBytes)
	if err !=nil {
		return nil,nil,err
	}
	caKeyBytes,err := ReadCert("pki/ca.key")
	if err !=nil {
		return nil,nil,err
	}
	caKey,err := cert.ParsePrivateKeyPEM(caKeyBytes)
	if err !=nil {
		return nil,nil,err
	}
	return caCert[0],caKey.(*rsa.PrivateKey),nil
}