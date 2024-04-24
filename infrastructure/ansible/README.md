# Installing Nomad on Server and Client Machines

This guide outlines the steps to install Nomad on multiple machines using Ansible. It covers setting up Nomad servers and clients according to the configuration specified in your config.ini file.

## Prerequisites
- Ansible installed on your control machine.
- Ansible's collections installed:
  ```shell
  ansible-galaxy collection install community.aws
  ansible-galaxy collection install community.general
  ```
- SSH access configured for the target machines listed in hosts.ini.
- Target machines have a user named ubuntu or adjust the ansible_user as per your setup.
- SSH keys are set up for authentication, or alternatively, you can use SSH passwords (ensure sshpass is installed if using passwords).

## Configuration

Prepare `config.ini` File: This file contains the IP addresses of your Nomad servers and clients. Replace the sample IP addresses with the actual IP addresses of your machines.

```
[local]
127.0.0.1 environment=devnet

[nomad_servers]
192.0.2.1 ansible_user=ubuntu
192.0.2.2 ansible_user=ubuntu
192.0.2.3 ansible_user=ubuntu

[nomad_clients]
198.51.100.1 ansible_user=ubuntu
198.51.100.2 ansible_user=ubuntu
198.51.100.3 ansible_user=ubuntu
```

> If you run ansible on your local machine add the following to the `[local]` section of your `config.ini` file: 127.0.0.1 ansible_connection=local

- Replace the 192.0.2.X and 198.51.100.X with the IP addresses of your Nomad server and client machines, respectively.
- Ensure the ansible_user matches the username on your target machines that has SSH access.

## Running the Playbook

Execute the Ansible playbook to install Nomad on the specified servers and clients. Navigate to the directory containing your playbook and run the following command:

```shell
cd ansible
ansible-playbook -i config.ini playbooks/nomad/init.yml --private-key=~/.ssh/your_private_key
```

- Replace ~/.ssh/your_private_key with the path to your SSH private key if not using the default SSH key.
- If you prefer to use SSH agent forwarding instead of directly specifying a private key, you can omit the --private-key option, assuming your SSH agent is running and loaded with your keys.

To install [crisis tools](https://www.brendangregg.com/blog/2024-03-24/linux-crisis-tools.html) run the following playbook:

```shell
ansible-playbook -i config.ini playbooks/nomad/install_linux_crisis_tools.yml --private-key=~/.ssh/your_private_key
```

And finally if you need certificates for development purposes, run the following playbook:

```shell
ansible-playbook -i config.ini playbooks/nomad/install_dev_ss_certificates.yml --private-key=~/.ssh/your_private_key
```