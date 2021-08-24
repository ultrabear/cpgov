package main

import (
	"os"
	"fmt"
	"strconv"
	"strings"
	"sort"
)

const (
	FATAL = "\033[91mFATAL:\033[0m "
)

func handle(e error, cond bool, n ...interface{}) {
	if e != nil || cond {
		fmt.Print(FATAL)
		fmt.Println(n...)
		os.Exit(1)
	}
}

func filter(cpu []os.DirEntry) []os.DirEntry {
	ns := cpu[:0]
	for _, f := range cpu {
		if strings.HasPrefix(f.Name(), "cpu") && f.IsDir() {
			if _, e := strconv.Atoi(f.Name()[3:]); e == nil {
				ns = append(ns, f)
			}
		}
	}

	// GC hint
	for i := len(ns); i < len(cpu); i++ {
		cpu[i] = nil
	}

	return ns
}

type OSFileList []*os.File

func (o OSFileList) Close() {
	for _, f := range o {
		f.Close()
	}
}

func getCPUfiles(flag int) OSFileList {

	cpus, e := os.ReadDir("/sys/devices/system/cpu/")

	handle(e, false, "Could not read cpu dir")

	cpus = filter(cpus)

	files := make([]*os.File, 0, len(cpus))
	for _, f := range cpus {
		fd, e := os.OpenFile("/sys/devices/system/cpu/" + f.Name() + "/cpufreq/scaling_governor", flag, 0311)
		handle(e, false, "Could not open gov for " + f.Name())
		files = append(files, fd)
	}

	return OSFileList(files)
}

func getCurrentGov() string {

	files := getCPUfiles(os.O_RDONLY)
	defer files.Close()

	govs := map[string]struct{}{}

	for _, f := range files {
		mxb := make([]byte, 1024)
		i, e := f.Read(mxb)
		handle(e, false, "Could not read from governor file")
		govs[string(mxb[:i])] = struct{}{}
	}

	if len(govs) > 1 {
		s := make([]string, 0, len(govs))
		for k, _ := range govs {
			s = append(s, strings.TrimSuffix(k, "\n"))
		}
		return fmt.Sprint(s)
	} else {
		for k, _ := range govs {
			return strings.TrimSuffix(k, "\n")
		}
		return ""
	}
}

func getValidGovs() []string {

	cpus, e := os.ReadDir("/sys/devices/system/cpu/")

	handle(e, false, "Could not read cpu dir")

	cpus = filter(cpus)

	files := make([]*os.File, 0, len(cpus))
	for _, f := range cpus {
		fd, e := os.OpenFile("/sys/devices/system/cpu/" + f.Name() + "/cpufreq/scaling_available_governors", os.O_RDONLY, 0311)
		defer fd.Close()
		handle(e, false, "Could not open gov for " + f.Name())
		files = append(files, fd)
	}

	vgovs := map[string]struct{}{}

	for _, f := range files {
		rb := make([]byte, 1024)
		n, re := f.Read(rb)
		handle(re, false, "Could not read from governor file")
		for _, v := range strings.Split(strings.Replace(string(rb[:n]), "\n", "", -1), " ") {
			if len(v) > 0 {
				vgovs[v] = struct{}{}
			}
		}
	}

	rstr := make([]string, 0, len(vgovs))
	for k, _ := range vgovs {
		rstr = append(rstr, k)
	}

	return rstr
}

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

	_, ok := validset[os.Args[1]]

	handle(nil, !ok, "Not valid gov, Valid: ", valid)

	return os.Args[1]
}

func main() {

	governor := validateGovs()

	files := getCPUfiles(os.O_WRONLY)
	defer files.Close()

	for _, f := range files {
		_, we := f.Write([]byte(governor))
		handle(we, false, "Could not write to governor file")
	}
}
