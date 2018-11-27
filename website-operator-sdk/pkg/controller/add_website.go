package controller

import (
	"github.com/jungho/k8s-crds/website-operator-sdk/pkg/controller/website"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, website.Add)
}
