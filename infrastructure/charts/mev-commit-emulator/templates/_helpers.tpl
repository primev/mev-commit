{{/*
Expand the name of the chart.
*/}}
{{- define "emulator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "emulator.fullname" -}}
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
{{- define "emulator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "emulator.labels" -}}
helm.sh/chart: {{ include "emulator.chart" . }}
{{ include "emulator.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: {{ .Values.job.type }}-emulator
{{- end }}

{{/*
Selector labels
*/}}
{{- define "emulator.selectorLabels" -}}
app.kubernetes.io/name: {{ include "emulator.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Validate emulator type and required values
*/}}
{{- define "emulator.validateValues" -}}
{{/* All validation removed */}}
{{- end }}

{{/*
Get the appropriate command based on emulator type
*/}}
{{- define "emulator.command" -}}
{{- if eq .Values.job.type "bidder" -}}
/usr/local/bin/bidder-emulator
{{- else if eq .Values.job.type "provider" -}}
/usr/local/bin/provider-emulator
{{- end -}}
{{- end }}
