provider "google-beta" {
  credentials = "${file("/Users/dfernandez/account.json")}"
  project     = "moneycol"
  region      = "europe-west1-b"
}

terraform {
  backend "gcs" {
    bucket = "moneycol-tf-state-dev"
    prefix = "terraform/state/dns"
  }
}

variable "dev_ip" {
  default = "35.187.176.24"
}

resource "google_dns_managed_zone" "moneycol_zone_free_domain" {
  description = "Managed by Terraform"
  dns_name    = "moneycol.ml."
  name        = "money-zone-free"
  project     = "moneycol"
  visibility  = "public"

  dnssec_config {
    kind          = "dns#managedZoneDnsSecConfig"
    non_existence = "nsec3"
    state         = "off"

    default_key_specs {
      algorithm  = "rsasha256"
      key_length = 2048
      key_type   = "keySigning"
      kind       = "dns#dnsKeySpec"
    }
    default_key_specs {
      algorithm  = "rsasha256"
      key_length = 1024
      key_type   = "zoneSigning"
      kind       = "dns#dnsKeySpec"
    }
  }

  timeouts {}
}

# terraform import google_dns_record_set.frontend {{project}}/{{zone}}/{{name}}/{{type}}
# terraform import google_dns_record_set.moneycol_dev "moneycol/money-zone-free/dev.moneycol.ml./A"
# terraform state show google_dns_record_set.moneycol_dev
resource "google_dns_record_set" "moneycol_dev" {
  managed_zone = "${google_dns_managed_zone.moneycol_zone_free_domain.name}"
  name         = "dev.${google_dns_managed_zone.moneycol_zone_free_domain.dns_name}"
  project      = "moneycol"
  rrdatas = [
    "${var.dev_ip}",
  ]
  ttl  = 30
  type = "A"
}
