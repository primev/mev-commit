clean:
	helm list --no-headers | awk '{print $$1}' | xargs helm uninstall
	kubectl get pvc --no-headers | awk '{print $$1}' | xargs kubectl delete pvc

check-images:
	docker images | grep primev

clean-images:
	docker rmi $$(docker images --filter="reference=primev/primev*" -q) --force

deploy-sc:
	kubectl apply -f minikube-sc.yaml
	kubectl get storageclass

deploy-pg:
	helm install oracle-devnet ./postgres -f postgres/values.yaml
	kubectl get pods | grep oracle-devnet
