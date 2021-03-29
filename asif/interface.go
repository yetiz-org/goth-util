package asif

import "reflect"

type InterfaceDef struct {
	fireThen bool
	fireElse bool
	base     interface{}
}

func Interface(intrf interface{}) *InterfaceDef {
	return &InterfaceDef{base: intrf}
}

func InterfaceDefault(intrf interface{}, defaultIntrf interface{}) *InterfaceDef {
	if reflect.ValueOf(intrf).IsNil() {
		return &InterfaceDef{base: defaultIntrf}
	}

	return Interface(intrf)
}

func (sd *InterfaceDef) Or(intrf interface{}) *InterfaceDef {
	if reflect.ValueOf(sd.base).IsNil() {
		sd.base = intrf
	}

	return sd
}

func (sd *InterfaceDef) IsEmpty() bool {
	if reflect.ValueOf(sd.base).IsNil() {
		return true
	}

	return false
}

func (sd *InterfaceDef) IsEqual(intrf interface{}) bool {
	return sd.base == intrf
}

func (sd *InterfaceDef) Val() interface{} {
	return sd.base
}

func (sd *InterfaceDef) Empty() *InterfaceDef {
	sd.fireThen = reflect.ValueOf(sd.base).IsNil()
	sd.fireElse = !sd.fireThen
	return sd
}

func (sd *InterfaceDef) Equal(intrf interface{}) *InterfaceDef {
	sd.fireThen = sd.base == intrf
	sd.fireElse = !sd.fireThen
	return sd
}

func (sd *InterfaceDef) Then(f func(sd *InterfaceDef)) *InterfaceDef {
	if f != nil && sd.fireThen {
		f(sd)
	}

	return sd
}

func (sd *InterfaceDef) Else(f func(sd *InterfaceDef)) *InterfaceDef {
	if f != nil && sd.fireElse {
		f(sd)
	}

	return sd
}
