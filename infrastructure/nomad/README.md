# Nomad Cluster

This guide describes the steps to install (and manage) Nomad and its dependencies on multiple machines using Ansible.
It covers setting up Nomad servers and clients according to the configuration specified in your `hosts.ini` file.

## Prerequisites

On your control machine:
- [AWS CLI tools](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
- [GoReleaser](https://goreleaser.com/install/)
- [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html)
- [Ansible's collections installed](https://docs.ansible.com/ansible/latest/collections_guide/collections_installing.html):
  ```shell
  ansible-galaxy collection install community.aws
  ansible-galaxy collection install community.general
  ```
- SSH access configured for the target machines listed in `hosts.ini`.
- Target machines have a user named ubuntu or adjust the ansible_user as per your setup.
- SSH keys are set up for authentication, or alternatively, you can use SSH passwords (ensure sshpass is installed if using passwords).

## Configuration

Ensure your AWS CLI is configured with the necessary credentials:
```shell
aws sts get-caller-identity
```

The output should display your AWS account ID, ARN, and user ID similar to the following:
```json
{
    "UserId": "AIDACKCEVSQ6C2EXAMPLE",
    "Account": "123456789012",
    "Arn": "arn:aws:iam::123456789012:user/Alice"
}
```

If you see the following error message your AWS CLI is not configured correctly:
```text
Unable to locate credentials. You can configure credentials by running "aws configure"
```

To configure your AWS CLI, create AWS access keys by following the instructions in the [AWS documentation](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html#Using_CreateAccessKey).
After creating the access keys, configure your AWS CLI with the access key ID, secret access key, and default region to `us-west-2` (optionally output format of your choice) by running the following command:
```shell
aws configure
```

Add your private key to the SSH agent:
```shell
ssh-add /path/to/your/private_key
```

Prepare `ansible.cfg` File: This file contains the configuration settings for Ansible and is also used to specify the location of the inventory file.
There is a sample `ansible.cfg.example` file located in this directory. Copy this file if you're happy with the default settings or create your own.

Prepare `hosts.ini` File: This file contains the IP addresses of your Nomad servers and clients. Replace the sample IP addresses with the actual IP addresses of your machines.
```ini
[nomad_servers]
192.0.2.1 ansible_user=ubuntu
192.0.2.2 ansible_user=ubuntu
192.0.2.3 ansible_user=ubuntu

[nomad_clients]
198.51.100.1 ansible_user=ubuntu
198.51.100.2 ansible_user=ubuntu
198.51.100.3 ansible_user=ubuntu
```
> Replace the 192.0.2.X and 198.51.100.X with the IP addresses of your Nomad server and client machines, respectively.
> Ensure the ansible_user matches the username on your target machines that has SSH access.
> If your host machine is the same as your control machine replace the IP addresses with `127.0.0.1 ansible_connection=local`.

If you do not want to use the SSH agent, another option is to add the following configuration to every `nomad_server` or
`nomad_client` record in the `host.ini` file: `ansible_ssh_private_key_file=/path/to/your/private_key`. For example:
```ini
[nomad_servers]
192.0.2.1 ansible_user=ubuntu ansible_ssh_private_key_file=/path/to/your/private_key

[nomad_clients]
198.51.100.1 ansible_user=ubuntu ansible_ssh_private_key_file=/path/to/your/private_key
198.51.100.2 ansible_user=ubuntu ansible_ssh_private_key_file=/path/to/your/private_key
```

Or if the private key is the same for all machines, you can set the `ansible_ssh_private_key_file` in the `[all:vars]` section:
```ini
[all:vars]
ansible_ssh_private_key_file=/path/to/your/private_key
```

Check ansible inventory:
```shell
ansible-inventory --inventory path/to/your/hosts.ini --list --yaml
```

Test the connection to the target machines:
```shell
ansible all --inventory path/to/your/hosts.ini --module-name ping
```

> The `ansible.cfg.example` refers to the `hosts.ini` file in the same directory.
> If you have a different location for your inventory file, update the `inventory` setting in your `ansible.cfg` file.

## Cluster Management

To manage the Nomad cluster, use the `cluster.sh` script. This script allows you to initialize, deploy, and destroy the Nomad cluster.
For detailed usage instructions on how to use the script, run the following command:
```shell
./cluster.sh --help
```

Example of initializing Nomad cluster:
```shell
./cluster.sh init
```

> The `init` should be run only once, before deploying the Nomad cluster.

Example of destroying running Nomad cluster:
```shell
./cluster.sh destroy
```

Example of deploying Nomad cluster:
```shell
./cluster.sh deploy --no-logs-collection --datadog-key <DATADOG_API_KEY> --l1-rpc-url <L1_RPC_URL>
```

> If only artifacts for a particular architecture and operating system are needed, you 
> can specify them using the following environment variables ARTIFACTS_GOOS and ARTIFACTS_GOARCH, respectively.
> This will reduce the number of artifacts generated and thus reduce deployment time.
>
> Example:
> ```shell
> export ARTIFACTS_GOOS=linux
> export ARTIFACTS_GOARCH=amd64
> ```

> After successful deployment, you can find all artifacts (including secrets in the secrets.json file) at: http://<TARGET_MACHINE_IP>:1111
