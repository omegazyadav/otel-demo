NAMESPACE_DEMO_APP      := default
NAMESPACE_POSTGRES      := default
NAMESPACE_INGRESS       := ingress-nginx

CHART_DEMO_APP   := ./helm/demo-app/
CHART_POSTGRES   := ./helm/postgres/
CHART_INGRESS    := dev/addons/ingress

install: install-postgres install-demo-app

install-demo-app:
	helm upgrade --install demo-app $(CHART_DEMO_APP) \
		--namespace $(NAMESPACE_DEMO_APP) \
		--create-namespace \
		--dependency-update

install-postgres:
	helm upgrade --install postgres $(CHART_POSTGRES) \
		--namespace $(NAMESPACE_POSTGRES) \
		--create-namespace \
		--dependency-update

install-ingress:
	helm upgrade --install ingress $(CHART_INGRESS) \
		--namespace $(NAMESPACE_INGRESS) \
		--create-namespace \
		--dependency-update
