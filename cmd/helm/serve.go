/*
Copyright 2016 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"k8s.io/helm/cmd/helm/helmpath"
	"k8s.io/helm/pkg/repo"
)

const serveDesc = `
This command starts a local chart repository server that serves charts from a local directory.

The new server will provide HTTP access to a repository. By default, it will
scan all of the charts in '$HELM_HOME/repository/local' and serve those over
the a local IPv4 TCP port (default '127.0.0.1:8879').
`

type serveCmd struct {
	out      io.Writer
	address  string
	repoPath string
}

func newServeCmd(out io.Writer) *cobra.Command {
	srv := &serveCmd{out: out}
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "start a local http web server",
		Long:  serveDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return srv.run()
		},
	}

	f := cmd.Flags()
	f.StringVar(&srv.repoPath, "repo-path", helmpath.Home(homePath()).LocalRepository(), "local directory path from which to serve charts")
	f.StringVar(&srv.address, "address", "127.0.0.1:8879", "address to listen on")

	return cmd
}

func (s *serveCmd) run() error {
	repoPath, err := filepath.Abs(s.repoPath)
	if err != nil {
		return err
	}
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return err
	}

	fmt.Fprintln(s.out, "Regenerating index. This may take a moment.")
	if err := index(repoPath, "http://"+s.address, ""); err != nil {
		return err
	}

	fmt.Fprintf(s.out, "Now serving you on %s\n", s.address)
	return repo.StartLocalRepo(repoPath, s.address)
}
