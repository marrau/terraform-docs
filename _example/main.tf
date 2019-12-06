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

  required_version = "~> 2.20"
}

provider "google-beta" {
  alias = "test"
}

// liasdhfjasodifuh
variable "subnet_ids" {
  description = "a comma-separated list of subnet IDs"
  type        = string
}

variable "security_group_ids" {
  description = "anitgher amore"
  default     = "sg-a, sg-b"
}

variable "something_list" {
  description = "A list"
  default     = ["abc"]
}

variable "required_list" {
  description = "A list"
  type        = list
}

variable "aboolean" {
  description = "A list"
  type        = bool
}


variable "richobject" {
  description = "more Fun"
  type = object({
    value = string
    test  = number
  })
  default = {
    value = "somevalue"
    test  = 123
  }
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

resource "google_compute_instance" "someinstance" {
  provider = test
}

module "somemodule" {
  source = "./module/dir"
}

locals {
  test = var.security_group_ids
}


// The VPC ID.
output "vpc_id" {
  description = "vpc output desc"
  value       = "vpc-5c1f55fd"
}
