// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"context"
	"fmt"

	"github.com/attic-labs/noms/cmd/util"
	flag "github.com/juju/gnuflag"
	"github.com/liquidata-inc/ld/dolt/go/store/go/config"
	"github.com/liquidata-inc/ld/dolt/go/store/go/d"
	"github.com/liquidata-inc/ld/dolt/go/store/go/diff"
	"github.com/liquidata-inc/ld/dolt/go/store/go/util/outputpager"
	"github.com/liquidata-inc/ld/dolt/go/store/go/util/verbose"
)

var stat bool

var nomsDiff = &util.Command{
	Run:       runDiff,
	UsageLine: "diff [--stat] <object1> <object2>",
	Short:     "Shows the difference between two objects",
	Long:      "See Spelling Objects at https://github.com/attic-labs/noms/blob/master/doc/spelling.md for details on the object arguments.",
	Flags:     setupDiffFlags,
	Nargs:     2,
}

func setupDiffFlags() *flag.FlagSet {
	diffFlagSet := flag.NewFlagSet("diff", flag.ExitOnError)
	diffFlagSet.BoolVar(&stat, "stat", false, "Writes a summary of the changes instead")
	outputpager.RegisterOutputpagerFlags(diffFlagSet)
	verbose.RegisterVerboseFlags(diffFlagSet)

	return diffFlagSet
}

func runDiff(ctx context.Context, args []string) int {
	cfg := config.NewResolver()
	db1, value1, err := cfg.GetPath(ctx, args[0])
	d.CheckErrorNoUsage(err)
	if value1 == nil {
		d.CheckErrorNoUsage(fmt.Errorf("Object not found: %s", args[0]))
	}
	defer db1.Close()

	db2, value2, err := cfg.GetPath(ctx, args[1])
	d.CheckErrorNoUsage(err)
	if value2 == nil {
		d.CheckErrorNoUsage(fmt.Errorf("Object not found: %s", args[1]))
	}
	defer db2.Close()

	if stat {
		diff.Summary(ctx, value1, value2)
		return 0
	}

	pgr := outputpager.Start()
	defer pgr.Stop()

	diff.PrintDiff(ctx, pgr.Writer, value1, value2, false)
	return 0
}
