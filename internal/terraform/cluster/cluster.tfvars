environment                  = "dev"
project                      = "qumu-dev"
cluster_zone                 = "europe-west1-b"
cluster_state_bucket         = "qumu-dev-tf-state-dev"
state_bucket_prefix          = "terraform/state/cluster"
applications_machine_type    = "f1-micro"
elasticsearch_machine_type   = "n1-standard-1"
applications_max_node_count  = 3
elasticsearch_max_node_count = 2
//terraform apply -var-file="cluster.tfvars" 
//https://github.com/jetstack/terraform-google-gke-cluster/blob/master/example/main.tf
