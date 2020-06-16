package templates

// APIDefinition is the template for the go definition of the API.
const APIDefinition = `package exports

// API defines the operations supported by the {{.Name}} service.
type API interface {
	{{- range $a := .API}}
		// {{$a.Path}}
		{{range $mname, $method := $a.Methods}}
		{{- if eq $method.Type "WS" -}}
			New{{$mname | capitalize | symbolize}}Handler() (*{{$mname | capitalize | symbolize}}Handler, error)
		{{else -}}
			{{- $mname | capitalize | symbolize}}(
				{{- if $method.InputType -}}
					*{{- $method.InputType | capitalize | symbolize}},
				{{- end -}}
				{{- if $a.Path | pathHasParameters -}}
					{{- with $params := $a.Path | pathParameters -}}
						{{- range $pnameidx := $params | indicesParameters -}}
							{{- with $ptypeidx := inc $pnameidx -}}
								{{- index $params $ptypeidx | typeName -}},
						{{- end -}}
					{{- end -}}
				{{- end -}}
			{{- end -}}
			{{- if $method.QueryParams -}}
				{{- range $method.QueryParams -}}
					{{- .Type | typeName -}},
				{{- end -}}
			{{- end -}}
			{{- if $method.HeaderParams -}}
				{{- range $method.HeaderParams -}}
					{{- .Type | typeName -}},
				{{- end -}}
			{{- end -}}
			)(*{{$method.ReturnType | capitalize | symbolize}}, error)
		{{end -}}
		{{end -}}
	{{- end -}}
}`
