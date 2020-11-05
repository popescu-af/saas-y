# saas-y
saas-y is a sassy SaaS framework code and configuration generator. It thinks developers should not spend time writing boilerplate for the cloud infrastructure and glue code.

Of course, saas-y is NOT very mature as of now, so one might find many problems / things to improve. But most definitely it strives to be a useful tool, at least for simple cloud apps, made up of one or more services.

## Prerequisites

* public repo for your tutorial project
* k8s cluster (or local minikube, microk8s or k3s)
* (optional) docker registry (one can use the docker-registry.yaml file to deploy an instance of the registry)

## Installation and Usage

```bash
git clone https://github.com/popescu-af/saas-y.git
cd saas-y
go install cmd/saas-y.go

# Basic example: edit hello-world.json to point to a repository you own
saas-y example/hello-world.json tutorial

# After the `tutorial` directory is generated, its contents need to be
# copied to a local clone of the repository mentioned in the hello-world.json
# and pushed.
cp tutorial/* /path/to/local/clone
cd /path/to/local/clone
git add .
git commit -m "Initial commit."
git push
```

After the last command, a directory named `tutorial` should have been created.
A few steps are required for the tutorial program to be implemented properly.

* first of all, the code 

* edit `services/time-svc/internal/logic/impl.go`
  * replace `"errors"` with `"time"` in the list of imports
  * replace _`GetTime`_ function with the following implementation
```go
func (i *Implementation) GetTime() (*exports.Time, error) {
	log.Info("called get_time")
	return &exports.Time{
		Value: time.Now().String(),
	}, nil
}
```

* edit `services/tutorial-svc/internal/logic/impl.go`
  * replace `"errors"` with `"fmt"` in the list of imports
  * replace _`Greet`_ function with the following implementation
```go
func (i *Implementation) Greet(name string) (*exports.Greeting, error) {
	log.Info("called greet")

	t, err := timesvc.GetTime()
	if err != nil {
		return nil, err
	}

	return &exports.Greeting{
		Message: fmt.Sprintf("Hello, %s! Current time is %s", name, t.Value),
	}, nil
}
```

* (optional) deploy the docker-registry service to your cluster
```bash
kubectl apply -f tutorial/deploy/docker-registry.yaml
# wait a bit for it to start (use k9s to monitor)
```

* port-forward your docker registy to `localhost:5000`
  * `docker.for.mac.localhost` (KEEP IN MIND FOR MAC)

cd services/time-svc && make deploy
cd ../tutorial-svc && make deploy

* port-forward the tutorial-svc
* curl it
  * that's it

## Input
The input consists of a JSON file with the following format
```json
{
    "domain": "example.com",
    "subdomains": [
        {
            "name": "api",
            "paths": [
                {
                    "value": "/f/o/o",
                    "endpoint": "foo-service"
                },
                {
                    "value": "/b/a/r",
                    "endpoint": "bar-service"
                }
            ]
        }
    ],
    "services": [
        {
            "name": "foo-service",
            "port": "80",
            "api": [
                {
                    "path": "/foo",
                    "methods": {
                        "method_name_0": {
                            "type": "GET",
                            "header_params": [
                                {
                                    "name": "header_param_name",
                                    "type": "int"
                                }
                            ],
                            "query_params": [
                                {
                                    "name": "query_param_name",
                                    "type": "string"
                                }
                            ],
                            "return_type": "return_struct_name"
                        },
                        "method_name_1": {
                            "type": "POST",
                            "input_type": "input_struct_name",
                            "return_type": "return_struct_name"
                        },
                    }
                },
                {
                    "path": "/bar/{rank:uint}/{price:float}",
                    "methods": {
                        "method_name_3": {
                            "type": "GET",
                            "return_type": "return_struct_name"
                        },
                        "method_name_4": {
                            "type": "DELETE",
                            "return_type": "return_struct_name"
                        },
                        "method_name_5": {
                            "type": "PATCH",
                            "input_type": "input_struct_name",
                            "return_type": "return_struct_name"
                        }
                    }
                }
            ],
            "structs": [
                {
                    "name": "input_struct_name",
                    "fields": [
                        {
                            "name": "a_field_name",
                            "type": "int"
                        },
                        {
                            "name": "another_field_name",
                            "type": "string"
                        }
                    ]
                }
            ],
            "env": [
                {
                    "name": "ENV_VAR_NAME",
                    "type": "int64",
                    "value": "42"
                }
            ],
            "dependencies" : ["baz-service"]
        }
    ],
    "external_services": [
        {
            "name": "external-service",
            "repository_url": "example.com/external-service",
            "port": "80",
            "image_url": "localhost:5000/external-service:latest",
            "env": [
                {
                    "name": "ENV_VAR_NAME",
                    "type": "string",
                    "value": "env_var_value"
                }
            ],
            "dependencies" : ["yet-another-external-service"]
        }
    ]
}
```

## Concepts
| Name | Definition | Example |
| ---- | ---------- | ------- |
| domain | where the SaaS is supposed to be deployed | example.com |
| subdomains | list of subdomains that will be accessible | api, www, ... |
| path | domain path to be routed to a service | /rest/1.0.0 |
| endpoint | service listening for requests, addressed by name only | cool-service |
| api | a description of the types of requests the service supports | see above JSON |
| env | list of environmental variables a service needs | see above JSON |
| dependencies | list of services the service depends on | see above JSON |
| structs | list of structures used in by the APIs of all services | see above JSON |
| external_services | list of services that are build elsewhere, to be directly used by means of pre-built docker images | see above JSON |

## More Detailed Description
_coming soon_

## License
MIT
