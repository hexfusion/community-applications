version: '3.8'
services:
{{- range .compose.services }} # using compose namespace
  {{ .name }}:
    image: "{{ .image.repository }}:{{ .image.tag }}"
    restart: unless-stopped
    environment:
    {{- range $key, $value := .envVars }}
      - {{ $key }}={{ $value }}
    {{- end }}
    volumes:
    {{- range .volumes }}
      - "{{ .source }}:{{ .path }}"
    {{- end }}
    ports:
    {{- range .ports }}
      - "{{ .hostPort }}:{{ .containerPort }}"
    {{- end }}
{{- end }}

networks:
  {{ .network.name }}: # using global
    driver: {{ .network.driver }}