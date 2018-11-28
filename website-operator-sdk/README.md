## Create a new Operator SDK Project

The benefit of Operator SDK is that it generates lot of the scaffolding to make creating CRDs and custom controllers much, much easier.  To create a new Operator SDK project run the following command somewhere within your GOPATH:

```sh
#my GOPATH is ~/go
#I ran this command within ~/go/src/github.com/jungho/k8s-crds
operator-sdk new website-operator-sdk --skip-git-init
```

This will create a directory called `website-operator-sdk` with the following directory structure (note, it won't create the README.md file which I added after I ran the command.):

```sh
[~/go/src/github.com/jungho/website-operator-sdk, master+1]: tree -L 1
.
├── Gopkg.lock
├── Gopkg.toml
├── README.md
├── build
├── cmd
├── deploy
├── pkg
├── vendor
└── version

6 directories, 3 files
```

See [project layout](https://github.com/operator-framework/operator-sdk/blob/master/doc/project_layout.md) for description of each directory.

## Create the scaffolding for the Website CRD and Custom Controller

To add the Website CRD, from the website-controller directory and run the following command:

```sh
#You MUST run this within the website-controller-sdk project directory.  Otherwise it will fail.  
#This is because the command expects to find the `cmd/manager/main.go` file. 
operator-sdk add api --api-version=example.architech.ca/v1beta1 --kind=Website
```

This will generate some golang code as well as resource yaml files for your CRD.  The yaml files generated in the deploy directory:

```sh
[~/go/src/github.com/jungho/k8s-crds/website-operator-sdk/deploy, master]: tree -L 2
.
├── crds
│   ├── example_v1beta1_website_crd.yaml  # your Website CustomResourceDefinition
│   └── example_v1beta1_website_cr.yaml   # your Website resource
├── operator.yaml # The deployment resource to deploy your operator
├── role_binding.yaml # The RoleBinding that binds your ServiceAccount to the Role 
├── role.yaml # The Role that your ServiceAccount will be bound to.  Has the necessary permissions to access the apiserver.
└── service_account.yaml #The ServiceAccount that the operator will execute as
```

The golang code to represent your Website resource is generated in the pkg/apis/GROUP/VERSION directory.  We will modify the website_types.go file to define your Website resource in golang.

```sh
~/go/src/github.com/jungho/k8s-crds/website-operator-sdk/pkg/apis, master+1]: tree -L 3
.
├── addtoscheme_example_v1beta1.go
├── apis.go
└── example
    └── v1beta1
        ├── doc.go
        ├── register.go
        ├── website_types.go  #You will modify this file so you can consume your Website resource in golang
        └── zz_generated.deepcopy.go

2 directories, 6 files
```

Next, add the controller that will watch and reconcile Website resources.  

```sh
operator-sdk add controller --api-version=example.architech.ca/v1beta1 --kind=Website
```

The sdk will generate the code for your controller in the pkg/controller directory.

```sh
[~/go/src/github.com/jungho/k8s-crds/website-operator-sdk/pkg/controller, master+1]: tree -L 2
.
├── add_website.go
├── controller.go
└── website
    └── website_controller.go #You will modify this code to add your controller logic.

1 directory, 3 files
```
