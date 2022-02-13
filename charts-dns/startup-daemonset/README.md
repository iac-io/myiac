## Startup daemonset

**Note**: the use of `traefik-dev` is now deprecated due to GKE Admission Controller 
[here](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#denyserviceexternalips)


This is a DaemonSet (a service/daemon that runs in each node of a Kubernetes cluster) used for syncing the IP addresses 
of a public GKE cluster with a DNS provider (Cloudflare). This is done so that the use of Load Balancers is avoided in exchange for 
potential (unlikely) downtime when a specific node is down.

### How it works

This setup assumes that the `traefik` Ingress Controller is deployed in the cluster (see `charts/traefik` in the 
charts repository).

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
- Finds the **external IP** of the cluster node that holds the pod for Traefik Ingress Controller
- Only this IP address can be used as the entrypoint for the cluster, so there can be intermittent downtime
while the node goes down and `traefik` pod gets rescheduled in another
- Cloudflare is updated for every subdomain in the zone `moneycol.net` with it

The operation runs every 6h, and on any node startup (DaemonSet). There's still downtime, as if the selected node whose IP
address has been set in the DNS provider goes down, that IP won't be updated until the next Node boots. 
Further work could be done to mitigate this, as the fact could be detected by running this process again when termination event is received on the VM/GKE.


### How to configure and deploy

For `startup-daemonset` to run the following is required:

- A API key and email address for Cloudflare
- A Kubernetes secret present in the cluster, holding the Cloudflare API key

```
myiac createSecret --secretName cloudflare-api-key --literal CF_API_KEY=xxx
```

- The environment variables `CF_API_KEY` and `CF_EMAIL` should be populated in the container running inside the DaemonSet
- A Docker image must exist / be built using the `Dockerfile` present at `Dockerfiles/Dockerfile`. The lines:
```
COPY --from=myiac-builder /workdir/*.json /home/app/
COPY --from=myiac-builder /workdir/charts /home/app/charts
```

aren't generic enough and assume that a) there's a JSON service account key on the root folder b) there's a `charts` folder
in the root as well.
- The `CMD` on this `myiac` container must be `entrypoint.sh`
- To build/tag the container correctly, the script `build-myiac-moneycol-docker.sh`should be used, setting up the 
adequate version to push inside:

```
./build-myiac-moneycol-docker.sh
```

The above will push a Docker image to the registry. The built tag should be used in the subsequent deploy

- Deploy the `startup-daemonset` with the built image

```
myiac setupEnvironment --provider gcp --project moneycol --keyPath ./account.json --zone europe-west1-b --env dev 
myiac deploy --project moneycol --env dev --app startup-daemonset --properties image.tag=<built-tag>
```

or 

```
# edit charts-dns/startup-daemonset/values.yaml, place the <tag> built
helm install startup-daemonset charts-dns/startup-daemonset
```