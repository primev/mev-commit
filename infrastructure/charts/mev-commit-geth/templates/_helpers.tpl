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

{{/*
Validate signer configuration
*/}}
{{- define "ethereum-poa.validateSigner" -}}
{{- if eq .Values.role "signer" -}}
  
  {{- /* Validate password configuration */ -}}
  {{- if and .Values.signer.password.externalSecret.enabled .Values.signer.password.value -}}
    {{- fail "ERROR: Cannot use both externalSecret and inline value for password. Please choose one method." -}}
  {{- end -}}
  
  {{- if not .Values.signer.password.externalSecret.enabled -}}
    {{- if not .Values.signer.password.value -}}
      {{- fail "ERROR: ESO is disabled but no inline password value provided. Please provide 'signer.password.value' or enable ESO with 'signer.password.externalSecret.enabled: true'" -}}
    {{- end -}}
  {{- end -}}
  
  {{- /* Validate private key configuration for signerPrivateKey method */ -}}
  {{- if eq .Values.signer.nodeBootstrapMethod "signerPrivateKey" -}}
    {{- if and .Values.signer.signerPrivateKey.externalSecret.enabled .Values.signer.signerPrivateKey.value -}}
      {{- fail "ERROR: Cannot use both externalSecret and inline value for private key. Please choose one method." -}}
    {{- end -}}
    
    {{- if not .Values.signer.signerPrivateKey.externalSecret.enabled -}}
      {{- if not .Values.signer.signerPrivateKey.value -}}
        {{- fail "ERROR: ESO is disabled but no inline private key value provided. Please provide 'signer.signerPrivateKey.value' or enable ESO with 'signer.signerPrivateKey.externalSecret.enabled: true'" -}}
      {{- end -}}
    {{- end -}}
  {{- end -}}
  
  {{- /* Validate signer address */ -}}
  {{- if not .Values.signer.address -}}
    {{- fail "ERROR: Signer address must be provided at 'signer.address'" -}}
  {{- end -}}
  
  {{- /* Validate keystore configuration for signerKeystore method */ -}}
  {{- if eq .Values.signer.nodeBootstrapMethod "signerKeystore" -}}
    {{- if not .Values.signer.signerKeystore.url -}}
      {{- if not .Values.signer.signerKeystore.content -}}
        {{- fail "ERROR: Keystore URL or content must be provided when using signerKeystore method" -}}
      {{- end -}}
    {{- end -}}
  {{- end -}}
  
{{- end -}}
{{- end }}

{{/*
Check if ESO is completely disabled for signer
*/}}
{{- define "ethereum-poa.esoDisabled" -}}
{{- if eq .Values.role "signer" -}}
  {{- if and (not .Values.signer.password.externalSecret.enabled) (not .Values.signer.signerPrivateKey.externalSecret.enabled) -}}
    {{- true -}}
  {{- end -}}
{{- end -}}
{{- end }}

{{/*
Get secret type based on configuration
*/}}
{{- define "ethereum-poa.secretType" -}}
{{- if .externalSecret.enabled -}}
external
{{- else if .value -}}
inline
{{- else -}}
existing
{{- end -}}
{{- end }}
