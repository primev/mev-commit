{{/*
Expand the name of the chart.
*/}}
{{- define "relay-emulator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "relay-emulator.fullname" -}}
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
{{- define "relay-emulator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "relay-emulator.labels" -}}
helm.sh/chart: {{ include "relay-emulator.chart" . }}
{{ include "relay-emulator.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: relay-emulator
{{- end }}

{{/*
Selector labels
*/}}
{{- define "relay-emulator.selectorLabels" -}}
app.kubernetes.io/name: {{ include "relay-emulator.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create log tags with service info
*/}}
{{- define "relay-emulator.logTags" -}}
{{- $baseTags := printf "service.name:%s,service.version:%s" .Values.job.name .Chart.AppVersion -}}
{{- if index .Values.job.env "log-tags" -}}
{{- printf "%s,%s" $baseTags (index .Values.job.env "log-tags") -}}
{{- else -}}
{{- $baseTags -}}
{{- end -}}
{{- end }}

{{/*
Validate required values
*/}}
{{- define "relay-emulator.validateValues" -}}
{{- if not .Values.job.l1RpcUrl }}
{{- fail "job.l1RpcUrl is required" }}
{{- end }}
{{- if not .Values.job.httpPort }}
{{- fail "job.httpPort is required" }}
{{- end }}
{{- end }}
