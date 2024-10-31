apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "oracle.fullname" . }}-test-connection"
  labels:
    {{- include "oracle.labels" . | nindent 4 }}
  annotations:
    "helm.sh/hook": test
spec:
  containers:
    - name: wget
      image: busybox
      command: ['wget']
      args: ['{{ include "oracle.fullname" . }}:{{ .Values.service.port }}']
  restartPolicy: Never
