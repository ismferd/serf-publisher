## Serf-Publisher

This k8s controller expose your application throught Serf.

# Generation of binary

```sh
make (arm|linux)
```


Before released it, you must ensure the version of your binary which will pushed to your registry, please check the [Dockerfile](Dockerfile) and set whatever you need.

# Docker release
```sh
make release
```



In order to run this operator you could deploy it using the [manifest](serf-publisher-deployment.yaml)

Or copy your binary to your master Kubernetes node and create an [unit](serf-publisher.service).

**NOTE:**
This project is purely academic.
