package web

import (
	"net/http"
	"fmt"
	"time"
	"github.com/bmizerany/pat"
	"k8s-webhook-admission-controller/pkg/certificate"
)

func performTask(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w,"Success!!!")
}

func RunServer(port,kubeconfig string) error {
	m := pat.New()
	m.Post("/",http.HandlerFunc(performTask))

	clienset, err:= getClient(kubeconfig)
	if err!=nil {
		return err
	}

	serverCert, err := certificate.ReadCert("/home/ac/go/src/k8s-webhook-admission-controller/pki/server.cert")
	if err!=nil {
		return fmt.Errorf("failed to read server.cert. Reason: %v.",err)
	}

	serverKey, err := certificate.ReadCert("/home/ac/go/src/k8s-webhook-admission-controller/pki/server.key")
	if err!=nil {
		return fmt.Errorf("failed to read server.key. Reason: %v.",err)
	}

	tlsConfig, err := getTlsConfig(clienset,serverCert,serverKey)
	if err!=nil {
		return fmt.Errorf("failed to configure tls. Reason: %v.",err)
	}
	//tlsConfig.BuildNameToCertificate()
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      m,
		TLSConfig:    tlsConfig,
	}
	fmt.Println("Server starting.....")
	err = server.ListenAndServeTLS("","")
	if err != nil {
		return fmt.Errorf("Reason: %v.",err)
	}
	return nil
}