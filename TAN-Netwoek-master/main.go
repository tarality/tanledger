package main

import (
	_ "embed"

	"github.com/tarality/tan-network/command/root"
	"github.com/tarality/tan-network/licenses"
)

var (
	//go:embed LICENSE
	license string
)

func main() {
	licenses.SetLicense(license)

	root.NewRootCommand().Execute()
}
