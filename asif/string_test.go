package asif

import (
	"reflect"
	"testing"
)

func TestString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want *StringDef
	}{
		{name: "empty", args: args{str: ""}, want: &StringDef{base: ""}},
		{name: "str", args: args{str: "str"}, want: &StringDef{base: "str"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringDefault(t *testing.T) {
	type args struct {
		str        string
		defaultStr string
	}
	tests := []struct {
		name string
		args args
		want *StringDef
	}{
		{name: "1", args: args{str: "1", defaultStr: ""}, want: &StringDef{base: "1"}},
		{name: "empty", args: args{str: "", defaultStr: "e"}, want: &StringDef{base: "e"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringDefault(tt.args.str, tt.args.defaultStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringDef_Else(t *testing.T) {
	type fields struct {
		fireThen bool
		fireElse bool
		base     string
	}
	type args struct {
		f func(sd *StringDef)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "e", fields: fields{base: ""}, args: args{f: func(sd *StringDef) { t.Error("") }}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := &StringDef{
				fireThen: tt.fields.fireThen,
				fireElse: tt.fields.fireElse,
				base:     tt.fields.base,
			}

			sd.Empty().Else(tt.args.f)
		})
	}
}

func TestStringDef_Empty(t *testing.T) {
	type fields struct {
		fireThen bool
		fireElse bool
		base     string
	}
	tests := []struct {
		name   string
		fields fields
		want   *StringDef
	}{
		{name: "e", fields: fields{base: ""}, want: &StringDef{base: "", fireThen: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := &StringDef{
				fireThen: tt.fields.fireThen,
				fireElse: tt.fields.fireElse,
				base:     tt.fields.base,
			}
			if got := sd.Empty(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringDef_Equal(t *testing.T) {
	type fields struct {
		fireThen bool
		fireElse bool
		base     string
	}
	type args struct {
		def string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *StringDef
	}{
		{name: "1", fields: fields{base: "1"}, args: args{def: "1"}, want: &StringDef{base: "1", fireThen: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := &StringDef{
				fireThen: tt.fields.fireThen,
				fireElse: tt.fields.fireElse,
				base:     tt.fields.base,
			}
			if got := sd.Equal(tt.args.def); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringDef_IsEmpty(t *testing.T) {
	type fields struct {
		fireThen bool
		fireElse bool
		base     string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "e", fields: fields{base: ""}, want: true},
		{name: "1", fields: fields{base: "1"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := &StringDef{
				fireThen: tt.fields.fireThen,
				fireElse: tt.fields.fireElse,
				base:     tt.fields.base,
			}
			if got := sd.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringDef_IsEqual(t *testing.T) {
	type fields struct {
		fireThen bool
		fireElse bool
		base     string
	}
	type args struct {
		def string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{name: "1", fields: fields{base: "1"}, args: args{def: "1"}, want: true},
		{name: "1", fields: fields{base: "1"}, args: args{def: "2"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := &StringDef{
				fireThen: tt.fields.fireThen,
				fireElse: tt.fields.fireElse,
				base:     tt.fields.base,
			}
			if got := sd.IsEqual(tt.args.def); got != tt.want {
				t.Errorf("IsEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringDef_Or(t *testing.T) {
	type fields struct {
		fireThen bool
		fireElse bool
		base     string
	}
	type args struct {
		str string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *StringDef
	}{
		{name: "e", fields: fields{base: ""}, args: args{str: "2"}, want: &StringDef{base: "2"}},
		{name: "1", fields: fields{base: "1"}, args: args{str: "2"}, want: &StringDef{base: "1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := &StringDef{
				fireThen: tt.fields.fireThen,
				fireElse: tt.fields.fireElse,
				base:     tt.fields.base,
			}
			if got := sd.Or(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Or() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringDef_Then(t *testing.T) {
	type fields struct {
		fireThen bool
		fireElse bool
		base     string
	}
	type args struct {
		f func(sd *StringDef)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "1", fields: fields{base: "1"}, args: args{f: func(sd *StringDef) { t.Error() }}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := &StringDef{
				fireThen: tt.fields.fireThen,
				fireElse: tt.fields.fireElse,
				base:     tt.fields.base,
			}

			sd.Empty().Then(tt.args.f)
		})
	}
}

func TestStringDef_Val(t *testing.T) {
	type fields struct {
		fireThen bool
		fireElse bool
		base     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "e", fields: fields{base: ""}, want: ""},
		{name: "1", fields: fields{base: "1"}, want: "1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := &StringDef{
				fireThen: tt.fields.fireThen,
				fireElse: tt.fields.fireElse,
				base:     tt.fields.base,
			}
			if got := sd.Val(); got != tt.want {
				t.Errorf("Val() = %v, want %v", got, tt.want)
			}
		})
	}
}
