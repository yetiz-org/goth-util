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
	sb := ""
	for _, sp := range strings.Split(strings.ToUpper(string(s)), "-") {
		if len(sp) > 0 {
			sb = sb + sp[0:1]
		}
	}

	return sb
}
