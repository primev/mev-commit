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
{{- if not .Values.consensus.leader.jwtSecret }}
{{- fail "ERROR: consensus.leader.jwtSecret is required" }}
{{- end }}
{{- if .Values.consensus.follower.enabled }}
{{- if not .Values.consensus.follower.jwtSecret }}
{{- fail "ERROR: consensus.follower.jwtSecret is required when follower is enabled" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Validate follower configuration
*/}}
{{- define "erigon-snode.validateFollower" -}}
{{- if .Values.consensus.follower.enabled }}
{{- if not .Values.consensus.follower.postgresDb.dsn }}
{{- fail "ERROR: consensus.follower.postgresDb.dsn is required when follower is enabled" }}
{{- end }}
{{- if not .Values.consensus.follower.instanceId }}
{{- fail "ERROR: consensus.follower.instanceId is required when follower is enabled" }}
{{- end }}
{{- $replicaCount := .Values.consensus.follower.replicaCount | default 1 }}
{{- if lt ($replicaCount | int) 1 }}
{{- fail "ERROR: consensus.follower.replicaCount must be at least 1" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Validate leader configuration
*/}}
{{- define "erigon-snode.validateLeader" -}}
{{- if not .Values.consensus.leader.instanceId }}
{{- fail "ERROR: consensus.leader.instanceId is required for leader node" }}
{{- end }}
{{- if not .Values.consensus.leader.priorityFeeRecipient }}
{{- fail "ERROR: consensus.leader.priorityFeeRecipient is required for leader node" }}
{{- end }}
{{- end }}

{{/*
Validate erigon configuration
*/}}
{{- define "erigon-snode.validateErigon" -}}
{{- if not .Values.erigon.networkId }}
{{- fail "ERROR: erigon.networkId is required" }}
{{- end }}
{{- if not .Values.erigon.datadir }}
{{- fail "ERROR: erigon.datadir is required" }}
{{- end }}
{{- if not .Values.erigon.ports.p2p }}
{{- fail "ERROR: erigon.ports.p2p is required" }}
{{- end }}
{{- if not .Values.erigon.ports.authrpc }}
{{- fail "ERROR: erigon.ports.authrpc is required" }}
{{- end }}
{{- end }}

{{/*
Validate storage configuration
*/}}
{{- define "erigon-snode.validateStorage" -}}
{{- if not .Values.storage.size }}
{{- fail "ERROR: storage.size is required" }}
{{- end }}
{{- if not .Values.storage.accessMode }}
{{- fail "ERROR: storage.accessMode is required" }}
{{- end }}
{{- end }}

{{/*
Validate image configuration
*/}}
{{- define "erigon-snode.validateImages" -}}
{{- if not .Values.image.erigon.repository }}
{{- fail "ERROR: image.erigon.repository is required" }}
{{- end }}
{{- if not .Values.image.erigon.tag }}
{{- fail "ERROR: image.erigon.tag is required" }}
{{- end }}
{{- if not .Values.image.snode.repository }}
{{- fail "ERROR: image.snode.repository is required" }}
{{- end }}
{{- if not .Values.image.snode.tag }}
{{- fail "ERROR: image.snode.tag is required" }}
{{- end }}
{{- end }}

{{/*
Run all validations
*/}}
{{- define "erigon-snode.validate" -}}
{{- include "erigon-snode.validateGenesis" . }}
{{- include "erigon-snode.validateAccountGeneration" . }}
{{- include "erigon-snode.validateJWT" . }}
{{- include "erigon-snode.validateFollower" . }}
{{- include "erigon-snode.validateLeader" . }}
{{- include "erigon-snode.validateErigon" . }}
{{- include "erigon-snode.validateStorage" . }}
{{- include "erigon-snode.validateImages" . }}
{{- end }}
