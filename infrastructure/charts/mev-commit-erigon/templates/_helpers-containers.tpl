{{/*
Common init containers for both leader and follower
*/}}
{{- define "erigon-snode.initContainers" -}}
- name: create-jwt
  image: busybox:latest
  command: ["/bin/sh", "-c"]
  args: ["cp /secrets/jwt /shared/jwt && chmod 600 /shared/jwt"]
  {{- with .Values.initContainerResources }}
  resources:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  securityContext:
    runAsUser: 0
    runAsGroup: 0
  volumeMounts:
    - name: jwt-secret
      mountPath: /secrets
    - name: erigon-volume
      mountPath: /shared
      subPath: shared
{{- if .Values.genesis.accountGeneration.enabled }}
- name: generate-accounts
  image: "ethereum/client-go:latest"
  command: ["/bin/sh", "-c"]
  args: ["/scripts/generate-account.sh"]
  {{- with .Values.initContainerResources }}
  resources:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  securityContext:
    runAsUser: 0
    runAsGroup: 0
  env:
    - name: ACCOUNT_COUNT
      value: "{{ .Values.genesis.accountGeneration.count }}"
    - name: ACCOUNT_BALANCE
      value: "{{ .Values.genesis.accountGeneration.balance | default "0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFD4A51000FDA0FFFF" }}"
  volumeMounts:
    - name: scripts-volume
      mountPath: /scripts
    - name: erigon-volume
      mountPath: /shared
      subPath: shared
    - name: genesis-template
      mountPath: /genesis-template
    - name: password-file
      mountPath: /password
    - name: erigon-volume
      mountPath: /keystore
      subPath: keystore
- name: create-genesis
  image: alpine:latest
  command: ["/bin/sh", "-c"]
  args: ["/scripts/create-genesis.sh"]
  {{- with .Values.initContainerResources }}
  resources:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  securityContext:
    runAsUser: 0
    runAsGroup: 0
  env:
    - name: ACCOUNT_BALANCE
      value: "{{ .Values.genesis.accountGeneration.balance | default "0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFD4A51000FDA0FFFF" }}"
  volumeMounts:
    - name: scripts-volume
      mountPath: /scripts
    - name: erigon-volume
      mountPath: /shared
      subPath: shared
    - name: genesis-template
      mountPath: /genesis-template
{{- else }}
- name: download-genesis
  image: alpine:latest
  command: ["/bin/sh", "-c"]
  args: ["/scripts/download-genesis.sh"]
  {{- with .Values.initContainerResources }}
  resources:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  securityContext:
    runAsUser: 0
    runAsGroup: 0
  env:
    - name: GENESIS_URL
      value: "{{ .Values.genesis.url }}"
  volumeMounts:
    - name: scripts-volume
      mountPath: /scripts
    - name: erigon-volume
      mountPath: /shared
      subPath: shared
{{- end }}
- name: init-erigon
  image: "{{ .Values.image.erigon.repository }}:{{ .Values.image.erigon.tag }}"
  imagePullPolicy: {{ .Values.image.erigon.pullPolicy }}
  command: ["/bin/sh", "-c"]
  args: ["/scripts/init-erigon.sh"]
  {{- with .Values.initContainerResources }}
  resources:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  securityContext:
    runAsUser: 0
    runAsGroup: 0
  env:
    - name: ERIGON_DATADIR
      value: "{{ .Values.erigon.datadir }}"
  volumeMounts:
    - name: scripts-volume
      mountPath: /scripts
    - name: erigon-volume
      mountPath: /shared
      subPath: shared
    - name: erigon-volume
      mountPath: {{ .Values.erigon.datadir }}
      subPath: erigon-data
{{- end }}

