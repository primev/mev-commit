{{/*
Expand the name of the chart.
*/}}
{{- define "mev-commit-p2p.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "mev-commit-p2p.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "mev-commit-p2p.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "mev-commit-p2p.labels" -}}
helm.sh/chart: {{ include "mev-commit-p2p.chart" . }}
{{ include "mev-commit-p2p.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "mev-commit-p2p.selectorLabels" -}}
app.kubernetes.io/name: {{ include "mev-commit-p2p.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Define TLS secret name based on settings
*/}}
{{- define "mev-commit-p2p.tlsSecretName" -}}
{{- if .Values.global.tls.selfSigned -}}
{{ include "mev-commit-p2p.fullname" . }}-tls
{{- else if .Values.global.tls.existingSecret -}}
{{ .Values.global.tls.existingSecret }}
{{- else -}}
{{ include "mev-commit-p2p.fullname" . }}-tls
{{- end -}}
{{- end }}

{{/*
Generate self-signed certificate if enabled
*/}}
{{- define "mev-commit-p2p.generateCertificate" -}}
{{- $altNames := list (printf "%s-%s" (include "mev-commit-p2p.fullname" .) .Values.node.type) (printf "%s-%s.%s" (include "mev-commit-p2p.fullname" .) .Values.node.type .Release.Namespace) (printf "%s-%s.%s.svc" (include "mev-commit-p2p.fullname" .) .Values.node.type .Release.Namespace) -}}
{{- $ca := genCA "mev-commit-ca" 365 -}}
{{- $cert := genSignedCert (include "mev-commit-p2p.fullname" .) nil $altNames 365 $ca -}}
tls.crt: {{ $cert.Cert | b64enc }}
tls.key: {{ $cert.Key | b64enc }}
{{- end -}}
