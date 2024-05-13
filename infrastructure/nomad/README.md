# Installing Nomad on Server and Client Machines

This guide outlines the steps to install Nomad on multiple machines using Ansible.
It covers setting up Nomad servers and clients according to the configuration specified in your hosts.ini file.

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
- SSH access configured for the target machines listed in hosts.ini.
- Target machines have a user named ubuntu or adjust the ansible_user as per your setup.
- SSH keys are set up for authentication, or alternatively, you can use SSH passwords (ensure sshpass is installed if using passwords).

## Configuration

Prepare `hosts.ini` File: This file contains the IP addresses of your Nomad servers and clients. Replace the sample IP addresses with the actual IP addresses of your machines.

```
[nomad_servers]
192.0.2.1 ansible_user=ubuntu
192.0.2.2 ansible_user=ubuntu
192.0.2.3 ansible_user=ubuntu

[nomad_clients]
198.51.100.1 ansible_user=ubuntu
198.51.100.2 ansible_user=ubuntu
198.51.100.3 ansible_user=ubuntu
```

If your host machine is the same as your control machine add the following to your `hosts.ini` file:
```
[local]
127.0.0.1 ansible_connection=local
```

- Replace the 192.0.2.X and 198.51.100.X with the IP addresses of your Nomad server and client machines, respectively.
- Ensure the ansible_user matches the username on your target machines that has SSH access.

## Running the Playbook

Export the following environment variable before executing any of the playbooks below: `export OBJC_DISABLE_INITIALIZE_FORK_SAFETY=YES`

Execute the Ansible playbook to install Nomad on the specified servers and clients. From the root of the repository run the following command:

```shell
ansible-playbook -i infrastructure/nomad/hosts.ini infrastructure/nomad/init.yml --private-key=~/.ssh/your_private_key
```

- Replace ~/.ssh/your_private_key with the path to your SSH private key if not using the default SSH key.
- If you prefer to use SSH agent forwarding instead of directly specifying a private key, you can omit the --private-key option, assuming your SSH agent is running and loaded with your keys.

> The init script also installs self-signed certificates. To avoid this step, add the `--skip-tags "certs"` option to the end of the command.

To generate and run Nomad jobs run the following playbook:

```shell
ansible-playbook -i infrastructure/nomad/hosts.ini infrastructure/nomad/deploy.yml --private-key=~/.ssh/your_private_key
```

To destroy all running Nomad jobs run the following playbook:

```shell
ansible-playbook -i infrastructure/nomad/hosts.ini infrastructure/nomad/destroy.yml --private-key=~/.ssh/your_private_key
```
