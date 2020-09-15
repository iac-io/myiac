# add a service account within a namespace to segregate tiller
kubectl --namespace kube-system create sa tiller

# create a cluster role binding for tiller
kubectl create clusterrolebinding tiller \
    --clusterrole cluster-admin \
    --serviceaccount=kube-system:tiller

# initialized helm within the tiller service account
helm init --service-account tiller

# updates the repos for Helm repo integration
helm repo update

echo "verify helm"
# verify that helm is installed in the cluster
kubectl get deploy,svc tiller-deploy -n kube-system

# Installing Redis
# helm install stable/redis \
#     --values values/values-production.yaml \
#     --name redis-system

# download a chart locally
# helm fetch stable/redis -d redis-chart --untar
# helm install ./redis-chart/redis --values ./redis-chart/redis/values-production.yaml --name redis-system
