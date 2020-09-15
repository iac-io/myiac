SERVICE_ACCOUNT=$1

# Get the ServiceAccount's token Secret's name
export SECRET=$(kubectl get serviceaccount ${SERVICE_ACCOUNT} -o json | jq -Mr '.secrets[].name | select(contains("token"))')
echo $SECRET

# Extract the Bearer token from the Secret and decode
export TOKEN=$(kubectl get secret ${SECRET} -o json | jq -Mr '.data.token' | base64 -D)
echo "export TOKEN=$TOKEN"

# Extract, decode and write the ca.crt to a temporary location
kubectl get secret ${SECRET} -o json | jq -Mr '.data["ca.crt"]' | base64 -D > /tmp/ca.crt

# Get the API Server location
export APISERVER=https://$(kubectl -n default get endpoints kubernetes --no-headers | awk '{ print $2 }')
echo "curl -i --cacert /tmp/ca.crt $APISERVER --header 'Authorization: Bearer $TOKEN"

# Example of proxied service
# http://localhost:8001/api/v1/namespaces/default/services/moneycol-server:80/proxy/graphql