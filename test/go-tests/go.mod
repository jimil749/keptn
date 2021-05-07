module github.com/keptn/keptn/test/go-tests

go 1.16

require (
	github.com/cloudevents/sdk-go/v2 v2.3.1
	github.com/google/uuid v1.2.0 // indirect
	github.com/imroc/req v0.3.0
	github.com/keptn/go-utils v0.8.4-0.20210506073402-95bb36d6f884
	github.com/keptn/keptn/shipyard-controller v0.0.0-20210503133401-8c1194432b46
	github.com/keptn/kubernetes-utils v0.8.1
	github.com/stretchr/testify v1.7.0
)

replace github.com/keptn/keptn/shipyard-controller => ../../shipyard-controller
