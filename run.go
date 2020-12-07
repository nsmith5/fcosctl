package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func compile(config string) (ign string, cleanup func(), err error) {
	var dir string
	dir, err = ioutil.TempDir(``, `fcosctl-*`)
	if err != nil {
		return
	}

	cleanup = func() { os.RemoveAll(dir) }

	ign = filepath.Join(dir, `ign.json`)
	cmd := exec.Command(`fcct`, `-d`, `.`, `--strict`, `--output`, ign, config)
	_, err = cmd.Output()
	return
}

func run(config, version string) {
	// Some Printers
	errf := color.New(color.FgRed).PrintfFunc()
	infof := color.New(color.FgCyan).PrintfFunc()

	// Pre-condition: verify the file config file exists
	if _, err := os.Stat(config); os.IsNotExist(err) {
		errf("Config %s does not exist\n", config)
		os.Exit(1)
	}

	infof("⚒️  Compiling config...\n")
	ign, cleanup, err := compile(config)
	if err != nil {
		exitError, wasExitError := err.(*exec.ExitError)
		if wasExitError {
			cleanup()
			errf("😬 Failed to compile config with error:\n%s\n", string(exitError.Stderr))
			os.Exit(1)
		}
		panic(err)
	}
	defer cleanup()

	imageFile := getImagePath(version)
	cmd := exec.Command(
		`qemu-kvm`,
		`-m`, `2048`,
		`-cpu`, `host`,
		`-nographic`,
		`-snapshot`,
		`-drive`, fmt.Sprintf(`if=virtio,file=%s`, imageFile),
		`-fw_cfg`, fmt.Sprintf(`name=opt/com.coreos/config,file=%s`, ign),
	)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	infof("🚀 Launching image!\n")
	cmd.Run()
}

func getImagePath(version string) string {
	if version == "latest" {
		versions, err := list()
		if err != nil {
			panic(err)
		}
		if len(versions) == 0 {
			panic(errors.New("no images available to run"))
		}
		sort.Strings(versions)
		version = versions[len(versions)-1]
	}

	dir, err := imgDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(dir, prefix+version+suffix)
}

func newRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `run [fcos config]`,
		Short: `Runs a config using qemu-kvm in an ephemeral virtual machine`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			version := cmd.Flags().Lookup("version").Value.String()
			run(args[0], version)
		},
	}
	cmd.Flags().String("version", "latest", "Image version to use as base. Use `latest` to run the most recent")

	return cmd
}
