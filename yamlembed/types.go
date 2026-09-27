package yamlembed

import (
	"strings"

//	"gopkg.in/yaml.v2"
)

type Foo struct {
	A string `yaml:"aa"`
	p int64  
}

type Bar struct {
	I      int64 `yaml:"-"`
	B      string
	UpperB string `yaml:"-"`
	OI     []string `yaml:"oi,omitempty,flow"`
	F      []any `yaml:"f,omitempty,flow"`
}

type Baz struct {
	Foo `yaml:",inline"`
	Bar `yaml:",inline"`
}

func (baz *Baz) UnmarshalYAML(unmarshal func(any) error) error {
	if err := unmarshal(&baz.Foo); err != nil {
		return err
	}
	if err := unmarshal(&baz.Bar); err != nil {
		return err
	}
	return nil
}

func (f *Foo) UnmarshalYAML(unmarshal func(any) error) error {
	type FooAlias Foo

	var foo FooAlias
    err := unmarshal(&foo)
	if err != nil {
		return err
	}

	*f = Foo(foo)
	return nil
}


func (b *Bar) UnmarshalYAML(unmarshal func(any) error) error {
	type BarAlias Bar

	var bar BarAlias
    err := unmarshal(&bar)
	if err != nil {
		return err
	}

	*b = Bar(bar)
	b.UpperB = strings.ToUpper(b.B)
	
	return nil
}
