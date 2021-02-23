## Startup daemonset

This is a DaemonSet used for the sync of DNS between a DNS provider and the Kubernetes nodes. This way it is avoided
the use of Load Balancers at the expense of temporarily losing access to the cluster.


### How it works


The `startup-daemonset` DaemonSet runs the following script every `CHECK_INTERVAL_SECONDS` (default 6h, 21600 seconds).

```
#! /bin/bash

set -o errexit
set -o pipefail
set -o nounset
myiac setupEnvironment --provider gcp --project moneycol --keyPath /home/app/account.json --zone europe-west1-b --env dev

# CF env vars from manifest
echo "CF email $CF_EMAIL"
export CHARTS_PATH=/home/app/charts
myiac updateDnsWithClusterIps --dnsProvider cloudflare --domain moneycol.net
touch /tmp/foo
echo done

```

The key operation is `updateDnsWithClusterIps`, it does the following:

- Selects Cloudflare as the DNS provider, zone `moneycol.net`
- Finds all the **internal** IPs of the Kubernetes cluster nodes
- Using those, it deploys an internal ClusterIP Service `traefik-dev` chart. This proxies access to `traefik` Ingress, giving
 access to any exposed service. It also repoints the **internal** IPs gathered in the previous step as `externalIps` 
 (this allows access to service by IP from any node, enabling the use of **any** of the nodes public IPs as the proxy, see [docs](https://kubernetes.io/docs/concepts/services-networking/service/#external-ips))
- Following up on the above, once the `traefik-dev` is deployed, **any** of the public IPs of the Kubernetes cluster nodes can be used
as A DNS record
- One of the IPs is picked up, and Cloudflare is updated for every subdomain in the zone `moneycol.net` with it

The operation runs every 6h, and on any node startup (DaemonSet). There's still downtime, as if the selected node whose IP
address has been set in the DNS provider goes down, that IP won't be updated until the next Node boots. This could be detected
by running this process again when termination event is received on the VM/GKE.


### How to configure and deploy

For `startup-daemonset` to run the following is required:

- A API key and email address for Cloudflare
- A Kubernetes secret present in the cluster, holding the API key

```
myiac createSecret --secretName cloudflare-api-key-sec --literal CF_API_KEY=xxx
```

- The environment variables `CF_API_KEY` and `CF_EMAIL` should be populated in the container running inside the DaemonSet
- A Docker image must exist / be built using the `Dockerfile` present at `Dockerfiles/Dockerfile`. The lines:
```
COPY --from=myiac-builder /workdir/*.json /home/app/
COPY --from=myiac-builder /workdir/charts /home/app/charts
```

aren't generic enough and assume that a) there's a JSON service account key on the root folder b) there's a `charts` folder
in the root as well.
- The `CMD` on this `myiac` container must be `entrypoint.sh`, provided at the root folder of the project
- To build/tag the container correctly, the script `build-myiac-moneycol-docker.sh`should be used:

```
./build-myiac-moneycol-docker.sh
```

The above will push a Docker image to the registry. The built tag should be used in the subsequent deploy

- Deploy the `startup-daemonset` with the built image

```
myiac setupEnvironment --provider gcp --project moneycol --keyPath ./account.json --zone europe-west1-b --env dev 
myiac deploy --project moneycol --env dev --app startup-daemonset --properties image.tag=<built-tag>
```