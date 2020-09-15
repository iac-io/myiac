# Accessing Kubernetes cluster (apiserver) from outside GKE

By default all access outside `kubectl` and GKE dashboard is disabled. Cannot do `curl` or access
the Kubernetes API server in browser at all. This is on purpose as the Kubernetes UI Dashboard has a
highly privileged administration role assigned by default.

In order to enable this and become a cluster-admin with a specific service account follow these steps.

## Create a service account

`kubectl create serviceaccount service-account-name
`

## Ensure your user can create RBAC resources

```
ACCOUNT=$(gcloud info --format='value(config.account)')
$ kubectl create clusterrolebinding owner-cluster-admin-binding \
    --clusterrole cluster-admin \
    --user $ACCOUNT
```

Without this, creating Roles/ClusterRoles/RoleBindings/ClusterRoleBindings may give you errors.

## Create a specific cluster role with limited permissions

For example, this ClusterRole allows getting, watching and listing pods and logs:

```
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: cluster-role-name
rules:
- apiGroups: [""] # "" indicates the core API group
  resources: ["pods", "pods/log"]
  verbs: ["get", "watch", "list"]
```

If you intend to become a cluster admin (unlimited access to everything), there is a existing clusterRole `cluster-admin` already. This is not advisable in general, but if this is a personal cluster or dev one, it can be useful.

## Bind you service account to the cluster role

```
kubectl create rolebinding service-account-name:cluster-role-name --clusterrole cluster-role-name --serviceaccount default:service-account-name
```

where `default` is the namespace, and those `kubectl` act in it by default.

For example, to become `cluster-admin`with moneycol service account:

```
kubectl create rolebinding moneycol-service-account:cluster-admin --clusterrole cluster-admin --serviceaccount default:moneycol-service-account
```

# Authenticate 

Use `authenticate-service-account.sh [service-account-name]`. Once done, curl can be used to browse api.

Links:

- https://medium.com/@nieldw/curling-the-kubernetes-api-server-d7675cfc398c
- https://medium.com/@lestrrat/accessing-the-kubernetes-api-sans-the-proxy-b24af1eb18a4
- https://medium.com/@lestrrat/configuring-rbac-for-your-kubernetes-service-accounts-c348b64eb242
