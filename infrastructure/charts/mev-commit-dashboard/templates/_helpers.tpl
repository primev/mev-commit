{{/*
Expand the name of the chart.
*/}}
{{- define "mev-commit-dashboard.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "mev-commit-dashboard.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "mev-commit-dashboard.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "mev-commit-dashboard.labels" -}}
helm.sh/chart: {{ include "mev-commit-dashboard.chart" . }}
{{ include "mev-commit-dashboard.selectorLabels" . }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- if .Values.commonLabels }}
{{ toYaml .Values.commonLabels }}
{{- end }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "mev-commit-dashboard.selectorLabels" -}}
app.kubernetes.io/name: {{ include "mev-commit-dashboard.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app: {{ include "mev-commit-dashboard.name" . }}
{{- end }}
