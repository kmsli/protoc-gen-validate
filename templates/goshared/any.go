package goshared

const anyTpl = `{{ $f := .Field }}{{ $r := .Rules }}
	{{ template "required" . }}

	if a := {{ accessor . }}; a != nil {
		{{ if $r.In }}
			if _, ok := {{ lookup $f "InLookup" }}[a.GetTypeUrl()]; !ok {
				err := {{ err . "类型 URL 必须在列表中 " $r.In }}
				if !all { return err }
				errors = append(errors, err)
			}
		{{ else if $r.NotIn }}
			if _, ok := {{ lookup $f "NotInLookup" }}[a.GetTypeUrl()]; ok {
				err := {{ err . "类型 URL 不能在列表中 " $r.NotIn }}
				if !all { return err }
				errors = append(errors, err)
			}
		{{ end }}
	}
`
