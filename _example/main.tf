/**
 * Module usage:
 *
 *      module "foo" {
 *        source = "github.com/foo/baz"
 *        subnet_ids = "${join(",", subnet.*.id)}"
 *      }
 *
 */
terraform {
  required_version = "~> 0.12"
}

provider "google" {
  alias = "test"
}

provider "google-beta" {
  alias = "test"
}

// liasdhfjasodifuh
variable "subnet_ids" {
  description = "subnet ids"
  description = "a comma-separated list of subnet IDs"
}

variable "security_group_ids" {
  description = "anitgher amore"
  default     = "sg-a, sg-b"
}

variable "amis" {
  description = "more things"

  default = {
    "us-east-1"      = "ami-8f7687e2"
    "us-west-1"      = "ami-bb473cdb"
    "us-west-2"      = "ami-84b44de4"
    "eu-west-1"      = "ami-4e6ffe3d"
    "eu-central-1"   = "ami-b0cc23df"
    "ap-northeast-1" = "ami-095dbf68"
    "ap-southeast-1" = "ami-cf03d2ac"
    "ap-southeast-2" = "ami-697a540a"
  }
}

data "google_compute_zones" "zones" {}

resource "google_compute_instance" "someinstance" {}

module "somemodule" {
  source = "./module/dir"
}

// The VPC ID.
output "vpc_id" {
  description = "vpc output desc"
  value       = "vpc-5c1f55fd"
}
