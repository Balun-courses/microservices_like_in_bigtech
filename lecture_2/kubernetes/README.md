# Деплой в kubernetes

## Установка Minikube
1. [Kubernetes: установка Minikube](https://kubernetes.io/ru/docs/tasks/tools/install-minikube/)
2. [Установка Minikube MacOS M1/M2](https://devopscube.com/minikube-mac/)

## Установка kubectl
1. https://kubernetes.io/ru/docs/tasks/tools/install-kubectl/
2. _ОПЦИОНАЛЬНО_: [quick kubernetes namespace and context switcher](https://github.com/blendle/kns)

## Deployment
1. `eval $(minikube docker-env)` - Вам нужно запустить `eval $(minikube docker-env)` на каждом терминале, который вы хотите использовать, поскольку он устанавливает только переменные среды для текущего сеанса оболочки.
2. `minikube image load main` - добавляем наш образ в minikube
3. `minikube image ls --format table` - просмотр образ в minikube (дожен быть docker.io/library/main)
4. `kubectl --context minikube apply ./namespace.yaml` - создаем namespace
5. `kubectl --context minikube --namespace main apply -f ./deployment.yaml` - создаем деплой
6. `kubectl --context minikube --namespace main get deploy`
7. `kubectl --context minikube --namespace main get pods --show-labels`
8. `kubectl --context minikube --namespace main get rs`
9. `kubectl --context minikube --namespace main set image deployment/main main-container=docker.io/library/main:latest` - устаналвиаем новый образ

## Service

### ClusterIP
1. `kubectl --context minikube --namespace main apply -f ./service_cluster_ip.yaml`
2. `kubectl --context minikube --namespace main get svc`
3. `kubectl --context minikube --namespace main port-forward service/main 8080:80`

### NodePort
1. `kubectl --context minikube --namespace main apply -f ./service_node_port.yaml`
2. `kubectl --context minikube --namespace main get svc`
3. `kubectl --context minikube --namespace main get nodes -o wide`
4. `minikube service main -n main` - быстрый путь :)

### Load
1. `kubectl --context minikube --namespace main apply -f ./service_load_balancer.yaml`
2. `kubectl --context minikube --namespace main get svc`
3. `minikube tunnel`
4. `kubectl --context minikube --namespace main get svc`

## Ingress
1. `minikube addons enable ingress`
2. `kubectl get pods -n ingress-nginx`
3. `kubectl apply -f ./ingress.yaml`
4. `kubectl get ingress`
5. `curl --resolve "main.info:80:$( minikube ip )" -i http://main.info`
6. `sudo -- sh -c "echo $( minikube ip ) main.info >> /etc/hosts"`

## Probes
1. `kubectl --context minikube --namespace main apply -f ./deployment_with_probes.yaml`
2. `kubectl --context minikube --namespace main describe deployment main`

# Requests and Limits
1. `kubectl --context minikube --namespace main apply -f ./deployment_with_requests_and_limits`
2. `kubectl --context minikube --namespace main describe deployment main`

# Полезные ссылки
1. https://codefresh.io/learn/kubernetes-deployment/
2. https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/
3. https://kubernetes.io/docs/concepts/services-networking/gateway/
4. https://gateway-api.sigs.k8s.io/guides/
5. https://alankrantas.medium.com/trying-out-kubernetes-gateway-api-beta-using-contour-with-kind-b5a6491096c1
6. https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