{{/*
Erigon container definition (shared by leader and follower)
Usage: include "erigon-snode.erigonContainer" (dict "root" . "component" "leader|follower")
*/}}
{{- define "erigon-snode.erigonContainer" -}}
- name: erigon
  image: "{{ .root.Values.image.erigon.repository }}:{{ .root.Values.image.erigon.tag }}"
  imagePullPolicy: {{ .root.Values.image.erigon.pullPolicy }}
  command: ["erigon"]
  {{- with .root.Values.resources.erigon }}
  resources:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  securityContext:
    runAsUser: 0
    runAsGroup: 0
  args:
    - --datadir={{ .root.Values.erigon.datadir }}
    - --nodekey={{ .root.Values.erigon.datadir }}/nodekey
    - --authrpc.jwtsecret=/shared/jwt
    - --authrpc.addr=0.0.0.0
    - --authrpc.port={{ .root.Values.erigon.ports.authrpc }}
    - --authrpc.vhosts=*
    {{- if .root.Values.erigon.zeroFeeTxList }}
    - --zerofee.tx.list={{ join "," .root.Values.erigon.zeroFeeTxList }}
    {{- end }}
    {{- if .root.Values.erigon.api.http.enabled }}
    - --http
    - --http.addr={{ .root.Values.erigon.api.http.addr }}
    - --http.port={{ .root.Values.erigon.ports.http }}
    - --http.vhosts={{ .root.Values.erigon.api.http.vhosts }}
    - --http.corsdomain={{ .root.Values.erigon.api.http.corsdomain }}
    - --http.api={{ .root.Values.erigon.api.http.api }}
    {{- end }}
    {{- if .root.Values.erigon.api.ws.enabled }}
    - --ws
    - --ws.port={{ .root.Values.erigon.ports.ws }}
    {{- end }}
    - --networkid={{ .root.Values.erigon.networkId }}
    {{- if .root.Values.erigon.api.metrics.enabled }}
    - --metrics
    - --metrics.addr={{ .root.Values.erigon.api.metrics.addr }}
    - --metrics.port={{ .root.Values.erigon.ports.metrics }}
    {{- end }}
    - --private.api.addr=0.0.0.0:{{ .root.Values.erigon.ports.privateApi }}
    {{- if .root.Values.erigon.nodiscover }}
    - --nodiscover
    {{- end }}
    - --prune.mode={{ .root.Values.erigon.pruneMode }}
    {{- range .root.Values.erigon.extraArgs.common }}
    - {{ . }}
    {{- end }}
    {{- if eq .component "leader" }}
    {{- range .root.Values.erigon.extraArgs.leader }}
    - {{ . }}
    {{- end }}
    {{- else if eq .component "follower" }}
    {{- range .root.Values.erigon.extraArgs.follower }}
    - {{ . }}
    {{- end }}
    {{- end }}
  ports:
    - name: http
      containerPort: {{ .root.Values.erigon.ports.http }}
    - name: ws
      containerPort: {{ .root.Values.erigon.ports.ws }}
    - name: authrpc
      containerPort: {{ .root.Values.erigon.ports.authrpc }}
    - name: metrics
      containerPort: {{ .root.Values.erigon.ports.metrics }}
    - name: p2p
      containerPort: {{ .root.Values.erigon.ports.p2p }}
      protocol: TCP
  volumeMounts:
    - name: erigon-volume
      mountPath: {{ .root.Values.erigon.datadir }}
      subPath: erigon-data
    - name: erigon-volume
      mountPath: /shared
      subPath: shared
    - name: nodekey
      mountPath: {{ .root.Values.erigon.datadir }}/nodekey
      subPath: nodekey
      readOnly: true
    {{- if .root.Values.genesis.accountGeneration.enabled }}
    - name: erigon-volume
      mountPath: /keystore
      subPath: keystore
    {{- end }}
{{- end }}

