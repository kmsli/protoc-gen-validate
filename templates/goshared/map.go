package goshared

const mapTpl = `
	{{ $f := .Field }}{{ $r := .Rules }}

	{{ if $r.GetIgnoreEmpty }}
		if len({{ accessor . }}) > 0 {
	{{ end }}

	{{ if $r.GetMinPairs }}
		{{ if eq $r.GetMinPairs $r.GetMaxPairs }}
			if len({{ accessor . }}) != {{ $r.GetMinPairs }} {
				err := {{ err . "值必须正好包含 " $r.GetMinPairs " 键值对" }}
				if !all { return err }
				errors = append(errors, err)
			}
		{{ else if $r.MaxPairs }}
			if l := len({{ accessor . }}); l < {{ $r.GetMinPairs }} || l > {{ $r.GetMaxPairs }} {
				err := {{ err . "值必须包含在 " $r.GetMinPairs " 和 " $r.GetMaxPairs " 键值对之间" }}
				if !all { return err }
				errors = append(errors, err)
			}
		{{ else }}
			if len({{ accessor . }}) < {{ $r.GetMinPairs }} {
				err := {{ err . "值必须至少包含 " $r.GetMinPairs " 键值对" }}
				if !all { return err }
				errors = append(errors, err)
			}
		{{ end }}
	{{ else if $r.MaxPairs }}
		if len({{ accessor . }}) > {{ $r.GetMaxPairs }} {
			err := {{ err . "值不能超过 " $r.GetMaxPairs " 键值对" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}

	{{ if or $r.GetNoSparse (ne (.Elem "" "").Typ "none") (ne (.Key "" "").Typ "none") }}
		{{- /* Sort the keys to make the iteration order (and therefore failure output) deterministic. */ -}}
		{
			sorted_keys := make([]{{ (typ .Field).Key }}, len({{ accessor . }}))
			i := 0
			for key := range {{ accessor . }} {
				sorted_keys[i] = key
				i++
			}
			sort.Slice(sorted_keys, func (i, j int) bool { return sorted_keys[i] < sorted_keys[j] })
			for _, key := range sorted_keys {
				val := {{ accessor .}}[key]
				_ = val

				{{ if $r.GetNoSparse }}
					if val == nil {
						err := {{ errIdx . "键" "值不能是稀疏的, 所有键值对不能为空值" }}
						if !all { return err }
						errors = append(errors, err)
					}
				{{ end }}

				{{ render (.Key "key" "key") }}

				{{ render (.Elem "val" "key") }}
			}
		}
	{{ end }}

	{{ if $r.GetIgnoreEmpty }}
		}
	{{ end }}

`
