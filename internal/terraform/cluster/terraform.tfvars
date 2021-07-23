environment                  = "dev"
project                      = "ozzy-playground"
cluster_zone                 = "europe-west2-b"
cluster_state_bucket         = "ozzy-playground-tfstate-dev"
state_bucket_prefix          = "terraform/state/cluster"
applications_machine_type    = "n1-standard-2"
applications_max_node_count  = 3
k8s_master_version           = "1.17.15-gke.800"
k8s_node_pool_version        = "1.17.15-gke.800"
//terraform apply -var-file="cluster.tfvars" 
//https://github.com/jetstack/terraform-google-gke-cluster/blob/master/example/main.tf
