deploy:
	kubectl apply -f k8s/

down: 
	kubectl delete -f k8s/
