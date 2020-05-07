package k8s

// Ingress is the template for the kubernetes ingress.
const Ingress = `{{with $spec := $}}apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: {{$spec.Domain | yamlify}}-ingress
  annotations:
    kubernetes.io/ingress.class: "traefik"
    certmanager.k8s.io/issuer: "letsencrypt-prod"
    ingress.kubernetes.io/auth-tls-insecure: "false"
    ingress.kubernetes.io/ssl-redirect: "true"
spec:
  tls:
  - hosts:{{range $spec.Subdomains}}
    - {{.Name}}.{{$spec.Domain}}{{end}}
    secretName: {{$spec.Domain | yamlify}}-tls
  rules:{{range $spec.Subdomains}}
  - host: {{.Name}}.{{$spec.Domain}}
    http:
      paths:{{range .Paths}}
      - path: {{.Value}}
        backend:
          serviceName: {{.Endpoint}}
          servicePort: 8000{{end}}{{end}}
{{end}}`