{{/*
Leader snode container definition
*/}}
{{- define "erigon-snode.snodeLeaderContainer" -}}
- name: snode
  image: "{{ .Values.image.snode.repository }}:{{ .Values.image.snode.tag }}"
  imagePullPolicy: {{ .Values.image.snode.pullPolicy }}
  command: ["snode"]
  {{- with .Values.resources.snode }}
  resources:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  securityContext:
    runAsUser: 0
    runAsGroup: 0
  args:
    - "leader"
    - --instance-id={{ .Values.consensus.leader.instanceId }}
    - --api-addr={{ .Values.consensus.leader.apiAddr }}
    - --health-addr={{ .Values.consensus.leader.healthAddr }}
    - --eth-client-url=http://localhost:{{ .Values.erigon.ports.authrpc }}
    - --jwt-secret={{ .Values.consensus.leader.jwtSecret }}
    - --non-auth-rpc-url={{ .Values.consensus.leader.nonAuthRpcUrl }}
    - --priority-fee-recipient={{ .Values.consensus.leader.priorityFeeRecipient }}
    - --evm-build-delay={{ .Values.consensus.leader.evmBuildDelay }}
    - --evm-build-delay-empty-block={{ .Values.consensus.leader.evmBuildDelayEmptyBlock }}
    - --tx-pool-polling-interval={{ .Values.consensus.leader.txPoolPollingInterval }}
    - --log-level={{ .Values.consensus.leader.logLevel }}
    - --log-fmt={{ .Values.consensus.leader.logFmt }}
    {{- if .Values.consensus.leader.logTags }}
    - --log-tags={{ .Values.consensus.leader.logTags }}
    {{- end }}
    {{- range .Values.consensus.leader.extraArgs }}
    - {{ . }}
    {{- end }}
  ports:
    - name: api
      containerPort: 9090
    - name: health
      containerPort: 8080
  {{- if .Values.consensus.liveness.enabled }}
  livenessProbe:
    httpGet:
      path: /health
      port: health
      scheme: HTTP
    initialDelaySeconds: {{ .Values.consensus.liveness.initialDelaySeconds | default 30 }}
    periodSeconds: {{ .Values.consensus.liveness.periodSeconds | default 10 }}
    timeoutSeconds: {{ .Values.consensus.liveness.timeoutSeconds | default 5 }}
    successThreshold: {{ .Values.consensus.liveness.successThreshold | default 1 }}
    failureThreshold: {{ .Values.consensus.liveness.failureThreshold | default 3 }}
  {{- end }}
  volumeMounts:
    - name: erigon-volume
      mountPath: /data/snode
      subPath: snode-data
    - name: erigon-volume
      mountPath: /shared
      subPath: shared
{{- end }}

{{/*
Follower snode container definition
*/}}
{{- define "erigon-snode.snodeFollowerContainer" -}}
- name: snode
  image: "{{ .Values.image.snode.repository }}:{{ .Values.image.snode.tag }}"
  imagePullPolicy: {{ .Values.image.snode.pullPolicy }}
  command: ["snode"]
  {{- with .Values.resources.snode }}
  resources:
    {{- toYaml . | nindent 4 }}
  {{- end }}
  securityContext:
    runAsUser: 0
    runAsGroup: 0
  args:
    - "follower"
    - --instance-id={{ .Values.consensus.follower.instanceId }}
    - --health-addr={{ .Values.consensus.follower.healthAddr }}
    - --eth-client-url=http://localhost:{{ .Values.erigon.ports.authrpc }}
    - --jwt-secret={{ .Values.consensus.follower.jwtSecret }}
    - --sync-batch-size={{ .Values.consensus.follower.syncBatchSize }}
    - --log-level={{ .Values.consensus.follower.logLevel }}
    - --log-fmt={{ .Values.consensus.follower.logFmt }}
    {{- if .Values.consensus.follower.logTags }}
    - --log-tags={{ .Values.consensus.follower.logTags }}
    {{- end }}
    {{- range .Values.consensus.follower.extraArgs }}
    - {{ . }}
    {{- end }}
  ports:
    - name: health
      containerPort: 8080
  {{- if .Values.consensus.liveness.enabled }}
  livenessProbe:
    httpGet:
      path: /health
      port: health
      scheme: HTTP
    initialDelaySeconds: {{ .Values.consensus.liveness.initialDelaySeconds | default 30 }}
    periodSeconds: {{ .Values.consensus.liveness.periodSeconds | default 10 }}
    timeoutSeconds: {{ .Values.consensus.liveness.timeoutSeconds | default 5 }}
    successThreshold: {{ .Values.consensus.liveness.successThreshold | default 1 }}
    failureThreshold: {{ .Values.consensus.liveness.failureThreshold | default 3 }}
  {{- end }}
  volumeMounts:
    - name: erigon-volume
      mountPath: /data/snode
      subPath: snode-data
    - name: erigon-volume
      mountPath: /shared
      subPath: shared
{{- end }}
