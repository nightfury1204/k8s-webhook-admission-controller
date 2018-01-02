package web

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"

	admissionv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s-webhook-admission-controller/pkg/certificate"
)

// get a clientset with in-cluster config.
func getClient(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create config. Reason: %v.",err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset. Reason: %v.",err)
	}
	return clientset,nil
}

func getAPIServerCert(clientset *kubernetes.Clientset) []byte {
	c, err := clientset.CoreV1().ConfigMaps("kube-system").Get("extension-apiserver-authentication", metav1.GetOptions{})
	if err != nil {
		glog.Fatal(err)
	}

	pem, ok := c.Data["requestheader-client-ca-file"]
	if !ok {
		glog.Fatalf(fmt.Sprintf("cannot find the ca.crt in the configmap, configMap.Data is %#v", c.Data))
	}
	glog.Info("client-ca-file=", pem)
	return []byte(pem)
}

func getTlsConfig(clientset *kubernetes.Clientset, serverCert []byte, serverKey []byte) (*tls.Config, error) {
	cert := getAPIServerCert(clientset)
	//for test
	cert ,err := certificate.ReadCert("/home/ac/go/src/k8s-webhook-admission-controller/pki/ca.cert")
	if err!=nil {
		return nil,err
	}
	apiserverCA := x509.NewCertPool()
	apiserverCA.AppendCertsFromPEM(cert)

	sCert, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		return nil,err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{sCert},
		ClientCAs:    apiserverCA,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	},nil
}

// register this example webhook admission controller with the kube-apiserver
// by creating .ValidatingWebhookConfigurations
func selfRegistration(clientset *kubernetes.Clientset, caCert []byte) {
	time.Sleep(10 * time.Second)
	client := clientset.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations()
	_, err := client.Get("example-config", metav1.GetOptions{})
	if err == nil {
		if err2 := client.Delete("example-config", nil); err2 != nil {
			glog.Fatal(err2)
		}
	}
	webhookConfig := &admissionv1beta1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "validating-webhook-config",
		},
		Webhooks: []admissionv1beta1.Webhook{
			{
				Name: "podwatch-image.nahid.try.com",
				Rules: []admissionv1beta1.RuleWithOperations{{
					Operations: []admissionv1beta1.OperationType{admissionv1beta1.Create},
					Rule: admissionv1beta1.Rule{
						APIGroups:   []string{"nahid.try.com"},
						APIVersions: []string{"v1alpha1"},
						Resources:   []string{"podwatchs"},
					},
				}},
				ClientConfig: admissionv1beta1.WebhookClientConfig{
					Service: &admissionv1beta1.ServiceReference{
						Namespace: "default",
						Name:      "validating-webhook",
					},
					CABundle: caCert,
				},
			},
		},
	}
	if _, err := client.Create(webhookConfig); err != nil {
		glog.Fatal(err)
	}
}
