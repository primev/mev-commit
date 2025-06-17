{{/*
Expand the name of the chart.
*/}}
{{- define "snode-geth.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "snode-geth.fullname" -}}
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
{{- define "snode-geth.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common selector labels. Keep these minimal and stable.
*/}}
{{- define "snode-geth.selectorLabels" -}}
app.kubernetes.io/name: {{ include "snode-geth.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app: {{ include "snode-geth.name" . }}
{{- end }}

{{/*
Common metadata labels (can be customized).
*/}}
{{- define "snode-geth.labels" -}}
{{ include "snode-geth.selectorLabels" . }}
helm.sh/chart: {{ include "snode-geth.chart" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
{{- with .Values.commonLabels }}
{{ toYaml . | nindent 0 }}
{{- end }}
{{- end }}

{{/*
Common annotations.
*/}}
{{- define "snode-geth.annotations" -}}
{{- with .Values.commonAnnotations }}
{{ toYaml . | nindent 0 }}
{{- end }}
{{- end }}

{{/*
Service selector override.
*/}}
{{- define "snode-geth.serviceSelector" -}}
{{- if .Values.service.selector }}
{{ toYaml .Values.service.selector | nindent 0 }}
{{- else }}
app: {{ include "snode-geth.name" . }}
{{- end }}
{{- end }}
