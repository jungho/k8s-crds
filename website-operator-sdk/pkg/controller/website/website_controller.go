package website

import (
	"context"

	examplev1beta1 "github.com/jungho/k8s-crds/website-operator-sdk/pkg/apis/example/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_website")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Website Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileWebsite{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("website-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Website
	err = c.Watch(&source.Kind{Type: &examplev1beta1.Website{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Website
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
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
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Website object and makes changes based on the state read
// and what is in the Website.Spec.  It will create a Deployment and Service if they do not exist.  This is the key
// method that you need to implement after you generate the scaffolding.
//
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileWebsite) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Website")

	// Fetch the Website instance
	instance := &examplev1beta1.Website{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
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

	deployment := newDeploymentForWebsite(instance)

	// Set Website instance as the owner for the deployment and controller
	if err := controllerutil.SetControllerReference(instance, deployment, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.client.Create(context.TODO(), deployment)
		if err != nil {
			return reconcile.Result{}, err
		}

		//Deployment created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	//Deployment already exists - don't requeue
	reqLogger.Info("Skip reconcile: Deployment already exists", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	return reconcile.Result{}, nil
}

// returns a Deployment that will manage a set of Pods owned by the Website custom resource
func newDeploymentForWebsite(ws *examplev1beta1.Website) *appsv1.Deployment {
	labels := map[string]string{
		"webserver": ws.Name,
	}

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
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ws.Name + "-website",
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
					Name:   ws.Name + "-website",
					Labels: labels,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						v1.Container{
							Name:  "main",
							Image: "nginx:alpine",
							Ports: []v1.ContainerPort{
								v1.ContainerPort{
									ContainerPort: 80,
									Protocol:      v1.ProtocolTCP,
								},
							},
							VolumeMounts: []v1.VolumeMount{
								v1.VolumeMount{
									Name:      "html",
									MountPath: "/usr/share/nginx/html",
									ReadOnly:  true,
								},
							},
						},
						v1.Container{
							Name: "git-sync",
							Env: []v1.EnvVar{
								v1.EnvVar{
									Name:  "GIT_SYNC_REPO",
									Value: ws.Spec.GitRepo,
								},
								v1.EnvVar{
									Name:  "GIT_SYNC_DEST",
									Value: "/gitrepo",
								},
								v1.EnvVar{
									Name:  "GIT_SYNC_BRANCH",
									Value: "master",
								},
								v1.EnvVar{
									Name:  "GIT_SYNC_REV",
									Value: "FETCH_HEAD",
								},
								v1.EnvVar{
									Name:  "GIT_SYNC_WAIT",
									Value: "10",
								},
							},
							VolumeMounts: []v1.VolumeMount{
								v1.VolumeMount{
									Name:      "html",
									MountPath: "/gitrepo",
								},
							},
						},
					},
					Volumes: []v1.Volume{
						v1.Volume{
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
}

// returns a Service that will manage a set of Pods owned by the Website custom resource
func newServiceForWebsite(ws *examplev1beta1.Website) *corev1.Service {
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
		  sessionAffinity: None
		  type: LoadBalancer

	*/
	return nil
}
