package main

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"
)

const ClusterName = "loadtest"
const ResourceGroup = "loadtest"

type runConf struct {
	Name              string
	CreateScriptPath  string
	DeleteScriptPath  string
	LogDir            string
	Locations         []string
	MaxConcurrency    int
	Version           string
	ClustersPerRegion int
	Delete            bool
	Create            bool
}

func makeConfig() *runConf {
	workdir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get workdir: %v", err)
	}

	flag.Usage = func() {
		w := flag.CommandLine.Output() // may be os.Stderr - but not necessarily
		fmt.Fprintf(w, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(w, "Examples: \n")
		fmt.Fprintf(w, "  - Create 6 clusters in the regions uksouth and australiawest with version 4.15.23 with 3 cluster creations simultaneuosly. No cleanup after cluster creation.:\n\n")
		fmt.Fprintf(w, "\tgo run run.go -regions='uksouth,australiawest' -version='4.15.23' -create='true' -delete='false' -num=6 -concurrency='3'\n\n")
		fmt.Fprintf(w, "  - Delete the 6 clusters in the regions uksouth and australiawest created by the previous run:\n\n")
		fmt.Fprintf(w, "\tgo run run.go -regions='uksouth,australiawest' -version='4.15.23' -create='false' -delete='true' -num=6 -concurrency='3'\n\n")
	}
	regionsInput := flag.String("regions", "", "comman separated list of regions in which to run the tests")
	clusterNum := flag.Int("num", 6, "number of clusters created/deleted per region")
	version := flag.String("version", "", "version to test")
	maxConcurrency := flag.Int("concurrency", 0, "Maximum number of parallel cluster creations/deletions. <= 0 means unlimited parallelity")
	delete := flag.Bool("delete", false, "Weather or not clusters should be deleted.")
	create := flag.Bool("create", true, "weather or not clusters should be created.")
	flag.Parse()

	if *regionsInput == "" {
		log.Fatalln("-regions can't be empty")
	}
	if *version == "" {
		log.Fatalln("-version can't be empty")
	}

	locations := strings.Split(*regionsInput, ",")

	return &runConf{
		Name:              "loadtest",
		CreateScriptPath:  path.Join(workdir, "create.sh"),
		DeleteScriptPath:  path.Join(workdir, "delete.sh"),
		LogDir:            path.Join(workdir, "logs"),
		Locations:         locations,
		MaxConcurrency:    *maxConcurrency,
		Version:           *version,
		ClustersPerRegion: *clusterNum,
		Delete:            *delete,
		Create:            *create,
	}
}

func main() {
	conf := makeConfig()

	err := os.Mkdir(conf.LogDir, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatalf("Unable to create log folder: %v", err)
	}

	clusterCreateRunners := []*CmdRunner{}
	clusterDeleteRunners := []*CmdRunner{}

	for _, loc := range conf.Locations {
		for i := 0; i < conf.ClustersPerRegion; i++ {
			if conf.Create {
				createRunner := NewRunner(conf.CreateScriptPath, loc, conf.Version, i)
				clusterCreateRunners = append(clusterCreateRunners, createRunner)
			}
			if conf.Delete {
				deleteRunner := NewRunner(conf.DeleteScriptPath, loc, conf.Version, i)
				clusterDeleteRunners = append(clusterDeleteRunners, deleteRunner)
			}
		}
	}

	log.Printf("Creating %d clusters...", len(clusterCreateRunners))
	ExecuteRunners(clusterCreateRunners, conf.MaxConcurrency)

	erroredRunners := []*CmdRunner{}
	for _, runner := range clusterCreateRunners {
		runnerName := fmt.Sprintf("%s-%d", runner.Location, runner.Num)
		err = os.WriteFile(path.Join(conf.LogDir, runnerName), []byte(runner.Output), 0644)
		if err != nil {
			log.Printf("Error writing logs: %v", err)
		}

		if runner.Err != nil {
			erroredRunners = append(erroredRunners, runner)
		}
	}

	log.Printf("%d / %d Create calls encountered an error.", len(erroredRunners), len(clusterCreateRunners))
	log.Println("Errored Runners:")
	for _, runner := range erroredRunners {
		log.Printf("\t- %s-%d", runner.Location, runner.Num)
	}

	log.Printf("deleting %d clusters", len(clusterDeleteRunners))
	ExecuteRunners(clusterDeleteRunners, conf.MaxConcurrency)
}

func ExecuteRunners(runners []*CmdRunner, maxConcurrency int) {
	wg := &sync.WaitGroup{}
	wg.Add(len(runners))
	if maxConcurrency <= 0 {
		maxConcurrency = len(runners)
	}
	concurrencyLimiter := make(chan int, maxConcurrency)
	for i := 0; i < maxConcurrency; i++ {
		concurrencyLimiter <- 0
	}

	// status routine to print progress
	statusDone := make(chan int)
	go func() {
		defer func() { statusDone <- 0 }()
		for {
			numDone := 0
			total := len(runners)
			for _, runner := range runners {
				if runner.Cmd.ProcessState != nil {
					numDone++
				}
			}
			log.Printf("%d / %d runners are done.\n", numDone, total)
			if numDone >= total {
				return
			}
			time.Sleep(5 * time.Second)
		}
	}()

	for _, runner := range runners {
		<-concurrencyLimiter
		go func(runner *CmdRunner) {
			defer func() { concurrencyLimiter <- 0 }()
			defer wg.Done()
			runner.Run()
		}(runner)
	}
	wg.Wait()
	<-statusDone
}

type CmdRunner struct {
	Cmd      *exec.Cmd
	Output   string
	Err      error
	Location string
	Num      int
}

func NewRunner(script string, location string, version string, num int) *CmdRunner {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "bash"
	}

	newCmd := exec.Command(shell, "-i", "-c", script)
	newCmd.Env = []string{
		"VERSION=" + version,
		"RESOURCEGROUP=" + fmt.Sprintf("%s-%s-%d", ResourceGroup, location, num),
		"CLUSTER=" + fmt.Sprintf("%s-%s-%d", ClusterName, location, num),
		"LOCATION=" + location,
	}

	return &CmdRunner{
		Cmd:      newCmd,
		Location: location,
		Num:      num,
	}
}

func (c *CmdRunner) Run() {
	out, err := c.Cmd.CombinedOutput()
	c.Output = string(out)
	c.Err = err
}
