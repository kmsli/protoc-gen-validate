package goshared

const hostTpl = `
	func (m {{ (msgTyp .).Pointer }}) _validateHostname(host string) error {
		s := strings.ToLower(strings.TrimSuffix(host, "."))

		if len(host) > 253 {
			return errors.New("主机名不能超过253个字符")
		}

		for _, part := range strings.Split(s, ".") {
			if l := len(part); l == 0 || l > 63 {
				return errors.New("主机名部分必须为非空且不能超过63个字符")
			}

			if part[0] == '-' {
				return errors.New("主机名部分不能以连字符开头")
			}

			if part[len(part)-1] == '-' {
				return errors.New("主机名部分不能以连字符结尾")
			}

			for _, r := range part {
				if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
					return fmt.Errorf("主机名部分只能包含字母数字字符或连字符, 取值 %q", string(r))
				}
			}
		}

		return nil
	}
`

const emailTpl = `
	func (m {{ (msgTyp .).Pointer }}) _validateEmail(addr string) error {
		a, err := mail.ParseAddress(addr)
		if err != nil {
			return err
		}
		addr = a.Address

		if len(addr) > 254 {
			return errors.New("email addresses cannot exceed 254 characters")
		}

		parts := strings.SplitN(addr, "@", 2)

		if len(parts[0]) > 64 {
			return errors.New("email address local phrase cannot exceed 64 characters")
		}

		return m._validateHostname(parts[1])
	}
`

const uuidTpl = `
	func (m {{ (msgTyp .).Pointer }}) _validateUuid(uuid string) error {
		if matched := _{{ snakeCase .File.InputPath.BaseName }}_uuidPattern.MatchString(uuid); !matched {
			return errors.New("invalid uuid format")
		}

		return nil
	}
`
