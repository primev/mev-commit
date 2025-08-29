{{/*
Expand the name of the chart.
*/}}
{{- define "bidder-emulator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "bidder-emulator.fullname" -}}
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
{{- define "bidder-emulator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "bidder-emulator.labels" -}}
helm.sh/chart: {{ include "bidder-emulator.chart" . }}
{{ include "bidder-emulator.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "bidder-emulator.selectorLabels" -}}
app.kubernetes.io/name: {{ include "bidder-emulator.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Generate log tags
*/}}
{{- define "bidder-emulator.logTags" -}}
{{- $tags := list }}
{{- if .Values.bidderEmulator.logTags }}
{{- $tags = append $tags .Values.bidderEmulator.logTags }}
{{- end }}
{{- $tags = append $tags (printf "service:%s" (include "bidder-emulator.name" .)) }}
{{- $tags = append $tags (printf "version:%s" .Values.version) }}
{{- $tags = append $tags (printf "namespace:%s" (.Values.namespace | default "default")) }}
{{- join "," $tags }}
{{- end }}

{{/*
Validate required values
*/}}
{{- define "bidder-emulator.validateValues" -}}
{{- if not .Values.bidderEmulator.l1RpcUrl }}
{{- fail "bidderEmulator.l1RpcUrl is required" }}
{{- end }}
{{- if not .Values.bidderEmulator.bidderRpcUrl }}
{{- fail "bidderEmulator.bidderRpcUrl is required" }}
{{- end }}
{{- if not .Values.keystores.urls }}
{{- fail "keystores.urls is required - at least one keystore URL must be provided" }}
{{- end }}
{{- if not .Values.keystores.password }}
{{- fail "keystores.password is required" }}
{{- end }}
{{- end }}
