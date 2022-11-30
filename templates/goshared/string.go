package goshared

const strTpl = `
	{{ $f := .Field }}{{ $r := .Rules }}

	{{ if $r.GetIgnoreEmpty }}
		if {{ accessor . }} != "" {
	{{ end }}

	{{ template "const" . }}
	{{ template "in" . }}

	{{ if or $r.Len (and $r.MinLen $r.MaxLen (eq $r.GetMinLen $r.GetMaxLen)) }}
		{{ if $r.Len }}
		if utf8.RuneCountInString({{ accessor . }}) != {{ $r.GetLen }} {
			err := {{ err . "值长度必须为 " $r.GetLen " runes" }}
			if !all { return err }
			errors = append(errors, err)
		{{ else }}
		if utf8.RuneCountInString({{ accessor . }}) != {{ $r.GetMinLen }} {
			err := {{ err . "值长度必须为 " $r.GetMinLen " runes" }}
			if !all { return err }
			errors = append(errors, err)
		{{ end }}
	}
	{{ else if $r.MinLen }}
		{{ if $r.MaxLen }}
			if l := utf8.RuneCountInString({{ accessor . }}); l < {{ $r.GetMinLen }} || l > {{ $r.GetMaxLen }} {
				err := {{ err . "值长度必须为 between " $r.GetMinLen " and " $r.GetMaxLen " runes, inclusive" }}
				if !all { return err }
				errors = append(errors, err)
			}
		{{ else }}
			if utf8.RuneCountInString({{ accessor . }}) < {{ $r.GetMinLen }} {
				err := {{ err . "值长度必须为 at least " $r.GetMinLen " runes" }}
				if !all { return err }
				errors = append(errors, err)
			}
		{{ end }}
	{{ else if $r.MaxLen }}
		if utf8.RuneCountInString({{ accessor . }}) > {{ $r.GetMaxLen }} {
			err := {{ err . "值长度必须为 at most " $r.GetMaxLen " runes" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}

	{{ if or $r.LenBytes (and $r.MinBytes $r.MaxBytes (eq $r.GetMinBytes $r.GetMaxBytes)) }}
		{{ if $r.LenBytes }}
			if len({{ accessor . }}) != {{ $r.GetLenBytes }} {
				err := {{ err . "值长度必须为 " $r.GetLenBytes " bytes" }}
				if !all { return err }
				errors = append(errors, err)
			}
		{{ else }}
			if len({{ accessor . }}) != {{ $r.GetMinBytes }} {
				err := {{ err . "值长度必须为 " $r.GetMinBytes " bytes" }}
				if !all { return err }
				errors = append(errors, err)
			}
		{{ end }}
	{{ else if $r.MinBytes }}
		{{ if $r.MaxBytes }}
			if l := len({{ accessor . }}); l < {{ $r.GetMinBytes }} || l > {{ $r.GetMaxBytes }} {
					err := {{ err . "值长度必须为 between " $r.GetMinBytes " and " $r.GetMaxBytes " bytes, inclusive" }}
					if !all { return err }
					errors = append(errors, err)
			}
		{{ else }}
			if len({{ accessor . }}) < {{ $r.GetMinBytes }} {
				err := {{ err . "值长度必须为 at least " $r.GetMinBytes " bytes" }}
				if !all { return err }
				errors = append(errors, err)
			}
		{{ end }}
	{{ else if $r.MaxBytes }}
		if len({{ accessor . }}) > {{ $r.GetMaxBytes }} {
			err := {{ err . "值长度必须为 at most " $r.GetMaxBytes " bytes" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}

	{{ if $r.Prefix }}
		if !strings.HasPrefix({{ accessor . }}, {{ lit $r.GetPrefix }}) {
			err := {{ err . "值缺少前缀 " (lit $r.GetPrefix) }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}

	{{ if $r.Suffix }}
		if !strings.HasSuffix({{ accessor . }}, {{ lit $r.GetSuffix }}) {
			err := {{ err . "值缺少后缀 " (lit $r.GetSuffix) }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}

	{{ if $r.Contains }}
		if !strings.Contains({{ accessor . }}, {{ lit $r.GetContains }}) {
			err := {{ err . "值不包含子字符 " (lit $r.GetContains) }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}

	{{ if $r.NotContains }}
		if strings.Contains({{ accessor . }}, {{ lit $r.GetNotContains }}) {
			err := {{ err . "值包含子字符 " (lit $r.GetNotContains) }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}

	{{ if $r.GetIp }}
		if strings.Contains({{ accessor . }}, "/") {
			if _, ipNet, _ := net.ParseCIDR({{ accessor . }}); ipNet == nil {
				err := {{ err . "值必须为有效的IP网段(CIDR)" }}
				if !all { return err }
				errors = append(errors, err)
			}
		} else {
			if ip := net.ParseIP({{ accessor . }}); ip == nil {
				err := {{ err . "值必须为有效的IP地址" }}
				if !all { return err }
				errors = append(errors, err)
			}
		}
	{{ else if $r.GetIpv4 }}
		if ip := net.ParseIP({{ accessor . }}); ip == nil || ip.To4() == nil {
			err := {{ err . "值必须为有效的IPv4地址" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ else if $r.GetIpv6 }}
		if ip := net.ParseIP({{ accessor . }}); ip == nil || ip.To4() != nil {
			err := {{ err . "值必须为有效的IPv6" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ else if $r.GetEmail }}
		if err := m._validateEmail({{ accessor . }}); err != nil {
			err = {{ errCause . "err" "值必须为有效的邮箱地址" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ else if $r.GetHostname }}
		if err := m._validateHostname({{ accessor . }}); err != nil {
			err = {{ errCause . "err" "值必须是有效的主机名" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ else if $r.GetAddress }}
		if err := m._validateHostname({{ accessor . }}); err != nil {
			if ip := net.ParseIP({{ accessor . }}); ip == nil {
				err := {{ err . "值必须是有效的主机名或ip地址" }}
				if !all { return err }
				errors = append(errors, err)
			}
		}
	{{ else if $r.GetUri }}
		if uri, err := url.Parse({{ accessor . }}); err != nil {
			err = {{ errCause . "err" "值必须是有效的URI" }}
			if !all { return err }
			errors = append(errors, err)
		} else if !uri.IsAbs() {
			err := {{ err . "值必须是绝对值" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ else if $r.GetUriRef }}
		if _, err := url.Parse({{ accessor . }}); err != nil {
			err = {{ errCause . "err" "值必须是有效的URI" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ else if $r.GetUuid }}
		if err := m._validateUuid({{ accessor . }}); err != nil {
			err = {{ errCause . "err" "值必须是有效的UUID" }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}

	{{ if $r.Pattern }}
		if !{{ lookup $f "Pattern" }}.MatchString({{ accessor . }}) {
			err := {{ err . "值与正则表达式模式不匹配 " (lit $r.GetPattern) }}
			if !all { return err }
			errors = append(errors, err)
		}
	{{ end }}

	{{ if $r.GetIgnoreEmpty }}
		}
	{{ end }}

`
