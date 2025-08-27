{{/*
Expand the name of the chart.
*/}}
{{- define "erigon-snode.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "erigon-snode.fullname" -}}
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
{{- define "erigon-snode.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "erigon-snode.labels" -}}
helm.sh/chart: {{ include "erigon-snode.chart" . }}
{{ include "erigon-snode.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "erigon-snode.selectorLabels" -}}
app.kubernetes.io/name: {{ include "erigon-snode.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Validate genesis configuration - ensure mutual exclusivity
*/}}
{{- define "erigon-snode.validateGenesis" -}}
{{- $accountGenEnabled := .Values.genesis.accountGeneration.enabled | default false }}
{{- $genesisUrlSet := .Values.genesis.url | default "" }}
{{- if and $accountGenEnabled (ne $genesisUrlSet "") }}
{{- fail "ERROR: Cannot use both genesis.accountGeneration.enabled=true AND genesis.url. Choose one method only." }}
{{- end }}
{{- if and (not $accountGenEnabled) (eq $genesisUrlSet "") }}
{{- fail "ERROR: Must configure either genesis.accountGeneration.enabled=true OR set genesis.url. One genesis method is required." }}
{{- end }}
{{- end }}

{{/*
Validate account generation configuration
*/}}
{{- define "erigon-snode.validateAccountGeneration" -}}
{{- if .Values.genesis.accountGeneration.enabled }}
{{- if not .Values.genesis.accountGeneration.count }}
{{- fail "ERROR: genesis.accountGeneration.count is required when account generation is enabled" }}
{{- end }}
{{- if lt (.Values.genesis.accountGeneration.count | int) 1 }}
{{- fail "ERROR: genesis.accountGeneration.count must be at least 1" }}
{{- end }}
{{- if not .Values.genesis.accountGeneration.password }}
{{- fail "ERROR: genesis.accountGeneration.password is required when account generation is enabled" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Validate JWT configuration
*/}}
{{- define "erigon-snode.validateJWT" -}}
{{- if not .Values.jwt.token }}
{{- fail "ERROR: jwt.token is required" }}
{{- end }}
{{- end }}

{{/*
Run all validations
*/}}
{{- define "erigon-snode.validate" -}}
{{- include "erigon-snode.validateGenesis" . }}
{{- include "erigon-snode.validateAccountGeneration" . }}
{{- include "erigon-snode.validateJWT" . }}
{{- end }}
