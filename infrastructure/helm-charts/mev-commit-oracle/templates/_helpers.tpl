{{/*
Expand the name of the chart.
*/}}
{{- define "mev-oracle.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "mev-oracle.fullname" -}}
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
{{- define "mev-oracle.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "mev-oracle.labels" -}}
helm.sh/chart: {{ include "mev-oracle.chart" . }}
{{ include "mev-oracle.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "mev-oracle.selectorLabels" -}}
app.kubernetes.io/name: {{ include "mev-oracle.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the keystore secret name
*/}}
{{- define "mev-oracle.keystoreSecretName" -}}
{{- if .Values.keystore.existingSecret }}
{{- .Values.keystore.existingSecret }}
{{- else }}
{{- include "mev-oracle.fullname" . }}-keystore
{{- end }}
{{- end }}

{{/*
Create the auth token secret name
*/}}
{{- define "mev-oracle.authTokenSecretName" -}}
{{- include "mev-oracle.fullname" . }}-auth-token
{{- end }}

{{/*
Validate required configuration
*/}}
{{- define "mev-oracle.validateConfig" -}}
{{- if not .Values.oracle.authToken }}
{{- fail "oracle.authToken is required" }}
{{- end }}
{{- if not .Values.keystore.existingSecret }}
{{- if not .Values.keystore.password }}
{{- fail "keystore.password is required when not using existingSecret" }}
{{- end }}
{{- end }}
{{- if not .Values.postgresql.host }}
{{- fail "postgresql.host is required" }}
{{- end }}
{{- if not .Values.postgresql.password }}
{{- fail "postgresql.password is required" }}
{{- end }}
{{- end }}
