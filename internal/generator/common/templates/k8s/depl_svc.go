package k8s

// DeplSvc is the template for the kubernetes deployment & service.
const DeplSvc = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Name}}
  labels:
    app: {{.Name}}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{.Name}}
  template:
    metadata:
      labels:
        app: {{.Name}}
    spec:
      containers:
        - name: {{.Name}}
          # Should be set to the right container registry if needed.
          image: localhost:32000/{{.Name}}:latest
          imagePullPolicy: Always
          ports:
            - containerPort: {{.Port}}
              name: http-port
          env:
            # Some environment variables might need to be read
            # from secrets or other entities.
            - name: APP_PORT
              value: "{{.Port}}"{{range .Environment}}
            - name: APP_{{.Name | toUpper}}
              value: "{{.Value}}"{{end}}
            {{range $d := .Dependencies -}}
            - name: APP_{{$d | replaceHyphens | toUpper}}_ADDR
              value: "{{$d}}:8000"
            {{- end}}
---
apiVersion: v1
kind: Service
metadata:
  name: {{.Name}}
spec:
  selector:
    app: {{.Name}}
  ports:
    - name: http
      protocol: TCP
      port: 8000
      targetPort: {{.Port}}`
