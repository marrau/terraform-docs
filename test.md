# _example


Terrform required version ~&gt; 0.12


## Providers

| Name | Alias | Version |
|------|-------|---------|
| [google](https://www.terraform.io/docs/providers/google/index.html) | test |  |
| [google-beta](https://www.terraform.io/docs/providers/google/index.html) | test |  |

## Modules

| Name | Version | Source |
|------|-------------|--------|
| somemodule |  | ./module/dir |

## Resources

| Name | Type |
|------|------|
| someinstance | [google_compute_instance](https://www.terraform.io/docs/providers/google/r/compute_instance.html) |
| zones | [google_compute_zones](https://www.terraform.io/docs/providers/google/d/compute_zones.html) |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| subnet_ids | a comma-separated list of subnet IDs | string |  | **yes** |
| amis | more things | map | map[ap-northeast-1:ami-095dbf68 ap-southeast-1:ami-cf03d2ac ap-southeast-2:ami-697a540a us-east-1:ami-8f7687e2 us-west-1:ami-bb473cdb us-west-2:ami-84b44de4 eu-west-1:ami-4e6ffe3d eu-central-1:ami-b0cc23df] | no |
| security_group_ids | anitgher amore | string | sg-a, sg-b | no |

## Outputs

| Name | Description |
|------|-------------|
| vpc_id | vpc output desc |

