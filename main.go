// this should pull config from env vars
// create job template for each service, which has settings pulled from manifest
// k8s cron job should be fine, per deployment
package main

import (
	"fmt"
	"os"
	"strconv"
)

type ScaleConfig struct {
	Namespace     string
	Deployment    string
	MinReplicas   int64
	MaxReplicas   int64
	Query         string
	QueryPeriod   int64
	LowWatermark  float64
	HighWatermark float64
}

func main() {
	config := loadConfigFromEnv()
	autoscale(config)
}

func loadConfigFromEnv() ScaleConfig {

	namespace := requiredEnv("NAMESPACE")
	deployment := requiredEnv("DEPLOYMENT")
	minReplicas := asInt(requiredEnv("MIN_REPLICAS"))
	maxReplicas := asInt(requiredEnv("MAX_REPLICAS"))
	ddQuery := requiredEnv("DD_QUERY")
	ddQueryPeriod := asInt(requiredEnv("DD_QUERY_PERIOD"))
	minWatermark := asFloat(requiredEnv("MIN_WATERMARK"))
	maxWatermark := asFloat(requiredEnv("MAX_WATERMARK"))

	config := ScaleConfig{namespace, deployment, minReplicas, maxReplicas, ddQuery, ddQueryPeriod, minWatermark, maxWatermark}
	return config
}

func asInt(value string) int64 {
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err.Error())
	}
	return n
}
func asFloat(value string) float64 {
	n, err := strconv.ParseFloat(value, 64)
	if err != nil {
		panic(err.Error())
	}
	return n
}

func requiredEnv(env string) string {
	value := os.Getenv(env)
	if len(value) == 0 {
		panic(fmt.Sprintf("Required env var missing: %s\n", env))
	}
	return value
}

func autoscale(config ScaleConfig) {
	queryResult := GetDataDogMetrics(config.Query, config.QueryPeriod)
	deployment := GetDeployment(config)
	qr := queryResult / float64(deployment.Replicas)

	fmt.Printf("replicas: %d, watermark: %f\n", deployment.Replicas, qr)

	if qr < config.LowWatermark { // scale down
		if deployment.Replicas > config.MinReplicas {
			fmt.Printf("Scaling down: %s\n", config)
			SetDeploymentScale(config, deployment.Replicas-1)
		}
	} else if qr > config.HighWatermark { // scale up
		if deployment.Replicas < config.MaxReplicas {
			fmt.Printf("Scaling up: %s\n", config)
			SetDeploymentScale(config, deployment.Replicas+1)
		}
	}
}
