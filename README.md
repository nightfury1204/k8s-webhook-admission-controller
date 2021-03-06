# k8s-webhook-admission-controller
validating webhook admission controller for k8s CRD

### Example

```
$ cat ./yaml/crd/example_crd_obj.yaml 
apiVersion: "nahid.try.com/v1alpha1"
kind: PodWatch
metadata:
  name: my-podwatch
spec:
  replicas: 4
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - name: ubuntu
        image: ubuntu:latest
        # Just spin & wait forever
        command: [ "/bin/bash", "-c", "--" ]
        args: [ "while true; do sleep 30; done;" ]
```

```     
# create podwatch resource
$ kubectl create -f ./yaml/crd/example_crd_obj.yaml
 Error from server (Number of replicas must be in between 1 and 3): error when creating "./yaml/crd/example_crd_obj.yaml": admission webhook "podwatch-image.nahid.try.com" denied the request: Number of replicas must be in between 1 and 3
 
```
### Others

Send ```--admission-control-config-file``` flag with admissionConfiguration file to kube apiserver.

```yaml
#admissionConfiguration demo
kind: AdmissionConfiguration
apiVersion: apiserver.k8s.io/v1alpha1
plugins:
- name: ValidatingAdmissionWebhook
  configuration:
    kind: WebhookAdmission
    apiVersion: apiserver.config.k8s.io/v1alpha1
    kubeConfigFile: <PATH TO KUBECONFIG>
```
