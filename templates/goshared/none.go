package goshared

const noneTpl = `// {{ name .Field }} 字段没有验证规则
	{{- if .Index }}[{{ .Index }}]{{ end }}
	{{- if .OnKey }} (key){{ end }}`
