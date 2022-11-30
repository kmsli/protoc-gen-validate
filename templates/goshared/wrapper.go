package goshared

const wrapperTpl = `
	{{ $f := .Field }}{{ $r := .Rules }}

	if wrapper := {{ accessor . }}; wrapper != nil {
		{{ render (unwrap . "wrapper") }}
	} {{ if .MessageRules.GetRequired }} else {
		err := {{ err . "值是必需的，不能为零。" }}
		if !all { return err }
		errors = append(errors, err)
	} {{ end }}
`
