//go:build ignore && gopher

package main

import (
	"context"
	"os"
	"time"

	. "github.com/ohhfishal/gopher/runtime"
)

var packages = []string{
	"./assert",
	"./color",
	"./log",
}

// Devel builds the gopher binary then runs it
func Devel(ctx context.Context, gopher *Gopher) error {
	// TODO: Don't assume gopher is in the path?
	if err := InstallGitHook(gopher.Stdout, GitPreCommit, "gopher cicd"); err != nil {
		return err
	}
	var status Status
	return gopher.Run(ctx, NowAnd(OnFileChange(1*time.Second, ".go")),
		status.Start(),
		&GoBuild{
			Packages: packages,
		},
		&GoFormat{},
		&GoTest{
			Packages: packages,
		},
		&GoVet{},
		&GoModTidy{},
		ExecCommand("go", "run", "."),
		status.Done(),
	)
}

// cicd runs the entire ci/cd suite
func CICD(ctx context.Context, gopher *Gopher) error {
	var status Status
	return gopher.Run(ctx, Now(),
		status.Start(),
		&GoBuild{
			Packages: packages,
		},
		&GoFormat{
			CheckOnly: true,
		},
		&GoTest{
			Packages: packages,
		},
		&GoVet{},
		status.Done(),
	)
}

// Removes all local build artifacts.
func Clean(ctx context.Context, gopher *Gopher) error {
	return os.RemoveAll("target")
}
