{{- define "sun-panel.config" -}}
  conf.ini: |
    [base]
    http_port=3002
    database_drive=mysql
    cache_drive=memory
    queue_drive=memory
    source_path=./uploads
    source_temp_path=./runtime/temp

    [mysql]
    host={{ .Values.mysql.host}}
    port={{ .Values.mysql.port}}
    username={{ .Values.mysql.username}}
    password={{ .Values.mysql.password}}
    db_name={{ .Values.mysql.dbname}}
    wait_timeout=100
{{- end }}