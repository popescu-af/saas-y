package k8s

// Certificate is the template for the cluster TLS certificate.
const Certificate = `{{with $spec := $}}apiVersion: cert-manager.io/v1alpha2
kind: Certificate
metadata:
  name: {{$spec.Domain | yamlify}}
spec:
  secretName: {{$spec.Domain | yamlify}}-tls
  issuerRef:
    name: letsencrypt
    kind: ClusterIssuer
  commonName: {{$spec.Domain}}
  dnsNames:
  - {{$spec.Domain}}
  - www.{{$spec.Domain}}{{range $spec.Subdomains}}{{if eq .Name "www"}}{{else}}
  - {{.Name}}.{{$spec.Domain}}{{end}}{{end}}{{end}}`

// ClusterIssuer is the template for the cert-manager certificate issuer.
const ClusterIssuer = `apiVersion: cert-manager.io/v1alpha2
kind: ClusterIssuer
metadata:
  name: letsencrypt
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    # email: TODO
    privateKeySecretRef:
      name: letsencrypt
    solvers:
    - http01:
        ingress:
          class: traefik`
