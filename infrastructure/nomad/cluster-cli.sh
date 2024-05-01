#!/usr/bin/env bash

# Function to display help information.
show_help() {
  echo "Usage:"
  echo "  $0 [command] [options]"
  echo
  echo "Commands:"
  echo "  start <testnet|integration> <path/to/profile_file>  Start the cluster with the given profile file"
  echo "  stop [-purge]                                       Stop the cluster. Optional [-purge] flag to remove all data"
  echo "  scale <job_name> <1-N>                              Scale the specified job up or down in the running cluster to the specified number of instances"
  echo
  echo "Job names for scale command:"
  echo "  geth-bootnode"
  echo "  geth-signer-node"
  echo "  geth-member-node"
  echo
  echo "Examples:"
  echo "  $0 start /path/to/profile_file  Start the cluster using the specified profile file"
  echo "  $0 stop                         Stop the cluster"
  echo "  $0 scale geth-member-node 3     Scale the 'geth-member-node' job to 3 instances in the running cluster"
}

# Function to start the cluster.
start_cluster() {
  echo "Starting cluster..."

  set -e
  nomad system gc
  nomad var put -in hcl @"$2"
  case "$1" in
    testnet)
        nomad run datadog-agent-logs-collector.nomad
        nomad run mev-commit-geth-bootnode1.nomad
        nomad run mev-commit-geth-signer-node1.nomad
        nomad run mev-commit-geth-member-node.nomad
        nomad run deploy-contracts.nomad && sleep 60
        nomad run mev-commit-bootnode1.nomad
        nomad run mev-commit-provider-node1.nomad
        nomad run mev-commit-provider-node1-funder.nomad
        nomad run mev-commit-provider-emulator-node1.nomad
        nomad run mev-commit-oracle.nomad
        nomad run mev-commit-bridge.nomad
        nomad run datadog-agent-metrics-collector.nomad
      ;;
    integration)
        nomad run datadog-agent-logs-collector.nomad
        nomad run mev-commit-geth-bootnode1.nomad
        nomad run mev-commit-geth-signer-node1.nomad
        nomad run mev-commit-geth-member-node.nomad
        nomad run deploy-contracts.nomad && sleep 60
        nomad run mev-commit-bootnode1.nomad
        nomad run mev-commit-provider-node1.nomad
        nomad run mev-commit-provider-node1-funder.nomad
        nomad run mev-commit-bidder-node1.nomad
        nomad run mev-commit-bidder-node1-funder.nomad
        nomad run mev-commit-provider-emulator-node1.nomad
        nomad run mev-commit-bidder-emulator-node1.nomad
        nomad run mev-commit-oracle.nomad
        nomad run mev-commit-bridge.nomad
        nomad run datadog-agent-metrics-collector.nomad
      ;;
  esac
  set +e

  echo "All jobs started successfully"
}

# Function to stop the cluster.
stop_cluster() {
  echo "Stopping cluster..."

  local nomad_opts
  if [ -n "$1" ]; then
    nomad_opts="-purge"
    nomad var purge "nomad/jobs"
  fi
  for job in $(nomad job status | awk 'NR > 1 {print $1}'); do
      nomad stop $nomad_opts "$job"
  done

  echo "All jobs have been stopped and purged successfully"
}

# Function to scale a job up/down.
scale_job() {
  echo "Scaling $1 to $2 instance(s)..."

  case "$1" in
    geth-member-node)
      nomad job scale "mev-commit-geth-member-node" "member-node-group" "$2"
      ;;
    *)
      echo "Error: scaling of $1 is not supported yet'" >&2
      exit 1
      ;;
  esac

  result=$?
  if [ $result -ne 0 ]; then
    echo "Scaling $1 to $2 instances has failed"
    exit $result
  fi

  echo "The $1 is now running $2 instance(s)"
}

case "$1" in
  start)
    case "$2" in
      testnet|integration)
        start_cluster "$2" "$3"
        ;;
      *)
        echo "Error: invalid parameter, use 'testnet', or 'integration'" >&2
        show_help
        exit 1
        ;;
    esac
    ;;
  stop)
    if [ -n "$2" ] && [ "$2" != "-purge" ]; then
      echo "Error: invalid flag, use '-purge'" >&2
      show_help
      exit 1
    fi
    stop_cluster "$([ -n "$2" ] && echo true || echo false)"
    ;;
  scale)
     case "$2" in
        geth-bootnode|geth-signer-node|geth-member-node)
          if [ -z "$3" ] || echo "$3" | grep -qvE '^-?[0-9]+(\.[0-9]+)?$'; then
            echo "Error: invalid number of instances: '$3'" >&2
            show_help
            exit 1
          fi
          scale_job "$2" "$3"
          ;;
        *)
          echo "Error: invalid job name, use 'geth-bootnode', 'geth-signer-node', or 'geth-member-node'" >&2
          show_help
          exit 1
          ;;
      esac
    ;;
  help|*)
    echo "Error: invalid command, use 'start', 'stop', or 'scale'" >&2
    show_help
    ;;
esac
