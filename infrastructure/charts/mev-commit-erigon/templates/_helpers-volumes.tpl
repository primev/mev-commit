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
- name: nodekey
  secret:
    {{- if eq .component "leader" }}
    secretName: {{ include "erigon-snode.fullname" .root }}-leader-nodekey
    {{- else if eq .component "follower" }}
    secretName: {{ include "erigon-snode.fullname" .root }}-follower-nodekey
    {{- end }}
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
Usage: include "erigon-snode.volumeClaimTemplate" (dict "root" . "component" "leader|follower")
*/}}
{{- define "erigon-snode.volumeClaimTemplate" -}}
- metadata:
    name: erigon-volume
  spec:
    accessModes: ["{{ .root.Values.storage.accessMode }}"]
    resources:
      requests:
        {{- if eq .component "leader" }}
        storage: {{ .root.Values.storage.leader.size | default .root.Values.storage.size }}
        {{- else if eq .component "follower" }}
        storage: {{ .root.Values.storage.follower.size | default .root.Values.storage.size }}
        {{- else }}
        storage: {{ .root.Values.storage.size }}
        {{- end }}
    {{- $storageClass := "" }}
    {{- if eq .component "leader" }}
    {{- $storageClass = .root.Values.storage.leader.storageClassName | default .root.Values.storage.storageClassName }}
    {{- else if eq .component "follower" }}
    {{- $storageClass = .root.Values.storage.follower.storageClassName | default .root.Values.storage.storageClassName }}
    {{- else }}
    {{- $storageClass = .root.Values.storage.storageClassName }}
    {{- end }}
    {{- if $storageClass }}
    storageClassName: {{ $storageClass }}
    {{- end }}
{{- end }}
