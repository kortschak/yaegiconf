// Copyright ©2019 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yaegiconf_test

import (
	"fmt"
	"log"

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
