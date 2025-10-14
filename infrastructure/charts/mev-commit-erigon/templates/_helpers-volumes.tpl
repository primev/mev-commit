{{/*
Common volumes for both leader and follower
*/}}
{{- define "erigon-snode.volumes" -}}
- name: scripts-volume
  configMap:
    name: {{ include "erigon-snode.fullname" .root }}-scripts
    defaultMode: 0755
- name: jwt-secret
  secret:
    secretName: {{ include "erigon-snode.fullname" .root }}-{{ .component }}-jwt
- name: leader-nodekey
  secret:
    secretName: {{ include "erigon-snode.fullname" .root }}-leader-nodekey
{{- if .root.Values.genesis.accountGeneration.enabled }}
- name: genesis-template
  configMap:
    name: {{ include "erigon-snode.fullname" .root }}-genesis-template
- name: password-file
  secret:
    secretName: {{ include "erigon-snode.fullname" .root }}-password
{{- end }}
{{- end }}

{{/*
Common volume claim template
*/}}
{{- define "erigon-snode.volumeClaimTemplate" -}}
- metadata:
    name: erigon-volume
  spec:
    accessModes: ["{{ .Values.storage.accessMode }}"]
    resources:
      requests:
        storage: {{ .Values.storage.size }}
    {{- if .Values.storage.storageClassName }}
    storageClassName: {{ .Values.storage.storageClassName }}
    {{- end }}
{{- end }}
