package goshared

const msgTpl = `
{{ if not (ignored .) -}}
{{ if disabled . -}}
	{{ cmt (msgTyp .) "的单字段验证(Validate)已禁用。此方法将始终返回nil" }}
{{- else -}}
	{{ cmt "单字段(Validate)检查 " (msgTyp .) "上具有此消息的原型定义中定义规则的字段值。 如果违反了任何规则，则返回遇到的第一个错误，如果没有违反，则返回零。" }}
{{- end -}}
func (m {{ (msgTyp .).Pointer }}) Validate() error {
	return m.validate(false)
}

{{ if disabled . -}}
	{{ cmt (msgTyp .) "的全字段验证(ValidateAll)已禁用。此方法将始终返回nil" }}
{{- else -}}
	{{ cmt "全字段(ValidateAll)检查" (msgTyp .) " 上具有此消息的原型定义中定义规则的字段值。 如果违反了任何规则, 结果是一个包含在 " (multierrname .) ", 中的违规错误列表，如果未找到，则为零。" }}
{{- end -}}
func (m {{ (msgTyp .).Pointer }}) ValidateAll() error {
	return m.validate(true)
}

{{/* Unexported function to handle validation. If the need arises to add more exported functions, please consider the functional option approach outlined in protoc-gen-validate#47. */}}
func (m {{ (msgTyp .).Pointer }}) validate(all bool) error {
	{{ if disabled . -}}
		return nil
	{{ else -}}
		if m == nil { return nil }

		var errors []error

		{{ range .NonOneOfFields }}
			{{ render (context .) }}
		{{ end }}

		{{ range .RealOneOfs }}
			{{- $oneof := . }}
			{{- if required . }}
			oneof{{ name $oneof }}Present := false
			{{- end }}
			switch v := m.{{ name . }}.(type) {
				{{- range .Fields }}
					{{- $context := (context .) }}
					case {{ oneof . }}:
						if v == nil {
							err := {{ errname .Message }}{
								field: "{{ name $oneof }}",
								reason: "oneof值不能是空类型",
							}
							if !all { return err }
							errors = append(errors, err)
						}
						{{- if required $oneof }}
						oneof{{ name $oneof }}Present = true
						{{- end }}
						{{ render $context }}
				{{- end }}
					default:
						_ = v // ensures v is used
			}
			{{- if required . }}
			if !oneof{{ name $oneof }}Present {
				err := {{ errname .Message }}{
					field: "{{ name $oneof }}",
					reason: "value is required",
				}
				if !all { return err }
				errors = append(errors, err)
			}
			{{- end }}
		{{- end }}

		{{ range .SyntheticOneOfFields }}
			if m.{{ name . }} != nil {
				{{ render (context .) }}
			}
		{{ end }}

		if len(errors) > 0 {
			return {{ multierrname . }}(errors)
		}

		return nil
	{{ end -}}
}

{{ if needs . "hostname" }}{{ template "hostname" . }}{{ end }}

{{ if needs . "email" }}{{ template "email" . }}{{ end }}

{{ if needs . "uuid" }}{{ template "uuid" . }}{{ end }}

{{ cmt (multierrname .) " 是一个包含多个验证问题的错误， 如果不满足指定的约束条件，被类型 " (msgTyp .) ".ValidateAll() 返回。" -}}
type {{ multierrname . }} []error

// Error returns a concatenation of all the error messages it wraps.
func (m {{ multierrname . }}) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m {{ multierrname . }}) AllErrors() []error { return m }

{{ cmt (errname .) " 是一个验证错误，如果不满足指定的约束条件，被类型 " (msgTyp .) ".Validate()返回" -}}
type {{ errname . }} struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e {{ errname . }}) Field() string { return e.field }

// Reason function returns reason value.
func (e {{ errname . }}) Reason() string { return e.reason }

// Cause function returns cause value.
func (e {{ errname . }}) Cause() error { return e.cause }

// Key function returns key value.
func (e {{ errname . }}) Key() bool { return e.key }

// ErrorName returns error name.
func (e {{ errname . }}) ErrorName() string { return "{{ errname . }}" }

// Error satisfies the builtin error interface
func (e {{ errname . }}) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %s{{ (msgTyp .) }}.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = {{ errname . }}{}

var _ interface{
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = {{ errname . }}{}

{{ range .Fields }}{{ with (context .) }}{{ $f := .Field }}
	{{ if has .Rules "In" }}{{ if .Rules.In }}
		var {{ lookup .Field "InLookup" }} = map[{{ inType .Field .Rules.In }}]struct{}{
			{{- range .Rules.In }}
				{{ inKey $f . }}: {},
			{{- end }}
		}
	{{ end }}{{ end }}

	{{ if has .Rules "NotIn" }}{{ if .Rules.NotIn }}
		var {{ lookup .Field "NotInLookup" }} = map[{{ inType .Field .Rules.In }}]struct{}{
			{{- range .Rules.NotIn }}
				{{ inKey $f . }}: {},
			{{- end }}
		}
	{{ end }}{{ end }}

	{{ if has .Rules "Pattern"}}{{ if .Rules.Pattern }}
		var {{ lookup .Field "Pattern" }} = regexp.MustCompile({{ lit .Rules.GetPattern }})
	{{ end }}{{ end }}

	{{ if has .Rules "Items"}}{{ if .Rules.Items }}
	{{ if has .Rules.Items.GetString_ "Pattern" }} {{ if .Rules.Items.GetString_.Pattern }}
		var {{ lookup .Field "Pattern" }} = regexp.MustCompile({{ lit .Rules.Items.GetString_.GetPattern }})
	{{ end }}{{ end }}
	{{ end }}{{ end }}

	{{ if has .Rules "Items"}}{{ if .Rules.Items }}
	{{ if has .Rules.Items.GetString_ "In" }} {{ if .Rules.Items.GetString_.In }}
		var {{ lookup .Field "InLookup" }} = map[string]struct{}{
			{{- range .Rules.Items.GetString_.In }}
				{{ inKey $f . }}: {},
			{{- end }}
		}
	{{ end }}{{ end }}
	{{ if has .Rules.Items.GetEnum "In" }} {{ if .Rules.Items.GetEnum.In }}
		var {{ lookup .Field "InLookup" }} = map[{{ inType .Field .Rules.Items.GetEnum.In }}]struct{}{
			{{- range .Rules.Items.GetEnum.In }}
				{{ inKey $f . }}: {},
			{{- end }}
		}
	{{ end }}{{ end }}
	{{ if has .Rules.Items.GetAny "In" }} {{ if .Rules.Items.GetAny.In }}
		var {{ lookup .Field "InLookup" }} = map[string]struct{}{
			{{- range .Rules.Items.GetAny.In }}
				{{ inKey $f . }}: {},
			{{- end }}
		}
	{{ end }}{{ end }}
	{{ end }}{{ end }}

	{{ if has .Rules "Items"}}{{ if .Rules.Items }}
	{{ if has .Rules.Items.GetString_ "NotIn" }} {{ if .Rules.Items.GetString_.NotIn }}
		var {{ lookup .Field "NotInLookup" }} = map[string]struct{}{
			{{- range .Rules.Items.GetString_.NotIn }}
				{{ inKey $f . }}: {},
			{{- end }}
		}
	{{ end }}{{ end }}
	{{ if has .Rules.Items.GetEnum "NotIn" }} {{ if .Rules.Items.GetEnum.NotIn }}
		var {{ lookup .Field "NotInLookup" }} = map[{{ inType .Field .Rules.Items.GetEnum.NotIn }}]struct{}{
			{{- range .Rules.Items.GetEnum.NotIn }}
				{{ inKey $f . }}: {},
			{{- end }}
		}
	{{ end }}{{ end }}
	{{ if has .Rules.Items.GetAny "NotIn" }} {{ if .Rules.Items.GetAny.NotIn }}
		var {{ lookup .Field "NotInLookup" }} = map[string]struct{}{
			{{- range .Rules.Items.GetAny.NotIn }}
				{{ inKey $f . }}: {},
			{{- end }}
		}
	{{ end }}{{ end }}
	{{ end }}{{ end }}

	{{ if has .Rules "Keys"}}{{ if .Rules.Keys }}
	{{ if has .Rules.Keys.GetString_ "Pattern" }} {{ if .Rules.Keys.GetString_.Pattern }}
		var {{ lookup .Field "Pattern" }} = regexp.MustCompile({{ lit .Rules.Keys.GetString_.GetPattern }})
	{{ end }}{{ end }}
	{{ end }}{{ end }}

	{{ if has .Rules "Values"}}{{ if .Rules.Values }}
	{{ if has .Rules.Values.GetString_ "Pattern" }} {{ if .Rules.Values.GetString_.Pattern }}
		var {{ lookup .Field "Pattern" }} = regexp.MustCompile({{ lit .Rules.Values.GetString_.GetPattern }})
	{{ end }}{{ end }}
	{{ end }}{{ end }}

{{ end }}{{ end }}
{{- end -}}
`
