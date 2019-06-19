## Serf-Publisher

This k8s controller expose your application throught Serf.


# Dependencies

You will need install GO, in order to build the project

You can follow [this guide](https://golang.org/doc/install)

then do: `export PATH=$PATH:/usr/local/go/bin`


# Generation of binary

```sh
make (arm|linux)
```


Before released it, you must ensure the version of your binary which will pushed to your registry, please check the [Dockerfile](Dockerfile) and set whatever you need.

# Docker release
```sh
make release
```

# Run 

1- Copy your binary to master node.

2- Just run the next command in the master node:
```sh
./serf-publisher --namespace your_namespace --kubeconfig your_kube_config_file
```

Or copy your binary to your master Kubernetes node and create an [unit](serf-publisher.service).

3- Create your service as NodePort

**NOTE:**
This project is purely academic.
