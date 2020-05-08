package k8s

// Registry is the template for the docker registry.
const Registry = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: docker-registry
  labels:
    app: docker-registry
spec:
  replicas: 1
  selector:
    matchLabels:
      app: docker-registry
  template:
    metadata:
      labels:
        app: docker-registry
    spec:
      containers:
        - name: docker-registry
          image: registry:latest
          ports:
            - containerPort: 5000
---
apiVersion: v1
kind: Service
metadata:
  name: docker-registry
spec:
  type: NodePort
  selector:
    app: docker-registry
  ports:
    - name: http
      protocol: TCP
      port: 5000
      nodePort: 32000
      targetPort: 5000
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: docker-registry-ui
  labels:
    app: docker-registry-ui
spec:
  replicas: 1
  selector:
    matchLabels:
      app: docker-registry-ui
  template:
    metadata:
      labels:
        app: docker-registry-ui
    spec:
      containers:
        - name: docker-registry-ui
          image: jc21/registry-ui
          ports:
            - containerPort: 80
          env:
            - name: REGISTRY_HOST
              value: docker-registry:5000
            - name: REGISTRY_SSL
              value: "false"
---
apiVersion: v1
kind: Service
metadata:
  name: docker-registry-ui
spec:
  selector:
    app: docker-registry-ui
  ports:
    - name: http
      protocol: TCP
      port: 6000
      targetPort: 80`
