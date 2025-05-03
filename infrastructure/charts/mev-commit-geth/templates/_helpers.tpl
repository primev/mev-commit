{{/*
Expand the name of the chart.
*/}}
{{- define "ethereum-poa.name" -}}
{{- .Chart.Name -}}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "ethereum-poa.fullname" -}}
{{- printf "%s-%s" .Release.Name (include "ethereum-poa.name" .) -}}
{{- end }}

{{/*
Create a resource name with the role as suffix (bootnode, signer, member).
*/}}
{{- define "ethereum-poa.fullnameWithRole" -}}
{{- printf "%s-%s" (include "ethereum-poa.fullname" .) .Values.role -}}
{{- end }}

{{/*
Standard labels.
*/}}
{{- define "ethereum-poa.labels" -}}
app.kubernetes.io/name: {{ include "ethereum-poa.name" . }}
helm.sh/chart: {{ include "ethereum-poa.chart" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Chart version label.
*/}}
{{- define "ethereum-poa.chart" -}}
{{ printf "%s-%s" .Chart.Name .Chart.Version }}
{{- end }}
