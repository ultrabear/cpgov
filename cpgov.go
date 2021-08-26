package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	FATAL = "\033[91mFATAL:\033[0m "
	// Read a maximum of 1024 bytes into memory at a time
	MaxRead = 1024
)

// Error handling function, since all errors are exiting
func handle(e error, cond bool, n ...interface{}) {
	if e != nil || cond {
		fmt.Print(FATAL)
		fmt.Println(n...)
		os.Exit(1)
	}
}

// Filters dirs to only include cpu folders
func filter(cpu []os.DirEntry) []os.DirEntry {
	ns := cpu[:0]
	for _, f := range cpu {
		if strings.HasPrefix(f.Name(), "cpu") && f.IsDir() {
			if _, e := strconv.Atoi(f.Name()[3:]); e == nil {
				ns = append(ns, f)
			}
		}
	}

	// GC rest of direntries
	for i := len(ns); i < len(cpu); i++ {
		cpu[i] = nil
	}

	return ns
}

// List of *os.File, defined to allow closing all at once
type OSFileList []*os.File

func (o OSFileList) Close() {
	for _, f := range o {
		f.Close()
	}
}

// Gets all the cpu scaling files in the given filemode with flag
// Exits program on error
func getCPUfiles(flag int) OSFileList {

	cpus, e := os.ReadDir("/sys/devices/system/cpu/")

	handle(e, false, "Could not read cpu dir")

	cpus = filter(cpus)

	files := make([]*os.File, 0, len(cpus))
	for _, f := range cpus {
		fd, e := os.OpenFile("/sys/devices/system/cpu/"+f.Name()+"/cpufreq/scaling_governor", flag, 0311)
		handle(e, false, "Could not open gov for", f.Name())
		files = append(files, fd)
	}

	return OSFileList(files)
}

// Grabs current governor that all cpus are running
// While its designed for a single governor, it supports multi governor printing anyways
// This tool is meant to be simple to use, for a more complex user, coding their own script is a good idea
func getCurrentGov() string {

	files := getCPUfiles(os.O_RDONLY)
	defer files.Close()

	govs := map[string]struct{}{}

	for _, f := range files {
		mxb := [MaxRead]byte{}
		i, e := f.Read(mxb[:])
		handle(e, false, "Could not read from governor file")
		govs[string(mxb[:i])] = struct{}{}
	}

	// If len >1 then strconcat and return it as a prettyprinted slice
	// if len == 1 then return the single string
	// if len = 0 return nil string, this shouldent be possible normally
	if len(govs) > 1 {
		s := make([]string, 0, len(govs))
		for k := range govs {
			s = append(s, strings.TrimSuffix(k, "\n"))
		}
		return fmt.Sprint(s)
	} else {
		for k := range govs {
			return strings.TrimSuffix(k, "\n")
		}
		return ""
	}
}

// Grabs valid governors from the cpu blocks
func getValidGovs() []string {

	cpus, e := os.ReadDir("/sys/devices/system/cpu/")

	handle(e, false, "Could not read cpu dir")

	cpus = filter(cpus)

	files := make([]*os.File, 0, len(cpus))
	for _, f := range cpus {
		fd, e := os.OpenFile("/sys/devices/system/cpu/"+f.Name()+"/cpufreq/scaling_available_governors", os.O_RDONLY, 0311)
		defer fd.Close()
		handle(e, false, "Could not open gov for", f.Name())
		files = append(files, fd)
	}

	// Compiled to map instead of slice to squash repeat values
	vgovs := map[string]struct{}{}

	for _, f := range files {
		rb := [MaxRead]byte{}
		n, re := f.Read(rb[:])
		handle(re, false, "Could not read from governor file")
		for _, v := range strings.Split(strings.Replace(string(rb[:n]), "\n", "", -1), " ") {
			if len(v) > 0 {
				vgovs[v] = struct{}{}
			}
		}
	}

	rstr := make([]string, 0, len(vgovs))
	for k := range vgovs {
		rstr = append(rstr, k)
	}

	return rstr
}

// Validates governor input as valid, and handles user input
func validateGovs() string {

	valid := getValidGovs()

	sort.Strings(valid)

	// Handle compiling to set dynamically
	validset := map[string]struct{}{}
	for _, v := range valid {
		validset[v] = struct{}{}
	}

	// return current gov if no input
	if len(os.Args) < 2 {
		fmt.Println("No gov specified, Current:", getCurrentGov(), "Valid:", valid)
		os.Exit(0)
	}

	// Ensure valid governor
	_, ok := validset[os.Args[1]]

	handle(nil, !ok, "Not valid gov, Valid:", valid)

	return os.Args[1]
}

// Main function, grabs governor and writes it to files
func main() {

	governor := validateGovs()

	files := getCPUfiles(os.O_WRONLY)
	defer files.Close()

	for _, f := range files {
		_, we := f.Write([]byte(governor))
		handle(we, false, "Could not write to governor file:", f.Name())
	}
}
