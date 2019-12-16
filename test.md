
# My Module

This is a test module and explains main documentation

## Example

```hcl
module "example" {
    source = "_example/"
    
    aboolean =  
    
    amis =  {
        ap-northeast-1 = "ami-095dbf68",
        ap-southeast-1 = "ami-cf03d2ac",
        ap-southeast-2 = "ami-697a540a",
        eu-central-1 = "ami-b0cc23df",
        eu-west-1 = "ami-4e6ffe3d",
        us-east-1 = "ami-8f7687e2",
        us-west-1 = "ami-bb473cdb",
        us-west-2 = "ami-84b44de4",
    }
    
    required_list = [
    ]
    
    richobject = {
        value = "somevalue"
        test =  123
    }
    
    security_group_ids = "sg-a, sg-b"
    
    something_list = [
        "abc",
    ]
    
    subnet_ids = ""
}
```

## Constraints

Terraform required version [~&gt; 0.12]

### Providers

* google
* google-beta
* test

## Variables

| Name | Type | Default | Required | Description |
| ---- | ---- | ------- | -------- | ----------- |
| aboolean | bool |   | **yes** | A list |
| required_list | list | [] | **yes** | A list |
| subnet_ids | string | "" | **yes** | a comma-separated list of subnet IDs<br>  asd |
| amis | map | {<br>&nbsp;&nbsp;&nbsp;&nbsp;ap-northeast-1 = "ami-095dbf68",<br>&nbsp;&nbsp;&nbsp;&nbsp;ap-southeast-1 = "ami-cf03d2ac",<br>&nbsp;&nbsp;&nbsp;&nbsp;ap-southeast-2 = "ami-697a540a",<br>&nbsp;&nbsp;&nbsp;&nbsp;eu-central-1 = "ami-b0cc23df",<br>&nbsp;&nbsp;&nbsp;&nbsp;eu-west-1 = "ami-4e6ffe3d",<br>&nbsp;&nbsp;&nbsp;&nbsp;us-east-1 = "ami-8f7687e2",<br>&nbsp;&nbsp;&nbsp;&nbsp;us-west-1 = "ami-bb473cdb",<br>&nbsp;&nbsp;&nbsp;&nbsp;us-west-2 = "ami-84b44de4",<br>} | no | This is a super long description.<br>  <br>  Using heredoc-syntax it is possible to<br>  create simple multiline description inside<br>  terraform.<br>  <br>  With this we can create a much more meaningful <br>  documentation of specific variables in case there<br>  is the need to describe a lot of stuff<br>  before we can actually use it. |
| richobject | object | {<br>&nbsp;&nbsp;&nbsp;&nbsp;value = "somevalue"<br>&nbsp;&nbsp;&nbsp;&nbsp;test =  123<br>} | no | more Fun |
| security_group_ids | string | "sg-a, sg-b" | no | anitgher amore |
| something_list | list | [<br>&nbsp;&nbsp;&nbsp;&nbsp;"abc",<br>] | no | A list |

## Outputs

| Name | Description |
| ---- | ----------- |
| vpc_id | vpc output desc |

## Resources

### Data

| ID  | Type | Name | Provider |
| --- | ---- | ---- | -------- |
| data.google_compute_zones.zones | google_compute_zones | zones | google |

### Managed

| ID  | Type | Name | Provider |
| --- | ---- | ---- | -------- |
| google_compute_instance.someinstance | google_compute_instance | someinstance | test |

## Modules

| Name  | Version | Source |
| ----- | ------- | ------ |
| somemodule | - | ./module/dir |

