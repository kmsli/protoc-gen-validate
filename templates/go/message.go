package golang

// Embedded message validation.
const messageTpl = `
	{{ $f := .Field }}{{ $r := .Rules }}
	{{ template "required" . }}

	{{ if .MessageRules.GetSkip }}
		// skipping validation for {{ $f.Name }}
	{{ else }}
		if all {
			switch v := interface{}({{ accessor . }}).(type) {
				case interface{ ValidateAll() error }:
					if err := v.ValidateAll(); err != nil {
						errors = append(errors, {{ errCause . "err" "嵌入消息验证失败" }})
					}
				case interface{ Validate() error }:
					{{- /* Support legacy validation for messages that were generated with a plugin version prior to existence of ValidateAll() */ -}}
					if err := v.Validate(); err != nil {
						errors = append(errors, {{ errCause . "err" "嵌入消息验证失败" }})
					}
			}
		} else if v, ok := interface{}({{ accessor . }}).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return {{ errCause . "err" "嵌入消息验证失败" }}
			}
		}
	{{ end }}
`
