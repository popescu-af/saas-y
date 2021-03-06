package templates

// Readme is the template for the service's README.md file.
const Readme = `# {{.Name}}
This service was generated with saas-y. Please modify a copy of this file, by adding to it:
- a proper description
- diagrams
- etc.

## Quick start
Run
` + "```" + `bash
# build docker image
make build

# run docker image locally
make run

# publish docker image to the container repository
# (default localhost:5000, see Makefile)
make publish

# deploy a new version of the service to k8s cluster
make deploy
` + "```" + `

## Documentation
TODO: add documentation (e.g. https://c4model.com/)

## Contributors
- John Doe
- Max Mustermann
- Ion Popescu
- Hong Gildong
- Fulano/Fulanita

## License
Please check LICENSE.md

## Changelog
Please check CHANGELOG.md
`
