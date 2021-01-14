variable "project" {
  type = string
}

variable "environment" {
  type = string
}

variable "cluster_zone" {
  type = string
}

variable "cluster_state_bucket" {
  type = string
}

variable "state_bucket_prefix" {
  type = string
}

variable "applications_machine_type" {
  type = string
}

variable "applications_max_node_count" {
  type = string
}

variable "k8s_master_version" {
  type = string
}

variable "k8s_node_pool_version" {
  type = string
}