## Creating a Kubebuilder Project

The benefit of Kubebuilder is that it generates lot of the scaffolding to make creating CRDs and custom controllers much, much easier.  To create a new Kubebuilder project run the following command somewhere within your GOPATH:

```sh
#my GOPATH is ~/go
#I ran this command within ~/go/src/github.com/jungho/k8s-crds/website-kubebuilder
kubebuilder init --domain architech.ca --owner "Jungho Kim"

#You will be asked the following, answer 'y'.  This will download all the required Go packages and place them in the vendor directory and generate the scaffolding
Run `dep ensure` to fetch dependencies (Recommended) [y/n]?

#Now create a new API resource and Controller, answer 'y' to both Resource and Controller
kubebuilder create api --group example --version v1beta1 --kind Website
```

Kubebuilder will generate the following directories as well as deployment yaml files, CRDs, golang types, functions for your API and Controller as well as tests.

```sh
[~/go/src/github.com/jungho/k8s-crds/kubebuilder, master+2]: tree -L 1
.
├── Dockerfile  #To containerize your Controller. You can modify to use a different base image.
├── Gopkg.lock
├── Gopkg.toml
├── Makefile    #Makefile to help with development workflow
├── PROJECT
├── bin         #the build output for your Controller
├── cmd         #The Go code to start the Manager that starts your Controller
├── config      #containers the deployment yaml files, CRD and sample instance, RBAC role, rolebindings
├── cover.out
├── hack
├── pkg         #Contains the Go code for your API, Controller and Reconciler
└── vendor      #required vendor libraries such as controller-runtime

6 directories, 6 files
```

Let's take a deeper look into the `pkg` directory.

```sh
[~/go/src/github.com/jungho/k8s-crds/kubebuilder/pkg, master+2]: tree -L 4
.
├── apis
│   ├── addtoscheme_example_v1beta1.go
│   ├── apis.go
│   └── example
│       ├── group.go
│       └── v1beta1
│           ├── doc.go
│           ├── register.go
│           ├── v1beta1_suite_test.go
│           ├── website_types.go  #You will modify this file so you can consume your Website resource in golang
│           ├── website_types_test.go
│           └── zz_generated.deepcopy.go
├── controller
│   ├── add_website.go
│   ├── controller.go
│   └── website
│       ├── website_controller.go #You will modify this code to add your reconciliation logic.
│       ├── website_controller_suite_test.go
│       └── website_controller_test.go
└── webhook
    └── webhook.go

6 directories, 15 files
```

## Modifying the generated code to add our reconciliation logic

First we modify [pkg/apis/example/v1beta1/website_types.go](https://github.com/jungho/k8s-crds/blob/master/website-kubebuilder/pkg/apis/example/v1beta1/website_types.go#L27:6). This file contains the generated golang struct types for your Website resource. The SDK generates a skeleton, you need to take it the rest of the way.  For our purposes, we want the GitRepo, Replicas, Port and TargetPort fields added the spec and Replicas to the status.

```go
// WebsiteSpec defines the desired state of Website
type WebsiteSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	GitRepo    string `json:"gitRepo"`
	Replicas   int32  `json:"replicas"`
	Port       int32  `json:"port"`
	TargetPort int32  `json:"targetPort"`
}
```

Note, as per the comments, whenever you makes changes to this file, you need to run `make` to update other sdk generated files such as 
[pkg/apis/example/v1beta1/website_types.go](./pkg/apis/example/v1beta1/zz_generated.deepcopy.go).

Next, you need to implement the reconciliation logic by updating [pkg/controller/website/website_controller.go](./pkg/controller/website/website_controller.go).

The key methods are:

```go
// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error 
```
This is the method that makes the Manager aware of your controller and is also where you specify which resources
your controller "watches" for changes.  We want to watch for Website, Deployment and Service resources.  Note, we don't
care about all Deployment and Service resources, only those "owned" by Website resources.  
See [pkg/controller/website/website_controller.go](./pkg/controller/website/website_controller.go).

```go
// Reconcile reads the state of the cluster for a Website object and makes changes based on the state read
// and what is in the Website.Spec.  It will create a Deployment and Service if they do not exist.  This is the key
// method that you need to implement after you generate the scaffolding.
//
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileWebsite) Reconcile(request reconcile.Request) (reconcile.Result, error) 
```

The Reconcile function is part of the [Reconciler](https://github.com/jungho/k8s-crds/blob/master/website-kubebuilder/vendor/sigs.k8s.io/controller-runtime/pkg/reconcile/reconcile.go#L79:6) interface. The generated [ReconcileWebsite struct](https://github.com/jungho/k8s-crds/blob/master/website-kubebuilder/pkg/controller/website/website_controller.go#L87:6) satisfies this interface.  It is responsible for implementing the reconciliation logic and will be invoked for each ADD, UPDATE, DELETE event for our Website resource.  See the [Controller Runtime Client API](https://github.com/operator-framework/operator-sdk/blob/master/doc/user/client.md) for the key interfaces.

## Build and Deploy the operator