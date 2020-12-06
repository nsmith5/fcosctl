package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	prefix = `fedora-coreos-`
	suffix = `-qemu.x86_64.qcow2`
)

func imgDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, `/.local/share/libvirt/images/`), nil
}

func list() (versions []string, err error) {

	dir, err := imgDir()
	if err != nil {
		return
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, suffix) {
			name = strings.TrimPrefix(name, prefix)
			name = strings.TrimSuffix(name, suffix)
			versions = append(versions, name)
		}
	}
	return
}

func pull(stream string) {
	errf := color.New(color.FgRed).PrintfFunc()
	infof := color.New(color.FgCyan).PrintfFunc()

	_, ok := map[string]bool{
		"stable":  true,
		"testing": true,
		"next":    true,
	}[stream]
	if !ok {
		errf("Stream must be one of stable, testing or next. Received %s\n", stream)
		os.Exit(1)
	}

	dir, err := imgDir()
	if err != nil {
		errf("Failed to find homedir %s\n", err)
		os.Exit(1)
	}

	infof("Pulling latest image from stream %s\n", stream)
	cmd := exec.Command(
		`coreos-installer`,
		`download`,
		`-s`, stream,
		`-p`, `qemu`,
		`-f`, `qcow2.xz`,
		`--decompress`,
		`-C`, dir,
	)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func deleteImg(version string) {
	errf := color.New(color.FgRed).PrintfFunc()
	infof := color.New(color.FgCyan).PrintfFunc()

	dir, err := imgDir()
	if err != nil {
		errf("Failed to get image dir: %s\n", err)
		os.Exit(1)
	}

	file := prefix + version + suffix
	path := filepath.Join(dir, file)
	err = os.Remove(path)
	if err != nil {
		errf("Failed to delete image: %s\n", err)
		os.Exit(1)
	}
	infof("✅ Version %s deleted\n", version)
}

func newImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   `image`,
		Short: `Manage FCOS images`,
	}

	list := &cobra.Command{
		Use:   `list`,
		Short: `List available images`,
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			errf := color.New(color.FgRed).PrintfFunc()
			infof := color.New(color.FgCyan).PrintfFunc()

			versions, err := list()
			if err != nil {
				errf("Failed to list images; %s\n", err)
				os.Exit(1)
			}

			for _, version := range versions {
				infof("* %s\n", version)
			}
		},
	}

	pull := &cobra.Command{
		Use:   `pull`,
		Short: `Download FCOS images`,
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			stream := cmd.Flags().Lookup("stream").Value.String()
			pull(stream)
		},
	}
	pull.Flags().String("stream", "stable", "Stream to pull from. Must be one of stable, testing or next")

	deleteCmd := &cobra.Command{
		Use:   `delete`,
		Short: `Delete FCOS images`,
		Run: func(cmd *cobra.Command, args []string) {
			version := cmd.Flags().Lookup("version").Value.String()
			deleteImg(version)
		},
	}
	deleteCmd.Flags().String("version", "", "Version to delete")
	cobra.MarkFlagRequired(deleteCmd.Flags(), "version")

	cmd.AddCommand(list)
	cmd.AddCommand(pull)
	cmd.AddCommand(deleteCmd)

	return cmd
}
