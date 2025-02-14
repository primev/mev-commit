# mev-commit

[![CI](https://github.com/primev/mev-commit/actions/workflows/ci.yml/badge.svg)](https://github.com/primev/mev-commit/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-BSL%201.1-blue.svg)](LICENSE)

Primev's mev-commit makes Ethereum FAST by enabling preconfirmations and new types of mev using an encrypted mempool, enhancing yield for opted-in validators.

## Documentation

For detailed documentation, visit the [mev-commit docs](https://docs.primev.xyz/).

## Main Components
  - [mev-commit client](p2p)
  - [mev-commit-oracle](oracle)
  - [mev-commit-bridge](bridge)
  - [mev-commit-geth](external/geth)
  - [contracts](contracts)

## Getting Started

The mev-commit repository is a mono-repository that also use submodules for external dependencies.
To clone the repository and its submodules, run the following command:

```shell
git clone --recurse-submodules <git@github.com:primev/mev-commit.git|https://github.com/primev/mev-commit.git>
```

If you have already cloned the repository and need to update the submodules, run the following command:

```shell
git submodule update --init --recursive
```

## Development

When working with submodules, you can use the `git submodule` command to list available submodules.
To make changes to a submodule, navigate to the submodule directory using the `cd` command.
Before making any changes, ensure that the submodule is up-to-date.
For example, to navigate to the `external/geth` submodule, run the following command:

```shell
cd external/geth
git submodule update --init --recursive
```

Make the necessary changes to the submodule and commit (and push) them.
To make the changes available in the main repository, you need to push the changes to the submodule.
After making changes to the submodule, navigate back to the main repository and commit the changes.
For example, to commit and push changes made to the `external/geth` submodule, run the following commands inside the submodule directory:

```shell
git add -p
git commit -m "<your-commit-message>"
git push
```
Go back to the main repository and commit (and push) the changes made to the submodule:

```shell
git add external/geth
git commit -m "<your-commit-message>"
git push
```

### Go Modules

Since this repository uses [Go Workspaces](https://go.dev/ref/mod#workspaces) to manage Go modules, when making changes to a Go module and its dependencies, ensure that the changes are reflected everywhere by running the following command:

```shell
go list -f '{{.Dir}}' -m | xargs -L1 go mod tidy -C
go work sync
```

> See the [go.work](go.work) file for all the Go modules used in this repository.
