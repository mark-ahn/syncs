package syncs

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
)

type context_args int

const (
	context_args_map context_args = iota
)

type Value interface {
	String() string
	Set(string) error
}

type ArgMap map[string]string

func WithArgMap(ctx context.Context, argMap ArgMap) context.Context {
	return context.WithValue(ctx, context_args_map, argMap)
}

func ArgMapFrom(ctx context.Context) ArgMap {
	d := ctx.Value(context_args_map)
	if d == nil {
		return ArgMap{}
	}
	argMap, ok := d.(ArgMap)
	if !ok {
		return ArgMap{}
	}
	return argMap
}

type parse_item struct {
	key   string
	value reflect.Value
	def   reflect.Value
}

type Parser struct {
	arg_map ArgMap
	items   []parse_item
}

func (__ *Parser) Var(p Value, key string, desc string) {
	__.items = append(__.items, parse_item{
		key:   key,
		value: reflect.ValueOf(p),
	})
}

func (__ *Parser) IntVar(p *int, key string, def int, desc string) {
	__.items = append(__.items, parse_item{
		key:   key,
		value: reflect.ValueOf(p).Elem(),
		def:   reflect.ValueOf(def),
	})
}

func (__ *Parser) FloatVar(p *float64, key string, def float64, desc string) {
	__.items = append(__.items, parse_item{
		key:   key,
		value: reflect.ValueOf(p).Elem(),
		def:   reflect.ValueOf(def),
	})
}

func (__ *Parser) BoolVar(p *bool, key string, def bool, desc string) {
	__.items = append(__.items, parse_item{
		key:   key,
		value: reflect.ValueOf(p).Elem(),
		def:   reflect.ValueOf(def),
	})
}

func (__ *Parser) StringVar(p *string, key string, def string, desc string) {
	__.items = append(__.items, parse_item{
		key:   key,
		value: reflect.ValueOf(p).Elem(),
		def:   reflect.ValueOf(def),
	})
}

func (__ *Parser) Parse() error {
	var err error
	match_check_map := make(map[string]struct{}, len(__.items))
	for _, d := range __.items {
		var value interface{}
		switch d.value.Kind() {
		case reflect.Bool:
			value, err = __.arg_map.Bool(d.key)
		case reflect.Int:
			value, err = __.arg_map.Int(d.key)
		case reflect.Float64:
			value, err = __.arg_map.Float(d.key)
		case reflect.String:
			value, err = __.arg_map.String(d.key)
		default:
			v_intf, ok := d.value.Interface().(Value)
			if !ok {
				return fmt.Errorf("value does not implement Value")
			}
			v_str := __.arg_map[d.key]
			err := v_intf.Set(v_str)
			if err != nil {
				return err
			}
		}
		switch err.(type) {
		case nil:
		case NoKeyInArgMapError:
			value, err = d.def.Interface(), nil
		default:
			return err
		}
		d.value.Set(reflect.ValueOf(value))
		match_check_map[d.key] = struct{}{}
	}
	for k := range __.arg_map {
		_, ok := match_check_map[k]
		if !ok {
			return fmt.Errorf("unknown argument key %s", k)
		}
	}
	return nil
}

func (__ ArgMap) Parser() *Parser {
	return &Parser{
		arg_map: __,
	}
}

func (__ ArgMap) Int(key string) (int, error) {
	var res int
	v, ok := __[key]
	if !ok {
		return res, NoKeyInArgMapErrorOf(key)
	}
	s64, err := strconv.ParseInt(v, 0, 0)
	if err != nil {
		return res, err
	}
	res = int(s64)
	return res, nil
}

func (__ ArgMap) Float(key string) (float64, error) {
	var res float64
	v, ok := __[key]
	if !ok {
		return res, NoKeyInArgMapErrorOf(key)
	}
	f64, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return res, err
	}
	res = f64
	return res, nil
}

func (__ ArgMap) String(key string) (string, error) {
	var res string
	v, ok := __[key]
	if !ok {
		return res, NoKeyInArgMapErrorOf(key)
	}
	res = v
	return res, nil
}

func (__ ArgMap) Bool(key string) (bool, error) {
	var res bool
	v, ok := __[key]
	if !ok {
		return res, NoKeyInArgMapErrorOf(key)
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return res, err
	}
	res = b
	return res, nil
}
