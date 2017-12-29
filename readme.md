# a naive autoscaler
auto-scale kubernetes deployments based on arbitrary datadog metric-series.

- Configure a k8s job to run at some interval, per k8s deployment, ie: 5 minutes.
- Set the job env vars for configuration

> note: it currently depends on having a kubeconf locally available.
> I will update it eventually, to use in-cluster authentication

# install dependencies
run deps.sh
`./deps.sh`
go get k8s.io/client-go/...
go get -u k8s.io/apimachinery/...


# build
```
export GOPATH=$(pwd)
go build
```
