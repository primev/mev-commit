#!/bin/bash

set -e

NAMESPACE="default"
HELMFILE_DIR="helmfiles"
HELMFILES=("helmfile-snode-geth.yaml" "helmfile-p2p-bootnode.yaml" "helmfile-oracle.yaml")

echo "Select mode:"
echo "1. Show diff"
echo "2. Dry-run (render manifests)"
echo "3. Apply (with confirmation)"
echo "4. Cleanup (uninstall all)"
read -p "Enter choice [1-4]: " mode

check_rollout() {
  local release=$1
  echo "⏳ Checking rollout for $release..."

  # check statefulsets
  sts=$(kubectl get statefulset -n "$NAMESPACE" -l app.kubernetes.io/instance=$release -o name)
  for s in $sts; do
    kubectl rollout status "$s" -n "$NAMESPACE" --timeout=60s || {
      echo "❌ StatefulSet $s failed"
      kubectl get pods -n "$NAMESPACE" -l app.kubernetes.io/instance=$release
      kubectl logs -n "$NAMESPACE" --tail=20 -l app.kubernetes.io/instance=$release --all-containers
      exit 1
    }
  done

  # check deployments
  deps=$(kubectl get deploy -n "$NAMESPACE" -l app.kubernetes.io/instance=$release -o name)
  for d in $deps; do
    kubectl rollout status "$d" -n "$NAMESPACE" --timeout=60s || {
      echo "❌ Deployment $d failed"
      kubectl get pods -n "$NAMESPACE" -l app.kubernetes.io/instance=$release
      kubectl logs -n "$NAMESPACE" --tail=20 -l app.kubernetes.io/instance=$release --all-containers
      exit 1
    }
  done

  echo "✅ $release is healthy"
}

case $mode in
  1)
    echo "🔍 Showing diffs:"
    for file in "${HELMFILES[@]}"; do
      echo "🔸 Diff: $file"
      helmfile -f "$HELMFILE_DIR/$file" diff || true
      echo "────────────"
    done
    ;;

  2)
    echo "📦 Rendering manifests (dry-run):"
    for file in "${HELMFILES[@]}"; do
      echo "🔸 Template: $file"
      helmfile -f "$HELMFILE_DIR/$file" template
      echo "────────────"
    done
    ;;

  3)
    echo "🚀 Starting deployment..."
    for file in "${HELMFILES[@]}"; do
      name=$(basename "$file" .yaml | sed 's/helmfile-//')
      echo "🔍 Diff for $name:"
      helmfile -f "$HELMFILE_DIR/$file" diff || true
      echo
      read -p "Apply $name? (y/n): " confirm
      if [[ $confirm == "y" ]]; then
        helmfile -f "$HELMFILE_DIR/$file" apply
        check_rollout "$name"
      else
        echo "❎ Skipped $name"
      fi
    done
    ;;

  4)
    echo "🧹 Destroying all releases..."
    for file in "${HELMFILES[@]}"; do
      name=$(basename "$file" .yaml | sed 's/helmfile-//')
      echo "🔥 Uninstalling $name..."
      helmfile -f "$HELMFILE_DIR/$file" destroy || true
    done
    ;;

  *)
    echo "❌ Invalid option. Use 1-4."
    exit 1
    ;;
esac

echo "🏁 Script completed!"
