/*
Copyright 2018 Jungho Kim.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package website

import (
	"context"

	examplev1beta1 "github.com/jungho/k8s-crds/website-kubebuilder/pkg/apis/example/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("controller_website")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Website Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
// USER ACTION REQUIRED: update cmd/manager/main.go to call this example.Add(mgr) to install this Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileWebsite{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("website-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Website
	err = c.Watch(&source.Kind{Type: &examplev1beta1.Website{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Watch for changes to secondary resource Deployment and Service and requeue the owner Website
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &examplev1beta1.Website{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &examplev1beta1.Website{},
	})

	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileWebsite{}

// ReconcileWebsite reconciles a Website object
type ReconcileWebsite struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Website object and makes changes based on the state read
// and what is in the Website.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.architech.ca,resources=websites,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileWebsite) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Website")
	// Fetch the Website ws
	ws := &examplev1beta1.Website{}
	err := r.Get(context.TODO(), request.NamespacedName, ws)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	labels := map[string]string{
		"webserver": ws.Name,
	}

	var deployment *appsv1.Deployment
	// Got the Website resource instance, now reconcile owned Deployment and Service resources
	deployment, err = r.newDeploymentForWebsite(ws, labels)

	if err != nil {
		return reconcile.Result{}, err
	}

	foundDeployment := &appsv1.Deployment{}
	// See if a Deployment already exists
	err = r.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, foundDeployment)

	// if the Deployment doesn't exist create it
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment",
			"Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.Create(context.TODO(), deployment)
		if err != nil {
			return reconcile.Result{}, err
		}

		//Deployment created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	//Deployment already exists, check the replica count int the status matches the desired replica count
	reqLogger.Info("Skip reconcile: Deployment already exists",
		"Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name)

	//make sure the replica count of the found deployment matches that of the spec
	if ws.Spec.Replicas != *foundDeployment.Spec.Replicas {
		foundDeployment.Spec.Replicas = &ws.Spec.Replicas

		if err = r.Update(context.TODO(), foundDeployment); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	var service *corev1.Service
	// Now reconcile the Service that is owned by the Website resource
	service, err = r.newServiceForWebsite(ws, labels)
	if err != nil {
		return reconcile.Result{}, err
	}

	//check if the service already exists
	foundService := &corev1.Service{}

	err = r.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, foundService)

	//If the service does not exist, create it
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Service", "Service.Namespace", service.Name, "Service.Name", service.Name)
		err = r.Create(context.TODO(), service)
		if err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	reqLogger.Info("Skip reconcile: Service already exists",
		"Service.Namespace", foundService.Namespace, "Service.Name", foundService.Name)

	return reconcile.Result{}, nil
}

// returns a Deployment that will manage a set of Pods owned by the Website custom resource
func (r *ReconcileWebsite) newDeploymentForWebsite(ws *examplev1beta1.Website, labels map[string]string) (*appsv1.Deployment, error) {

	/*
		We need to create Deployment resource that will be owned by our Website resource.
		We need to model the following deployment resource using the Go types defined in
		the k8s.io/api package which defines all the resource types.

		apiVersion: apps/v1
		kind: Deployment
		metadata:
		  labels:
		    webserver: kubia-website
		  name: kubia-website
		  namespace: default
		spec:
		  replicas: 1
		  selector:
		    matchLabels:
		      webserver: kubia-website
		  strategy:
		    rollingUpdate:
		      maxSurge: 1
		      maxUnavailable: 1
		    type: RollingUpdate
		  template:
		    metadata:
		      labels:
		        webserver: kubia-website
		      name: kubia-website
		    spec:
		      containers:
		      - image: nginx:alpine
		        imagePullPolicy: IfNotPresent
		        name: main
		        ports:
		        - containerPort: 80
		          protocol: TCP
		        resources: {}
		        volumeMounts:
		        - mountPath: /usr/share/nginx/html
		          name: html
		          readOnly: true
		      - image: openweb/git-sync
		        imagePullPolicy: Always
						name: git-sync
						env:
		        - name: GIT_SYNC_REPO
		          value: https://github.com/luksa/kubia-website-example.git
		        - name: GIT_SYNC_DEST
		          value: /gitrepo
		        - name: GIT_SYNC_BRANCH
		          value: master
		        - name: GIT_SYNC_REV
		          value: FETCH_HEAD
		        - name: GIT_SYNC_WAIT
		          value: "10"
		        resources: {}
		        volumeMounts:
		        - mountPath: /gitrepo
		          name: html
		      volumes:
		      - emptyDir: {}
		        name: html
	*/
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ws.Name,
			Namespace: ws.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &ws.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   ws.Name,
					Labels: labels,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "main",
							Image: "nginx:alpine",
							Ports: []v1.ContainerPort{
								{
									ContainerPort: 80,
									Protocol:      v1.ProtocolTCP,
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      "html",
									MountPath: "/usr/share/nginx/html",
									ReadOnly:  true,
								},
							},
						},
						{
							Name:  "git-sync",
							Image: "openweb/git-sync",
							Env: []v1.EnvVar{
								{
									Name:  "GIT_SYNC_REPO",
									Value: ws.Spec.GitRepo,
								},
								{
									Name:  "GIT_SYNC_DEST",
									Value: "/gitrepo",
								},
								{
									Name:  "GIT_SYNC_BRANCH",
									Value: "master",
								},
								{
									Name:  "GIT_SYNC_REV",
									Value: "FETCH_HEAD",
								},
								{
									Name:  "GIT_SYNC_WAIT",
									Value: "10",
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      "html",
									MountPath: "/gitrepo",
								},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "html",
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(ws, deployment, r.scheme); err != nil {
		return nil, err
	}

	return deployment, nil
}

// returns a Service that will manage a set of Pods owned by the Website custom resource
func (r *ReconcileWebsite) newServiceForWebsite(ws *examplev1beta1.Website, labels map[string]string) (*corev1.Service, error) {
	/*
		apiVersion: v1
		kind: Service
		metadata:
		  name: kubia-website-lb
		  namespace: default
		spec:
		  ports:
		  - name: http
		    port: 8080
		    targetPort: 80
		  selector:
		    webserver: kubia-website
		  type: LoadBalancer
	*/

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ws.Name + "-service",
			Namespace: ws.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       ws.Spec.Port,
					TargetPort: intstr.FromInt(int(ws.Spec.TargetPort)),
				},
			},
			Selector: map[string]string{"webserver": ws.Name},
			Type:     corev1.ServiceTypeLoadBalancer,
		},
	}

	// SetControllerReference sets owner as a Controller OwnerReference on owned.
	// This is used for garbage collection of the owned object and for
	// reconciling the owner object on changes to owned (with a Watch + EnqueueRequestForOwner).
	// Since only one OwnerReference can be a controller, it returns an error if
	// there is another OwnerReference with Controller flag set.
	if err := controllerutil.SetControllerReference(ws, service, r.scheme); err != nil {
		return nil, err
	}

	return service, nil
}
