package xtype

import "strings"

type String string

func (s String) String() string {
	return string(s)
}

func (s String) Upper() string {
	return strings.ToUpper(string(s))
}

func (s String) Lower() string {
	return strings.ToLower(string(s))
}

func (s String) Short() string {
	return s.ShortBySign("-")
}

func (s String) ShortBySign(sign string) string {
	sb := ""
	for _, sp := range strings.Split(strings.ToUpper(string(s)), sign) {
		if len(sp) > 0 {
			sb = sb + sp[0:1]
		}
	}

	return sb
}
