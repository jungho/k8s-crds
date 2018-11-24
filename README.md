# Kubernetes Custom Resource Definitions #

Kubernetes is a highly extensible platform.  It supports many [extension points](https://kubernetes.io/docs/concepts/extend-kubernetes/extend-cluster/) to extend and customize your Kubernetes deployment without touching the core kubernetes source code.  Here I focus on Custom Resource Definitions, what they are, why they are useful, and how to implement them using [Kubebuilder](https://book.kubebuilder.io/) and [Operator SDK](https://github.com/operator-framework/operator-sdk).  I will explain CRDs using the excellent example from the book [Kubernetes in Action](https://www.manning.com/books/kubernetes-in-action) by [Marko Lukša](https://github.com/luksa) - This is a fantastic book, I really recommend you get it!  

We will first dive into Marko's example as it is simple and clear.  We will then reimplement his example using Kubebuilder and Operator SDK which enable you to build production grade CRDs and Controllers quickly.  

## What and Why? ##

## Kubebuilder ##

## Operator SDK ##

# References #

- [Kuberenetes Deep Dive Series by RedHat.  Excellent overview of the kube-apiserver, code generation for CRDs, and how state is stored in etcd](https://blog.openshift.com/kubernetes-deep-dive-api-server-part-1/)
- [Writing Kubernetes Custom Controllers. Describes how to implement custom controllers using client-go.  Read this prior to diving into the sample-controller as it describes an established pattern for implementing controllers.](https://medium.com/@cloudark/kubernetes-custom-controllers-b6c7d0668fdf)
- [sample-controller.  Example implementing a customer controller using client-go library.](https://github.com/kubernetes/sample-controller)
