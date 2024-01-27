{{/*
Expand the name of the chart.
*/}}

{{- define "zbi.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "zbi.fullname" -}}
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
{{- define "zbi.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "zbi.labels" -}}
helm.sh/chart: {{ include "zbi.chart" . }}
{{ include "zbi.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "zbi.selectorLabels" -}}
app.kubernetes.io/name: {{ include "zbi.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Controller name
*/}}
{{- define "zbi-controller.fullname" -}}
{{- printf "%s-controller" (include "zbi.name" .) }}
{{- end }}

{{/*
Blockchain domain
*/}}
{{- define "zbi-blockchain.domain" -}}
{{- printf "%sblockchain.%s" ( .Values.prefix ) ( .Values.domain ) }}
{{- end }}

{{/*
Control plane domain
*/}}
{{- define "zbi-control-plane.domain" -}}
{{- printf "%sapp.%s" ( .Values.prefix ) ( .Values.domain ) }}
{{- end }}

{{/*
Database name
*/}}
{{- define "zbi-db.fullname" -}}
{{- printf "%s-db" (include "zbi.name" .) }}
{{- end }}


{{/*
Control Plane name
*/}}
{{- define "zbi-control-plane.fullname" -}}
{{- printf "%s-control-plane" (include "zbi.name" .) }}
{{- end }}

{{/*
Database Selector labels
*/}}
{{- define "zbi-db.selectorLabels" -}}
app.kubernetes.io/name: {{ include "zbi-db.fullname" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}


{{/*
Controller Selector labels
*/}}
{{- define "zbi-controller.selectorLabels" -}}
app.kubernetes.io/name: {{ include "zbi-controller.fullname" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Control-Plane Selector labels
*/}}
{{- define "zbi-control-plane.selectorLabels" -}}
app.kubernetes.io/name: {{ include "zbi-control-plane.fullname" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "zbi.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "zbi-control-plane.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
