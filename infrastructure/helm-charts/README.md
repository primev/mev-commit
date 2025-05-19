
# 🚀 MEV Commit Helm Deployment

This repository contains Helm charts and a deployment workflow to manage the following components:

- `snode-geth`
- `p2p-bootnode`
- `oracle`

---

## ⚙️ Setup Overview

- All release definitions are having their own  [`helmfile.yaml`](./helmfile.yaml).
- A companion script [`deploy.sh`](./deploy.sh) orchestrates applying releases **in order**, with **delays and health checks** to ensure readiness.
- This is necessary because certain components depend on others being fully provisioned before starting.

---

## 📦 Prerequisites

Make sure the following tools are installed before proceeding:

- [Helm](https://helm.sh/docs/intro/install/)
- [Helmfile](https://github.com/helmfile/helmfile)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [helm-diff plugin](https://github.com/databus23/helm-diff):
  ```bash
  helm plugin install https://github.com/databus23/helm-diff
  ```

---

## 📝 How to Use

1. **Review and Customize Values**
   - Each release has a `values.yaml` file in its chart directory.
   - Before running the deployment, **inspect and update these values** to match your environment:
     - `mev-commit-snode-geth/values.yaml`
     - `mev-commit-p2p/values-bootnode.yaml`
     - `mev-commit-oracle/values.yaml`

2. **Run Deployment Script**

```bash
./deploy.sh
```

You will be prompted to select a mode:

```
Select mode:
1. Show diff
2. Dry-run (render manifests)
3. Apply (with confirmation)
4. Cleanup (uninstall all)
Enter choice [1-4]:
```

---

## 🧹 Cleanup

To remove all deployed resources:
Pick option `4` while running bash script
```bash
4. Cleanup (uninstall all)
```

---

## ✅ Best Practices

- Always start with `diff` or `dry-run` to verify what’s being deployed.
- When modifying `values.yaml`, commit with clear change descriptions.
- For new environments, copy existing value files and tweak overrides.

