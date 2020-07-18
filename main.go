// Copyright 2020 Simen A. W. Olsen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/karrick/godirwalk"
	"github.com/urfave/cli/v2"
)

type remoteInfo struct {
	Org  string
	Name string
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	app := &cli.App{
		Name:   "organize-repo-folder",
		Usage:  "moves every repo into each github org",
		Action: OrganizeRepoFolder,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Value: homeDir + "/Repos",
				Usage: "path to your repositories",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Value: false,
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getRepo() {

}

func runDirectory() {

}

func getRepoName(r *git.Remote) (*remoteInfo, error) {
	var re = regexp.MustCompile(`(?m)^(https|git)(:\/\/|@)([^\/:]+)[\/:]([^\/:]+)\/(.+).git`)
	c := r.Config()
	u := c.URLs[0]
	f := re.FindStringSubmatch(u)

	if f == nil {
		return nil, errors.New("Could not deconstruct git name (nil)")
	}

	// if len(f) >= 5 {
	// 	return nil, errors.New("Could not deconstruct git name (len was not enough)")
	// }

	if f[4] == "" || f[5] == "" {
		return nil, errors.New("Could not deconstruct git name")
	}

	return &remoteInfo{
		Org:  strings.ToLower(f[4]),
		Name: strings.ToLower(f[5]),
	}, nil
}

// OrganizeRepoFolder runs through your repositories and cleans it!
func OrganizeRepoFolder(c *cli.Context) error {
	repoPath, err := filepath.Abs(c.String("path"))
	if err != nil {
		return err
	}

	dryRun := c.Bool("dry-run")

	if dryRun {
		fmt.Println("#### Running in dry-run mode ####")
	}

	paths := make(map[string]bool)
	wo := make(map[string]string)

	err = godirwalk.Walk(repoPath, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsDir() && de.Name() == ".git" {

				r, err := git.PlainOpen(osPathname)
				if err != nil {
					fmt.Printf("[error] Was not able to read git on %s\n", osPathname)
					return err
				}
				remote, err := r.Remote("origin")

				// If repo doesn't have remote, we'll skip.
				if err != nil {
					fmt.Printf("[error] Missing remote(s) on %s\n", osPathname)
					return nil
				}

				re, err := getRepoName(remote)

				// If repo failed to get repo name, skip it.
				if err != nil {
					fmt.Printf("[error] Not able to get repository name on %s\n", osPathname)
					return nil
				}

				correctPath := repoPath + "/" + re.Org + "/" + re.Name

				if correctPath+"/.git" == osPathname {
					paths[correctPath] = true
					wo[strings.ReplaceAll(osPathname, "/.git", "")] = correctPath
					return nil
				}

				isAvail := false
				i := 0
				for isAvail == false {
					if !paths[correctPath] {
						paths[correctPath] = true
						isAvail = true
					} else {
						i = i + 1
						correctPath = correctPath + "-" + strconv.Itoa(i)
						fmt.Printf("[duplicate] %s will be named %s\n", osPathname, correctPath)
					}
				}

				wo[strings.ReplaceAll(osPathname, "/.git", "")] = correctPath
			}
			// fmt.Printf("%s %s\n", de.ModeType(), osPathname)
			return nil
		},
		Unsorted: true,
	})

	fmt.Printf("mv %s %s\nmkdir %s\n", repoPath, repoPath+"-bak", repoPath)
	if !dryRun {
		os.Rename(repoPath, repoPath+"-bak")
		_ = os.Mkdir(repoPath, 0700)
	}

	// Loop through the work order
	for s, d := range wo {
		fmt.Printf("mkdir %s\n", path.Dir(d))
		if !dryRun {
			_ = os.Mkdir(path.Dir(d), 0700)
		}
		ns := strings.ReplaceAll(s, repoPath, repoPath+"-bak")
		fmt.Printf("mv %s %s\n", ns, d)
		if !dryRun {
			os.Rename(ns, d)
		}
	}

	fmt.Printf(
		"\n\n### DONE ###\nRepositories had issues, are kept in %s.\n"+
			"You should clean that up (at least rename the folder before you run this command again).", repoPath+"-bak")

	if err != nil {
		return err
	}
	return nil

}
