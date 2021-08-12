# Kubernetes Autoscaling Example
In this project, I try to implement Horizontal Pod Autoscaler[HPA](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) provided by [Kubernetes](https://kubernetes.io/). The Horizontal Pod Autoscaler automatically scales the number of [Pods](https://kubernetes.io/docs/concepts/workloads/pods/) in a [replication controller](https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller/), [deployment](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/), [replica set](https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/) or [stateful set](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) based on observed CPU utilization (or, with **custom metrics** support, on some other application-provided metrics).

We start by creating a [golang](https://golang.org/) web-server that publishes messages to a [rabbitmq queue](https://www.rabbitmq.com/). These messages are then consumed by a [NodeJS](https://nodejs.org/en/) worker. The worker takes 5s to consume 1 message, representing artificial time to process a request. We then try to scale the worker pods on the basis of the rabbitmq queue length. We'll discuss the entire infra in the upcoming section and role of each part in detail. 

# Architecture
<img align="center" src="./assets/HPA.png"/>
<br>

## Web-server
A simple golan server with two routes, one for pushing messages to rabbitmq queue and one for exposing metrics. We are collecting 2 metrics, `http_request_total` and `request_status` to monitor number of request on and status returned by each route respectively. 

### Routes
```
GET /generate    # push message to queue
GET /metrics     # expose metrics for prometheus server

```

## RabbitMQ
Rabbit MQ is a simple open-source message-broker that can be deployed with ease. The broker has wide support of client libraries across different programming languages. Rabbitmq is widely used in industry and has proven it's mettle. We use Advance Message Queuing Protocol([AMQP](https://www.amqp.org)) as our messaging protocol for both web-server and worker.



# How to run 
