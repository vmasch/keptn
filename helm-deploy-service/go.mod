module github.com/keptn/keptn/helm-deploy-service

go 1.13

require (
	github.com/cloudevents/sdk-go v0.10.0
	github.com/google/uuid v1.1.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/keptn/go-utils v0.6.3-0.20200730110526-91896e4be8e0
	github.com/keptn/kubernetes-utils v0.2.0
	helm.sh/helm/v3 v3.1.2
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/cli-runtime v0.17.2
	k8s.io/client-go v0.17.2
	k8s.io/kubectl v0.17.2
)

// Transitive requirement from Helm: See https://github.com/helm/helm/blob/v3.1.2/go.mod
replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
)
