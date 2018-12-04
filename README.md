# Kubernetes Custom Resource Definitions

Kubernetes is a highly extensible platform.  It supports many [extension points](https://kubernetes.io/docs/concepts/extend-kubernetes/extend-cluster/) to extend and customize your Kubernetes deployment without touching the core kubernetes source code.  Here I focus on Custom Resource Definitions, what they are, why they are useful, and how to implement them using [Kubebuilder](https://book.kubebuilder.io/) and [Operator SDK](https://github.com/operator-framework/operator-sdk).  I will explain CRDs using the excellent example from the book [Kubernetes in Action](https://www.manning.com/books/kubernetes-in-action) by [Marko Lukša](https://github.com/luksa).

**Note, this is an EXCELLENT book on Kubernetes and I highly recommend you read it!**

We will first dive into Marko's example as it is simple and clear.  However, his implementation is intentionally very simple and as he notes "barely working".  So we will reimplement his example using Kubebuilder and Operator SDK which enable you to build production grade CRDs and Controllers quickly.  

## What and Why?

But first, let's talk about what Custom Resource Definition are and why they are useful. WIP


## Marko Lukša Website Example ##

To demonstrate the process of creating a CRD, let's first deploy [Marko Lukša's](https://github.com/luksa) example.

The scenario is as follows:

We want to create a new resource called `WebSite`.  When we create this resource, a new WebSite based on the source code located at the specified git repo will be deployed and exposed in Kubernetes.  This will require a Deployment and Service resources to be created.

1. First step is to define our Website Custom Resource.  You do so by creating a `CustomResourceDefinition` resource. 

```sh

#We have a CRD defined in website-crd.yaml, create it like any other resource
kubectl create -f website-crd.yaml

#Verify it has been created
kubectl get customresourcedefinitions

NAME                              AGE
websites.extensions.example.com   16s

```

2. Now that the CRD has been created, create an instance of our Website resource.

```sh
kubectl create -f website.yaml

#Verify it has been created
kubectl get ws

NAME      AGE
kubia     4s

```
3. At this point, nothing happens because there is no controller that watches for the Website resource.  So we need to deploy a custom controller.  The source code for the controller is [here](https://github.com/luksa/k8s-website-controller)

```sh
#The website-controller will need a serviceAccount with sufficient privileges to access the kube-apiserver
kubectl create serviceaccount website-controller
kubectl create clusterrolebinding website-controller --clusterrole=cluster-admin --serviceaccount=default:website-controller

#create the website-controller as a deployment
kubectl create -f website-controller.yaml

#see what pods get deployed
kubectl get pods

#Notice the kubia-website pod has been deployed
NAME                                  READY     STATUS    RESTARTS   AGE
kubia-website-5645d5dc9-np68m         2/2       Running   0          3m
website-controller-84f9785c68-lzgmn   2/2       Running   1          3m

```

Below is a diagram from the [Kubernetes in Action](https://www.manning.com/books/kubernetes-in-action) book that describes the series of events that occur when the Website custom resource is deployed.

![Website Controller](./images/website-controller.png "Website Controller")

The website-controller pod contains 2 containers.

1. The controller itself that watches for 'Website' resources and deploys the website pod (which is just an nginx container that serves a static page)
2. A kubectl-proxy container as a 'side-car' container.  

The controller needs to communicate with the kube-apiserver in an authenticated manner.  The simplest way is to start up kubectl in proxy mode.  Since containers within the same pod share the same network namespace, the controller container can access the proxy via 'localhost'.  On start up, the controller gets a list of all 'Website' resource by sending a GET request like so:

```sh
#Notice the api path is group/version/resource as defined in the CRD
http://localhost:8001/apis/extensions.example.com/v1/websites?watch=true
```

See the following diagram from the [Kubernetes in Action](https://www.manning.com/books/kubernetes-in-action) book.

![Controller Pod](./images/controller-pod.png)

## Kubebuilder

TODO: finish implementation

## Operator SDK

**Note that Operator SDK is alpha software!  Proceed with caution!**

See the [user guide](https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md) to install Operator SDK locally on your machine.  Note the dependencies!  After you have installed the sdk, see [here](./website-operator-sdk/README.md) for instructions how to implement the Website CRD and custom controller and to deploy to AKS.

## References ##

- [Kuberenetes Deep Dive Series by RedHat.  Excellent overview of the kube-apiserver, code generation for CRDs, and how state is stored in etcd](https://blog.openshift.com/kubernetes-deep-dive-api-server-part-1/)
- [Writing Kubernetes Custom Controllers. Describes how to implement custom controllers using client-go.  Read this prior to diving into the sample-controller as it describes an established pattern for implementing controllers.](https://medium.com/@cloudark/kubernetes-custom-controllers-b6c7d0668fdf)
- [sample-controller.  Example implementing a customer controller using client-go library.](https://github.com/kubernetes/sample-controller)
