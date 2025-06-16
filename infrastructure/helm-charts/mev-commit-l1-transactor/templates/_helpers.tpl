{{/*
Expand the name of the chart.
*/}}
{{- define "l1-transactor.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "l1-transactor.fullname" -}}
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
{{- define "l1-transactor.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "l1-transactor.labels" -}}
helm.sh/chart: {{ include "l1-transactor.chart" . }}
{{ include "l1-transactor.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "l1-transactor.selectorLabels" -}}
app.kubernetes.io/name: {{ include "l1-transactor.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Generate log tags
*/}}
{{- define "l1-transactor.logTags" -}}
{{- $tags := list -}}
{{- $tags = append $tags (printf "service.name:%s" .Values.job.name) -}}
{{- $tags = append $tags (printf "service.version:%s" .Values.version) -}}
{{- if .Values.job.env.logTags -}}
{{- $tags = append $tags .Values.job.env.logTags -}}
{{- end -}}
{{- join "," $tags -}}
{{- end }}

{{/*
Get secret name
*/}}
{{- define "l1-transactor.secretName" -}}
{{- if .Values.existingSecret -}}
{{- .Values.existingSecret -}}
{{- else -}}
{{- include "l1-transactor.fullname" . }}-secrets
{{- end -}}
{{- end }}

{{/*
Validate required values
*/}}
{{- define "l1-transactor.validateValues" -}}
{{- if not .Values.l1RpcUrl -}}
{{- fail "l1RpcUrl is required" -}}
{{- end -}}
{{- if not .Values.image.repository -}}
{{- fail "image.repository is required" -}}
{{- end -}}
{{- if not .Values.image.tag -}}
{{- fail "image.tag is required" -}}
{{- end -}}
{{- end }}
