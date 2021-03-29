package asif

type StringDef struct {
	fireThen bool
	fireElse bool
	base     string
}

func String(str string) *StringDef {
	return &StringDef{base: str}
}

func StringDefault(str string, defaultStr string) *StringDef {
	if str == "" {
		return &StringDef{base: defaultStr}
	}

	return String(str)
}

func (sd *StringDef) Or(str string) *StringDef {
	if sd.base == "" {
		sd.base = str
	}

	return sd
}

func (sd *StringDef) IsEmpty() bool {
	if sd.base == "" {
		return true
	}

	return false
}

func (sd *StringDef) IsEqual(str string) bool {
	return sd.base == str
}

func (sd *StringDef) Val() string {
	return sd.base
}

func (sd *StringDef) Empty() *StringDef {
	sd.fireThen = sd.base == ""
	sd.fireElse = !sd.fireThen
	return sd
}

func (sd *StringDef) Equal(str string) *StringDef {
	sd.fireThen = sd.base == str
	sd.fireElse = !sd.fireThen
	return sd
}

func (sd *StringDef) Then(f func(sd *StringDef)) *StringDef {
	if f != nil && sd.fireThen {
		f(sd)
	}

	return sd
}

func (sd *StringDef) Else(f func(sd *StringDef)) *StringDef {
	if f != nil && sd.fireElse {
		f(sd)
	}

	return sd
}
