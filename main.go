// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package main

import (
	"flag"
	"fmt"

	"github.com/wantnotshould/byelog/cmd/flags"
)

func main() {
	flag.StringVar(&flags.Data, "data", "data", "Data directory")
	flag.BoolVar(&flags.Debug, "debug", false, "Enable debug mode")
	flag.Parse()

	fmt.Println("hello, byelog")
}
