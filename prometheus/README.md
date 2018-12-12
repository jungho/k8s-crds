## Deploying the Prometheus Operator

1. Install Helm.  See See https://docs.helm.sh/using_helm/#installing-helm for your OS.
2. Install the necessary role/role-bindings for Tiller.  Note this is only required for AKS `kubectel create -f tiller-rolebinding.yaml`
2. Init Helm from your AKS cluster.  `helm init --upgrade --serviceaccount tiller`.  If using Minikube do `helm init --upgrade`
3. Run the `deploy-operator.sh` script.  Note, the script will deploy the operator to the `monitoring` namespace.
4. Now look at what was deployed

```bash

#What did helm deploy?
[~/go/src/github.com/jungho/k8s-crds/prometheus, master]: helm ls
NAME                    REVISION        UPDATED                         STATUS          CHART  
kube-prometheus         1               Wed Dec 12 12:56:07 2018        DEPLOYED        kube-prometheus-0.0.105
prometheus-operator     1               Wed Dec 12 12:54:52 2018        DEPLOYED        prometheus-operator-0.0.29


#What CRDs were deployed?
[~/go/src/github.com/jungho/k8s-crds, master, 3s]: kubectl get crds                                                                    
NAME                                    CREATED AT
alertmanagers.monitoring.coreos.com     2018-12-04T15:37:50Z
prometheuses.monitoring.coreos.com      2018-12-04T15:37:50Z
prometheusrules.monitoring.coreos.com   2018-12-04T15:37:50Z
servicemonitors.monitoring.coreos.com   2018-12-04T15:37:50Z

#What Deployments, Statefulsets, Services, Pods were deployed?
.com/jungho/k8s-crds, master+1]: kubectl get pods -n monitoring          
NAME                                                   READY     STATUS    RESTARTS   AGE
alertmanager-kube-prometheus-0                         2/2       Running   0          52m
kube-prometheus-exporter-kube-state-76f498d465-4bz68   2/2       Running   0          52m
kube-prometheus-exporter-node-zftlk                    1/1       Running   0          52m
kube-prometheus-grafana-57d5b4d79f-9gdxh               2/2       Running   0          52m
prometheus-kube-prometheus-0                           3/3       Running   1          52m
prometheus-operator-d75587d6-mrtfg                     1/1       Running   0          53m

[~/go/src/github.com/jungho/k8s-crds, master+1]: kubectl get deployments -n monitoring
NAME                                  DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
kube-prometheus-exporter-kube-state   1         1         1            1           53m
kube-prometheus-grafana               1         1         1            1           53m
prometheus-operator                   1         1         1            1           54m

[~/go/src/github.com/jungho/k8s-crds, master+1]: kubectl get services -n monitoring   
NAME                                  TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)             AGE
alertmanager-operated                 ClusterIP   None             <none>        9093/TCP,6783/TCP   53m
kube-prometheus                       ClusterIP   10.100.185.114   <none>        9090/TCP            53m
kube-prometheus-alertmanager          ClusterIP   10.101.39.218    <none>        9093/TCP            53m
kube-prometheus-exporter-kube-state   ClusterIP   10.101.32.201    <none>        80/TCP              53m
kube-prometheus-exporter-node         ClusterIP   10.100.69.114    <none>        9100/TCP            53m
kube-prometheus-grafana               ClusterIP   10.106.106.127   <none>        80/TCP              53m
prometheus-operated                   ClusterIP   None             <none>        9090/TCP            53m

[~/go/src/github.com/jungho/k8s-crds, master+1]: kubectl get statefulsets -n monitoring
NAME                           DESIRED   CURRENT   AGE
alertmanager-kube-prometheus   1         1         54m
prometheus-kube-prometheus     1         1         54m
```

To tear everything down, run `helm delete prometheus-operator && helm delete kube-prometheus && kubectl delete namespace monitoring`