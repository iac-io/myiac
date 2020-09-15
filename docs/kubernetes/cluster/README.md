## SSL for Free setup with Traefik

* From SSL free files bundle, concatenate `certificate.crt` + `ca-bundle.crt` into `cert.crt`
```
cat certificate.crt ca-bundle.crt > cert.crt
# edit cert.crt to add a \n after the first certificate or use: 
# https://stackoverflow.com/questions/8183191/concatenating-files-and-insert-new-line-in-between-files
```
* Rename `cert.crt` to `tls.crt`
```
mv cert.crt tls.crt
```
* Rename `private.key` to `tls.key`
```
mv private.key tls.key
```
* Run the following to create a secret in Kubernetes for the certificate
```
kubectl -n default create secret tls traefik-dev-tls-cert --key=tls.key --cert=tls.crt
```
* Reference in the Traefik Ingress for the service needing SSL:

```
apiVersion: extensions/v1beta1
kind: Ingress
...
spec:
  rules:
    - host: moneycolfrontend
    - http:
        paths:
          - path: /
            backend:
              serviceName: moneycolfrontend
              servicePort: 80
  # the secret was created with kubectl previously (see notes/readme)
  tls:
  - hosts:
    - dev.moneycol.ml
    secretName: traefik-dev-tls-cert
```
