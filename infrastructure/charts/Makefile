clean:
	helm list --no-headers | awk '{print $$1}' | xargs helm uninstall
	kubectl get pvc --no-headers | awk '{print $$1}' | xargs kubectl delete pvc
	kubectl delete storageclass premium-rwo

check-images:
	docker images | grep primev

clean-images:
	docker rmi $$(docker images --filter="reference=primev/primev*" -q) --force

deploy-devnet:
	kubectl apply -f minikube-sc.yaml
	kubectl get storageclass
	sleep 10
	helm install oracle-devnet ./postgres -f postgres/values.yaml
	kubectl get pods | grep oracle-devnet
	@BRANCH=$$(git rev-parse --abbrev-ref HEAD) && \
	echo "Deploying with branch: $$BRANCH" && \
	bash deploy.sh --branch $$BRANCH --password CXBMf4xEDO0I --local true
