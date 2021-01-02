// Copyright ©2019 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yaegiconf_test

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/traefik/yaegi/interp"

	"github.com/kortschak/yaegiconf"
)

func Example_struct() {
	type Config struct {
		N int
		F float64
	}
	var c Config
	err := yaegiconf.EvalTo(&c, `config.Value{N: 5, F: 0.1}`)
	if err != nil {
		log.Fatalf("failed to parse configuration: %v", err)
	}

	fmt.Printf("%#v\n", c)

	// Output:
	//
	// yaegiconf_test.Config{N:5, F:0.1}
}

func Example_nestedstruct() {
	type Sub struct {
		S string
	}
	type Config struct {
		N int
		F float64
		X Sub
	}
	s := interp.Exports{"xcfg": map[string]reflect.Value{
		"Value": reflect.Zero(reflect.TypeOf(&Config{})),
		"Sub":   reflect.Zero(reflect.TypeOf(&Sub{})),
	}}
	var c Config
	err := yaegiconf.EvalExtWithContextTo(context.Background(), &c, `xcfg.Value{
		N: 5, F: 0.1,
		X: xcfg.Sub{S: "set"},
}`,
		interp.Options{}, s)
	if err != nil {
		log.Fatalf("failed to parse configuration: %v", err)
	}

	fmt.Printf("%#v\n", c)

	// Output:
	//
	// yaegiconf_test.Config{N:5, F:0.1, X:yaegiconf_test.Sub{S:"set"}}
}

func Example_string() {
	type Config string
	var c Config
	err := yaegiconf.EvalTo(&c, `config.Value("Configured")`)
	if err != nil {
		log.Fatalf("failed to parse configuration: %v", err)
	}

	fmt.Printf("%#v\n", c)

	// Output:
	//
	// "Configured"
}

func Example_map() {
	type Config map[string]interface{}
	var c Config
	err := yaegiconf.EvalTo(&c, `config.Value{"int": 5, "float64": 0.1}`)
	if err != nil {
		log.Fatalf("failed to parse configuration: %v", err)
	}

	fmt.Printf("%#v\n", c)

	// Output:
	//
	// yaegiconf_test.Config{"float64":0.1, "int":5}
}

func Example_timeout() {
	type Config string
	var c Config
	err := yaegiconf.EvalTo(&c, `for {}; config.Value("Configured")`)
	if err != nil {
		fmt.Printf("failed to parse configuration: %v", err)
	}

	fmt.Println(c)

	// Output:
	//
	// failed to parse configuration: context deadline exceeded
}
