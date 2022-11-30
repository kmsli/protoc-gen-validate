package golang

const requiredTpl = `
	{{ if .Rules.GetRequired }}
		if {{ accessor . }} == nil {
			err := {{ err . "值是必需的" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}
`
