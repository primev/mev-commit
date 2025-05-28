#!/usr/bin/env python3

import subprocess
import time
import re
import argparse
import sys
from pathlib import Path

class PoADeployer:
    def __init__(self, namespace="ethereum-poa"):
        self.namespace = namespace
        self.releases = {
            'bootnode': 'poa-bootnode',
            'signer': 'poa-signer', 
            'member': 'poa-member'
        }
        
    def run_cmd(self, cmd, capture_output=True):
        """Run shell command and return result"""
        try:
            result = subprocess.run(cmd, shell=True, capture_output=capture_output, text=True)
            if result.returncode != 0 and capture_output:
                print(f"❌ Command failed: {cmd}")
                print(f"Error: {result.stderr}")
                return None
            return result.stdout if capture_output else result.returncode == 0
        except Exception as e:
            print(f"❌ Error running command: {e}")
            return None

    def wait_for_pod_ready(self, release_name, timeout=300):
        """Wait for pod to be ready"""
        print(f"⏳ Waiting for {release_name} pod to be ready...")
        cmd = f"kubectl wait --for=condition=ready pod -l app.kubernetes.io/instance={release_name} -n {self.namespace} --timeout={timeout}s"
        return self.run_cmd(cmd, capture_output=False)

    def extract_enode(self, release_name, max_retries=30):
        """Extract enode from bootnode logs"""
        print(f"🔍 Extracting enode from {release_name}...")
        
        for attempt in range(1, max_retries + 1):
            print(f"   Attempt {attempt}/{max_retries}")
            
            # Get pod name
            pod_cmd = f"kubectl get pods -n {self.namespace} -l app.kubernetes.io/instance={release_name} -o jsonpath='{{.items[0].metadata.name}}'"
            pod_name = self.run_cmd(pod_cmd)
            
            if not pod_name:
                time.sleep(10)
                continue
                
            # Get logs and extract node ID
            logs_cmd = f"kubectl logs {pod_name.strip()} -c init-nodekey -n {self.namespace}"
            logs = self.run_cmd(logs_cmd)
            
            if logs:
                node_id_match = re.search(r'Node ID: ([a-f0-9]{128})', logs)
                if node_id_match:
                    node_id = node_id_match.group(1)
                    service_name = f"{release_name}-ethereum-poa-bootnode"
                    enode = f"enode://{node_id}@{service_name}:30303"
                    print(f"✅ Found enode: {enode}")
                    return enode
                    
            time.sleep(10)
            
        print("❌ Failed to extract enode")
        return None

    def validate_values_file(self, role, values_file):
        """Validate that required values are present"""
        print(f"🔍 Validating {values_file}...")
        
        if not Path(values_file).exists():
            print(f"❌ {values_file} not found")
            return False
            
        with open(values_file, 'r') as f:
            content = f.read()
            
        # Check for required values based on role
        required_checks = {
            'bootnode': ['role:', 'chainId:'],
            'signer': ['role:', 'chainId:', 'signer:'],
            'member': ['role:', 'chainId:']
        }
        
        missing_values = []
        
        for check in required_checks.get(role, []):
            if check not in content:
                missing_values.append(check)
            elif check == 'chainId:':
                # Special check for chainId - make sure it's not empty
                chain_id_match = re.search(r'chainId:\s*(\S+)', content)
                if not chain_id_match or not chain_id_match.group(1):
                    missing_values.append('chainId (empty value)')
                    
        if missing_values:
            print(f"❌ Missing required values in {values_file}:")
            for missing in missing_values:
                print(f"   - {missing}")
            return False
            
        print(f"✅ {values_file} validation passed")
        return True

    def update_values_with_enode(self, enode):
        """Update signer and member values files with enode"""
        print("📝 Updating values files with enode...")
        
        for role in ['signer', 'member']:
            values_file = f"values-{role}.yaml"
            if not Path(values_file).exists():
                print(f"⚠️  {values_file} not found, skipping...")
                continue
                
            # Read, update, and write back
            with open(values_file, 'r') as f:
                content = f.read()
                
            # Replace placeholder with actual enode
            updated_content = content.replace('ENODE_PLACEHOLDER', enode)
            
            # Backup original
            with open(f"{values_file}.bak", 'w') as f:
                f.write(content)
                
            # Write updated content
            with open(values_file, 'w') as f:
                f.write(updated_content)
                
            print(f"✅ Updated {values_file}")

    def deploy_release(self, role, release_name):
        """Deploy a single release"""
        print(f"🚀 Deploying {role} ({release_name})...")
        
        values_file = f"values-{role}.yaml"
        
        # Validate values file first
        if not self.validate_values_file(role, values_file):
            return False
            
        # Create namespace if needed
        self.run_cmd(f"kubectl create namespace {self.namespace} --dry-run=client -o yaml | kubectl apply -f -", capture_output=False)
        
        # Deploy with helm
        cmd = f"helm upgrade --install {release_name} . -f {values_file} -n {self.namespace} --wait --timeout=5m"
        success = self.run_cmd(cmd, capture_output=False)
        
        if success:
            print(f"✅ {role} deployed successfully")
            return True
        else:
            print(f"❌ Failed to deploy {role}")
            return False

    def cleanup(self):
        """Clean up all deployments"""
        print("🧹 Cleaning up deployments...")
        
        # Get PVCs before deleting releases
        print("📝 Collecting PVCs to delete...")
        pvcs_to_delete = []
        for role, release_name in self.releases.items():
            pvc_cmd = f"kubectl get pvc -n {self.namespace} -l app.kubernetes.io/instance={release_name} -o jsonpath='{{.items[*].metadata.name}}'"
            pvcs = self.run_cmd(pvc_cmd)
            if pvcs and pvcs.strip():
                pvcs_to_delete.extend(pvcs.strip().split())
        
        # Delete releases in reverse order
        for role, release_name in reversed(list(self.releases.items())):
            print(f"   Removing {release_name}...")
            self.run_cmd(f"helm uninstall {release_name} -n {self.namespace}", capture_output=False)
            
        # Delete PVCs (StatefulSet PVCs don't get auto-deleted)
        if pvcs_to_delete:
            print(f"🗑️  Deleting PVCs: {', '.join(pvcs_to_delete)}")
            pvc_list = ' '.join(pvcs_to_delete)
            self.run_cmd(f"kubectl delete pvc {pvc_list} -n {self.namespace}", capture_output=False)
        else:
            print("📋 No PVCs found to delete")
            
        # Restore backup files
        for role in ['signer', 'member']:
            backup_file = f"values-{role}.yaml.bak"
            values_file = f"values-{role}.yaml"
            if Path(backup_file).exists():
                self.run_cmd(f"mv {backup_file} {values_file}")
                print(f"✅ Restored {values_file}")
                
        print("✅ Cleanup complete")

    def dry_run(self):
        """Show what would be deployed"""
        print("📋 Deployment Plan (Dry Run)")
        print("=" * 30)
        print("1. Deploy bootnode")
        print("2. Extract enode from bootnode logs")
        print("3. Update signer/member values with enode")
        print("4. Deploy signer")
        print("5. Deploy member")
        print()
        
        # Check and validate files
        print("📁 File Check & Validation:")
        all_valid = True
        for role in ['bootnode', 'signer', 'member']:
            values_file = f"values-{role}.yaml"
            if Path(values_file).exists():
                if self.validate_values_file(role, values_file):
                    print(f"   ✅ {values_file} (valid)")
                else:
                    print(f"   ❌ {values_file} (validation failed)")
                    all_valid = False
            else:
                print(f"   ❌ {values_file} (missing)")
                all_valid = False
        print()
        
        # Show current cluster context
        context = self.run_cmd("kubectl config current-context")
        print(f"🎯 Target cluster: {context.strip() if context else 'Unknown'}")
        print(f"🎯 Target namespace: {self.namespace}")
        print()
        
        if all_valid:
            print("✅ All validations passed - ready to deploy!")
        else:
            print("❌ Fix validation errors before deploying")
        
        print()
        print("⚠️  This was a dry-run. Use --install to deploy.")

    def install(self):
        """Deploy the full PoA stack"""
        print("🚀 Starting PoA Network Deployment")
        print("=" * 35)
        
        try:
            # 1. Deploy bootnode
            if not self.deploy_release('bootnode', self.releases['bootnode']):
                return False
                
            # 2. Wait for bootnode and extract enode
            if not self.wait_for_pod_ready(self.releases['bootnode']):
                print("❌ Bootnode pod failed to become ready")
                return False
                
            enode = self.extract_enode(self.releases['bootnode'])
            if not enode:
                return False
                
            # 3. Update values files
            self.update_values_with_enode(enode)
            
            # 4. Deploy signer
            if not self.deploy_release('signer', self.releases['signer']):
                return False
                
            # 5. Deploy member
            if not self.deploy_release('member', self.releases['member']):
                return False
                
            print()
            print("🎉 PoA Network Deployment Complete!")
            print("=" * 35)
            
            # Show status
            print("📊 Pod Status:")
            self.run_cmd(f"kubectl get pods -n {self.namespace}", capture_output=False)
            print()
            print("🔗 Services:")
            self.run_cmd(f"kubectl get svc -n {self.namespace}", capture_output=False)
            
            return True
            
        except KeyboardInterrupt:
            print("\n⚠️  Deployment interrupted by user")
            return False
        except Exception as e:
            print(f"❌ Deployment failed: {e}")
            return False

def main():
    parser = argparse.ArgumentParser(description='Deploy Ethereum PoA Network')
    parser.add_argument('--dry-run', action='store_true', help='Show deployment plan without executing')
    parser.add_argument('--install', action='store_true', help='Deploy the PoA network')
    parser.add_argument('--cleanup', action='store_true', help='Clean up all deployments')
    parser.add_argument('--namespace', default='ethereum-poa', help='Kubernetes namespace (default: ethereum-poa)')
    
    args = parser.parse_args()
    
    if not any([args.dry_run, args.install, args.cleanup]):
        parser.print_help()
        sys.exit(1)
        
    deployer = PoADeployer(namespace=args.namespace)
    
    if args.cleanup:
        deployer.cleanup()
    elif args.dry_run:
        deployer.dry_run()
    elif args.install:
        success = deployer.install()
        sys.exit(0 if success else 1)

if __name__ == "__main__":
    main()
