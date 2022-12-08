# auto-ingress-operator

## Description
为k8s集群里的service自动创建对应的ingress，可以指定service前缀，以及通过黑白名单来指定要生成的namespace

域名规则: `<serviceName>---<namespace>.<autoIngressName>`

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1.安装Operator

```sh
kubectl apply -f deploy/auto-ingress-operator.yaml
```


2.创建自己的域名规则
	
```sh
kubectl apply -f deploy/auto-ingress.yaml
```

annotations 中定义的标签可以认为是公共标签， 最终将被继承到生成的 Ingress 中。 因此可以通过 annotation 选择的 IngressController， 并为该 Controller 配置一些公共标签。

配置文件如下
```yml
apiVersion: apps.ingress.com/v1
kind: AutoIngress
metadata:
  name: autoingress-sample
  annotations:
    kubernetes.io/ingress.class: nginx
spec:
  rootDomain: example.cn
  #servicePrefixes:
  #  - "srv"
  #  - "web"
  namespaces:
    - "infra--staging"
  namespaceblacklist:
    - "kube-system"
    - "ingress-nginx"
  tlsSecretName: "tls-test"
```
+ `rootDomain`:         (必须), 后缀域名, 必须。
+ `servicePrefixes`:    (可选), 指定适配以 特定 关键字开头的 service,不指定匹配所有。
+ `namespaces`:         (可选), 只匹配特定namespace 下的service，不指定匹配所有namespace下svc。
+ `namespaceblacklist`: (可选), 适配特定namespace 下的service不生成ingress
+ `tlsSecretName`: （可选） 指定使用的 https 证书在 k8s 集群中的名字。

```sh
#cat auto-ingress.yaml

apiVersion: apps.ingress.com/v1
kind: AutoIngress
metadata:
  name: autoingress-sample
  annotations:
    kubernetes.io/ingress.class: nginx
spec:
  rootDomain: example.cn
  #servicePrefixes:
  #  - "srv"
  #  - "web"
  namespaces:
    - "infra--staging"
  namespaceblacklist:
    - "kube-system"
    - "ingress-nginx"

 
#kg svc -n infra--staging

NAME          TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
simple-demo   ClusterIP   172.16.106.30   <none>        80/TCP    89d

#kg ingress -n infra--staging

NAME                              CLASS    HOSTS                                     ADDRESS   PORTS   AGE
simple-demo--autoingress-sample   <none>   simple-demo---infra--staging.example.cn             80      42h

```