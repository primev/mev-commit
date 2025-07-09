{{/*
Expand the name of the chart.
*/}}
{{- define "mev-commit-validator-monitor.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "mev-commit-validator-monitor.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "mev-commit-validator-monitor.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "mev-commit-validator-monitor.labels" -}}
helm.sh/chart: {{ include "mev-commit-validator-monitor.chart" . }}
{{ include "mev-commit-validator-monitor.selectorLabels" . }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- if .Values.commonLabels }}
{{ toYaml .Values.commonLabels }}
{{- end }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "mev-commit-validator-monitor.selectorLabels" -}}
app.kubernetes.io/name: {{ include "mev-commit-validator-monitor.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app: {{ include "mev-commit-validator-monitor.name" . }}
{{- end }}

{{/*
Create the name of the database secret to use
*/}}
{{- define "mev-commit-validator-monitor.dbSecretName" -}}
{{- if .Values.database.credentials.existingSecret }}
{{- .Values.database.credentials.existingSecret }}
{{- else }}
{{- printf "%s-db-credentials" (include "mev-commit-validator-monitor.fullname" .) }}
{{- end }}
{{- end }}

{{/*
Create the name of the slack secret to use
*/}}
{{- define "mev-commit-validator-monitor.slackSecretName" -}}
{{- if .Values.notifications.slack.existingSecret }}
{{- .Values.notifications.slack.existingSecret }}
{{- else }}
{{- printf "%s-slack-webhook" (include "mev-commit-validator-monitor.fullname" .) }}
{{- end }}
{{- end }}
