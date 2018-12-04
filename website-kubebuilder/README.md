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

## Build and Deploy the operator