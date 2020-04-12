# saas-y
saas-y generates framework code and configuration for a SaaS platform.

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
            "port": 80,
            "api": {
                "type": "http",
                "paths": [
                    {
                        "path": "/getfoo",
                        "parameters": ["param0", "param1"],
                        "return": "return_struct_name"
                    },
                    {
                        "path": "/postfoo",
                        "input": "input_struct_name",
                        "return": "return_struct_name"
                    }
                ]
            },
            "env": [
                {
                    "name": "ENV_VAR_NAME",
                    "value": "env_var_value"
                }
            ],
            "dependencies" : ["baz-service"]
        }
    ],
    "structs": [
        {
            "name": "input_struct_name",
            "fields": [
                {
                    "name": "a_field_name",
                    "type": "int"
                }
            ]
        }
    ],
    "external_services": [
        {
            "name": "external-service",
            "port": 80,
            "image_url": "localhost:5000/external-service:latest",
            "env": [
                {
                    "name": "ENV_VAR_NAME",
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

## Framework generation
TODO: description


## Installation
TODO

## Usage
```bash
# TODO
```

## License
Proprietary
