// Copyright ©2026 cdme. All rights reserved.
// Author: https://cdme.cn
// Email: hi@cdme.cn

package main

import (
	"flag"
	"fmt"

	"github.com/wantnotshould/byelog/cmd/flags"
	"github.com/wantnotshould/byelog/internal/bootstrap"
)

func main() {
	flag.StringVar(&flags.Data, "data", "data", "Data directory")
	flag.BoolVar(&flags.Debug, "debug", false, "Enable debug mode")
	flag.Parse()

	bootstrap.Run()
	defer bootstrap.Release()

	fmt.Println("hello, byelog")
}
