# Nomad Dev Cluster

## Pre-requisites

See the [ansible/README](../ansible/README.md) in the root of this repository for instructions on how to set up the environment.

## Configuration

Rename [devnet.hcl.tmp](./devnet.hcl.tmp) file to `devnet.hcl` and update the values as needed.

> Note: to uniquely identify log messages in the DadaDog logs, change the `DATA_DOG_VERSION_TAG` value in the `devnet.hcl` file.

## Running the Cluster

To start the Nomad cluster `integration` with devnet.hcl profile, run the following command where the `cluster-cli.sh` script is located:

```shell
./cluster-cli.sh start integration devnet.hcl
```

To stop the running Nomad cluster run the following command where the `cluster-cli.sh` script is located:

```shell
./cluster-cli.sh stop -purge
```

For more information on the `cluster-cli.sh` script, run the following command:

```shell
./cluster-cli.sh help
```

## Using Custom Binaries

To use or test custom binaries for different artifacts during development process, change the `source` value of the specific `artifact` stanza in the corresponding nomad script.