package goshared

const timestampcmpTpl = `{{ $f := .Field }}{{ $r := .Rules }}
			{{  if $r.Const }}
				if !ts.Equal({{ tsLit $r.Const }}) {
					err := {{ err . "value must equal " (tsStr $r.Const) }}
					if !all { return err }
					errors = append(errors, err)
				}
			{{ end }}

			{{ if or $r.LtNow $r.GtNow $r.Within }} now := time.Now(); {{ end }}
			{{- if $r.Lt }}  lt  := {{ tsLit $r.Lt }};  {{ end }}
			{{- if $r.Lte }} lte := {{ tsLit $r.Lte }}; {{ end }}
			{{- if $r.Gt }}  gt  := {{ tsLit $r.Gt }};  {{ end }}
			{{- if $r.Gte }} gte := {{ tsLit $r.Gte }}; {{ end }}
			{{- if $r.Within }} within := {{ durLit $r.Within }}; {{ end }}

			{{ if $r.Lt }}
				{{ if $r.Gt }}
					{{  if tsGt $r.GetLt $r.GetGt }}
						if ts.Sub(gt) <= 0 || ts.Sub(lt) >= 0 {
							err := {{ err . "值必须在范围内 (" (tsStr $r.GetGt) ", " (tsStr $r.GetLt) ")" }}
							if !all { return err }
							errors = append(errors, err)
						}
					{{ else }}
						if ts.Sub(lt) >= 0 && ts.Sub(gt) <= 0 {
							err := {{ err . "值必须超出范围 [" (tsStr $r.GetLt) ", " (tsStr $r.GetGt) "]" }}
							if !all { return err }
							errors = append(errors, err)
						}
					{{ end }}
				{{ else if $r.Gte }}
					{{  if tsGt $r.GetLt $r.GetGte }}
						if ts.Sub(gte) < 0 || ts.Sub(lt) >= 0 {
							err := {{ err . "值必须在范围内 [" (tsStr $r.GetGte) ", " (tsStr $r.GetLt) ")" }}
							if !all { return err }
							errors = append(errors, err)
						}
					{{ else }}
						if ts.Sub(lt) >= 0 && ts.Sub(gte) < 0 {
							err := {{ err . "值必须超出范围 [" (tsStr $r.GetLt) ", " (tsStr $r.GetGte) ")" }}
							if !all { return err }
							errors = append(errors, err)
						}
					{{ end }}
				{{ else }}
					if ts.Sub(lt) >= 0 {
						err := {{ err . "值必须小于 " (tsStr $r.GetLt) }}
						if !all { return err }
						errors = append(errors, err)
					}
				{{ end }}
			{{ else if $r.Lte }}
				{{ if $r.Gt }}
					{{  if tsGt $r.GetLte $r.GetGt }}
						if ts.Sub(gt) <= 0 || ts.Sub(lte) > 0 {
							err := {{ err . "值必须在范围内 (" (tsStr $r.GetGt) ", " (tsStr $r.GetLte) "]" }}
							if !all { return err }
							errors = append(errors, err)
						}
					{{ else }}
						if ts.Sub(lte) > 0 && ts.Sub(gt) <= 0 {
							err := {{ err . "值必须超出范围 (" (tsStr $r.GetLte) ", " (tsStr $r.GetGt) "]" }}
							if !all { return err }
							errors = append(errors, err)
						}
					{{ end }}
				{{ else if $r.Gte }}
					{{ if tsGt $r.GetLte $r.GetGte }}
						if ts.Sub(gte) < 0 || ts.Sub(lte) > 0 {
							err := {{ err . "值必须在范围内 [" (tsStr $r.GetGte) ", " (tsStr $r.GetLte) "]" }}
							if !all { return err }
							errors = append(errors, err)
						}
					{{ else }}
						if ts.Sub(lte) > 0 && ts.Sub(gte) < 0 {
							err := {{ err . "值必须超出范围 (" (tsStr $r.GetLte) ", " (tsStr $r.GetGte) ")" }}
							if !all { return err }
							errors = append(errors, err)
						}
					{{ end }}
				{{ else }}
					if ts.Sub(lte) > 0 {
						err := {{ err . "值必须小于 or equal to " (tsStr $r.GetLte) }}
						if !all { return err }
						errors = append(errors, err)
					}
				{{ end }}
			{{ else if $r.Gt }}
				if ts.Sub(gt) <= 0 {
					err := {{ err . "值必须大于 " (tsStr $r.GetGt) }}
					if !all { return err }
					errors = append(errors, err)
				}
			{{ else if $r.Gte }}
				if ts.Sub(gte) < 0 {
					err := {{ err . "值必须大于 or equal to " (tsStr $r.GetGte) }}
					if !all { return err }
					errors = append(errors, err)
				}
			{{ else if $r.LtNow }}
				{{ if $r.Within }}
					if ts.Sub(now) >= 0 || ts.Sub(now.Add(-within)) < 0 {
						err := {{ err . "值必须小于 now within " (durStr $r.GetWithin) }}
						if !all { return err }
						errors = append(errors, err)
					}
				{{ else }}
					if ts.Sub(now) >= 0 {
						err := {{ err . "值必须小于 now" }}
						if !all { return err }
						errors = append(errors, err)
					}
				{{ end }}
			{{ else if $r.GtNow }}
				{{ if $r.Within }}
					if ts.Sub(now) <= 0 || ts.Sub(now.Add(within)) > 0 {
						err := {{ err . "值必须大于 now within " (durStr $r.GetWithin) }}
						if !all { return err }
						errors = append(errors, err)
					}
				{{ else }}
					if ts.Sub(now) <= 0 {
						err := {{ err . "值必须大于 now" }}
						if !all { return err }
						errors = append(errors, err)
					}
				{{ end }}
			{{ else if $r.Within }}
				if ts.Sub(now.Add(within)) >= 0 || ts.Sub(now.Add(-within)) <= 0 {
					err := {{ err . "value must be within " (durStr $r.GetWithin) " of now" }}
					if !all { return err }
					errors = append(errors, err)
				}
			{{ end }}
`
