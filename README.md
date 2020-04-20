# saas-y
saas-y generates framework code and configuration for a SaaS platform.

## Installation
TODO

## Usage
```bash
# TODO
```

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
                        "get": {
                            "header_params": [
                                {
                                    "name": "header_param_name",
                                    "type": "type",
                                    "value": "value"
                                }
                            ],
                            "query_params": [
                                {
                                    "name": "query_param_name",
                                    "type": "type",
                                    "value": "value"
                                }
                            ],
                            "return_type": "return_struct_name"
                        },
                        "post": {
                            "input_type": "input_struct_name",
                            "return_type": "return_struct_name"
                        },
                        "options": {
                            "header_params": [
                                {
                                    "name": "header_param_name",
                                    "value": "value"
                                }
                            ]
                        }
                    }
                },
                {
                    "path": "/bar/{param(type)}",
                    "methods": {
                        "get": {
                            "return_type": "return_struct_name"
                        },
                        "delete": {
                            "return_type": "return_struct_name"
                        },
                        "patch": {
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
                        }
                    ]
                }
            ],
            "env": [
                {
                    "name": "ENV_VAR_NAME",
                    "type": "type",
                    "value": "env_var_value"
                }
            ],
            "dependencies" : ["baz-service"]
        }
    ],
    "external_services": [
        {
            "name": "external-service",
            "port": "80",
            "image_url": "localhost:5000/external-service:latest",
            "env": [
                {
                    "name": "ENV_VAR_NAME",
                    "type": "type",
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

### services
TODO: services' JSON -> AST services -> generated directories and files

## License
Proprietary
