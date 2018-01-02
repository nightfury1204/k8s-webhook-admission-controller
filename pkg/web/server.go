package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s-webhook-admission-controller/pkg/certificate"
	"log"
	"net/http"
	"time"

	"github.com/bmizerany/pat"
	"github.com/golang/glog"
	adv1beta1 "k8s.io/api/admission/v1beta1"
)

func performTask(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		glog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	reviewResponse, err := validateResource(body)
	if err != nil {
		glog.Error(err)
	}
	review := adv1beta1.AdmissionReview{
		Response: reviewResponse,
	}

	resp, err := json.Marshal(review)
	if err != nil {
		glog.Error(err)
	}
	if _, err := w.Write(resp); err != nil {
		glog.Error(err)
	}
}

func RunServer(port, kubeconfig string) error {
	m := pat.New()
	m.Post("/", http.HandlerFunc(performTask))

	clienset, err := getClient(kubeconfig)
	if err != nil {
		return err
	}

	serverCert, err := certificate.ReadCert("pki/server.cert")
	if err != nil {
		return fmt.Errorf("failed to read server.cert. Reason: %v.", err)
	}

	serverKey, err := certificate.ReadKey("pki/server.key")
	if err != nil {
		return fmt.Errorf("failed to read server.key. Reason: %v.", err)
	}

	tlsConfig, err := getTlsConfig(clienset, serverCert, serverKey)
	if err != nil {
		return fmt.Errorf("failed to configure tls. Reason: %v.", err)
	}
	//tlsConfig.BuildNameToCertificate()
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      m,
		TLSConfig:    tlsConfig,
	}
	log.Println("Server starting.....")
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		return fmt.Errorf("Reason: %v.", err)
	}
	return nil
}
