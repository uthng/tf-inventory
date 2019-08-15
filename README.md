# tf-inventory
**tf-inventory** is a little tool to generate Ansible dynamic inventory file from Terraform state file. It is inspired of the work of adammck: [terraform-inventory](https://github.com/adammck/terraform-inventory).

**tf-inventory** parses the resources in each module. If the resource type is supported, it will get the `hostname` and `ip` from the instance. In Ansible, `hostname` and `ip` will be used respectively as the host's name and `ansible_ssh_host`. 

If Terraform module is not used in TF files (only one module with `Path` root), the name of supported resource will be the group name even in the case of multiple instances of a resource. Otherwise, servers will be grouped under the name given by the 2nd element of `Module.Path`.

**tf-inventory** considers Terraform outputs as group variables in Ansible. If the output is a `map`, then each element of the map will be considered as a group var distinct.

Actually, the following providers are supported:
- Vsphere
- Scaleway

## Installation
To use **tf-inventory**, you have to clone this repository and build it as following:
```shell
# git clone https://github.com/uthng/tf-inventory
# cd tf-inventory
# GO111MODULE=on make deps build
# ./releases/<yourplatform>/tf-inventory
```

## Usage
You can use **tf-inventory** alone or with ansible.

#### CLI
```shell
tf-inventory --help
usage: tf-inventory [<flags>] [<tfsfile>]

An utility to generate ansible inventory file from Terraform state

Flags:
  --help       Show context-sensitive help (also try --help-long and --help-man).
  --list       List mode
  --host=HOST  Host mode

Args:
  [<tfsfile>]  Terraform state file
```

**tf-inventory** gets the content of Terraform state file in 3 following ways in order:

1. `tfsfile`: if it is specified, **tf-inventory** will read and parse its contents
2. `consul`: if the following environment variables are specified, **tf-inventory** will get the tfstate content with the given key from Consul
	- **CONSUL_URL:** Consul server's url
	- **CONSUL_KEY:** Consul key containing tfstate
	- **CONSUL_TOKEN:** Consul token to access to key
3. `terraform`: if 2 above cases donot satisfy, **tf-inventory** will perform the command `terraform state pull` to generate tfstate instead using the following environment variables:
	- **TF_BIN**: path to folder containing `terraform` binary
	- **TF_DIR**: path to Terraform working folder.

Otherwise, it will return an error.

##### Examples:
Using `tfstate` argument:
```shell
./tf-inventory --list ./fixtures/terraform.tfstate
```

Using `TF_DIR`:
```shell
TF_DIR=/path/to/terraform ./tf-inventory --list
```

#### With Ansible
When using with Ansible, the current directory must be a Terraform working folder:
```shell
ansible -i /path/to/tf-inventory my-host -m ping
```
If not, `TF_DIR` must be used to specify it:
```shell
TF_DIR=/path/to/terraform/ ansible -i tf-inventory my-host -m ping
```