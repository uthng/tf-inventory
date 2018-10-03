# tf-inventory
**tf-inventory** is a little tool to generate Ansible dynamic inventory file from Terraform state file. It is inspired of the work of adammck: [terraform-inventory](https://github.com/adammck/terraform-inventory).

**tf-inventory** parses the resources in each module. If the resource type is supported, it will get the `hostname` and `ip` from the instance. In Ansible, `hostname` and `ip` will be used respectively as the host's name and `ansible_ssh_host`. The name of supported resource will be the group name even in the case of multiple instances of a resource (using `count`).

**tf-inventory** considers Terraform outputs as host variables in Ansible. If the output is a `map`, then each element of the map will be considered as a host var distinct.

## Installation
To use **tf-inventory**, you have to clone this repository and build it as following:
```shell
# git clone https://github.com/uthng/tf-inventory
# glide install
# go build -o tf-inventory
```

## Usage
You can use **tf-inventory** alone or with ansible.
#### CLI
```shell
./tf-inventory is a utility to generate ansible inventory file from Terraform state
Syntax: ./tf-inventory [options] [hostname] [tfstate]
  -host string
    	Host mode
  -list
    	List mode
```
If the `tfstate` is omitted,  **tf-inventory** will try to read `terraform.tfstate` in the current directory. If no file is found, it will try to execute the command `terraform state pull` to get using the information specified in the 2 following environment variables:
- **TF_BIN**: path to folder containing `terraform` binary
- **TF_DIR**: path to Terraform working folder.

After all, it will return an error.

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