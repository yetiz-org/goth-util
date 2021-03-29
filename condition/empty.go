package condition

func NonNilDo(par interface{}, f func(interface{})) {
	if par != nil {
		f(par)
	}
}

func NilDo(par interface{}, f func()) {
	if par == nil {
		f()
	}
}

func NonNilDoReturn(par interface{}, f func(interface{}) interface{}) interface{} {
	if par != nil {
		return f(par)
	}

	return nil
}

func NilDoReturn(par interface{}, f func() interface{}) interface{} {
	if par == nil {
		return f()
	}

	return nil
}

func NonEmptyStringDo(str string, f func(string)) {
	if str != "" {
		f(str)
	}
}

func EmptyStringDo(str string, f func()) {
	if str == "" {
		f()
	}
}

func NonEmptyStringDoReturn(str string, f func(string) interface{}) interface{} {
	if str != "" {
		return f(str)
	}

	return nil
}

func EmptyStringDoReturn(str string, f func() interface{}) interface{} {
	if str == "" {
		return f()
	}

	return nil
}
