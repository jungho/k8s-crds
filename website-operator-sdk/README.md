## Create a new Operator SDK Project

The benefit of Operator SDK is that it generates lot of the scaffolding to make creating CRDs and custom controllers much, much easier.  To create a new Operator SDK project run the following command somewhere within your GOPATH:

```sh
#my GOPATH is ~/go
#I ran this command within ~/go/src/github.com/jungho/k8s-crds
operator-sdk new website-operator-sdk --skip-git-init
```

This will create a directory called `website-operator` with the following directory structure (note, it won't create the README.md file which I added after I ran the command.):

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

## Create the scaffolding for the Website CRD and Custom Controller

To add the Website CRD and Controller cd into the website-controller directory and run the following command:

```sh
#You MUST run this within the website-controller project directory.  Otherwise it will fail.  
#This is because the command expects to find the `cmd/manager/main.go` file. 
operator-sdk add api --api-version=example.architech.ca/v1 --kind=Website
```
