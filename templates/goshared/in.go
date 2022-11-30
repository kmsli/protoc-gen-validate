package goshared

const inTpl = `{{ $f := .Field }}{{ $r := .Rules }}
	{{ if $r.In }}
		if _, ok := {{ lookup $f "InLookup" }}[{{ accessor . }}]; !ok {
			err := {{ err . "值必须在此列表中 " $r.In }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ else if $r.NotIn }}
		if _, ok := {{ lookup $f "NotInLookup" }}[{{ accessor . }}]; ok {
			err := {{ err . "值必须不在此列表中 " $r.NotIn }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}
`
