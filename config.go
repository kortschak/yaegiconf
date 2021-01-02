// Copyright ©2019 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package yaegiconf provides a simple interface to the yaegi
// interpreter for use as a configuration parser.
package yaegiconf

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/traefik/yaegi/interp"
)

// EvalTo evaluates the configuration source in src and stores the
// result of the evaluation into dst, which must be a pointer to
// a value. The type of the value is accessible within the src via
// the label config.Value.
//
// EvalTo will timeout if it has not completed evaluating src within
// ten seconds.
func EvalTo(dst interface{}, src string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return EvalWithContextTo(ctx, dst, src)
}

// EvalWithContextTo evaluates the configuration source in src and stores the
// result of the evaluation into dst, which must be a pointer to a value.
//
// EvalWithContextTo allows the user to provide a cancellation context to the
// evaluation.
func EvalWithContextTo(ctx context.Context, dst interface{}, src string) error {
	return EvalExtWithContextTo(ctx, dst, src, interp.Options{}, interp.Exports{
		"config": map[string]reflect.Value{
			"Value": reflect.Zero(reflect.TypeOf(dst))}})
}

// EvalExtWithContextTo evaluates the configuration source in src and stores the
// result of the evaluation into dst, which must be a pointer to a value.
//
// EvalExtWithContextTo allows the user to provide a cancellation context and
// access to additional types and interpreter configuration.
func EvalExtWithContextTo(ctx context.Context, dst interface{}, src string, options interp.Options, symbols ...interp.Exports) error {
	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr {
		return errors.New("yaegiconf: invalid config type")
	}

	i := interp.New(options)
	for _, m := range symbols {
		i.Use(m)
		for p := range m {
			_, err := i.EvalWithContext(ctx, fmt.Sprintf("import %q", p))
			if err != nil {
				return err
			}
		}
	}

	v, err := i.EvalWithContext(ctx, src)
	if err != nil {
		return err
	}

	if !v.IsValid() {
		return errors.New("yaegiconf: no configuration value in src")
	}
	rv = rv.Elem()
	if v.Type() != rv.Type() {
		return fmt.Errorf("yaegiconf: cannot use src (type %v) as type %v in configuration", v.Type(), rv.Type())
	}
	rv.Set(v)
	return nil
}
