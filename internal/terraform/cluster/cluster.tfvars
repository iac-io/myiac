environment                  = "dev"
project                      = "ozzy-playground"
cluster_zone                 = "europe-west1-b"
cluster_state_bucket         = "ozzy-playground-tfstate-dev"
state_bucket_prefix          = "terraform/state/cluster"
applications_machine_type    = "n1-standard-2"
elasticsearch_machine_type   = "n1-standard-1"
applications_max_node_count  = 3
elasticsearch_max_node_count = 1
//terraform apply -var-file="cluster.tfvars" 
//https://github.com/jetstack/terraform-google-gke-cluster/blob/master/example/main.tf
