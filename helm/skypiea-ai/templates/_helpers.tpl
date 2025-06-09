{{/*
Expand the name of the chart.
*/}}
{{- define "skypiea-ai.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "skypiea-ai.fullname" -}}
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
{{- define "skypiea-ai.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "skypiea-ai.labels" -}}
helm.sh/chart: {{ include "skypiea-ai.chart" . }}
{{ include "skypiea-ai.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "skypiea-ai.selectorLabels" -}}
app.kubernetes.io/name: {{ include "skypiea-ai.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "skypiea-ai.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "skypiea-ai.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Database secrets
*/}}
{{- define "skypiea-ai.secret.db.env" -}}
- name: DATABASE_{{ .Values.database.dialect | upper }}_USERNAME
  valueFrom:
    secretKeyRef:
      name: {{ .Values.db_secret }}
      key: {{ .Values.database.dialect | lower }}_username
- name: DATABASE_{{ .Values.database.dialect | upper }}_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ .Values.db_secret }}
      key: {{ .Values.database.dialect | lower }}_password
{{- end -}}