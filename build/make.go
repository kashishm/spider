package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/getgauge/spider/version"
)

const (
	CGO_ENABLED = "CGO_ENABLED"
)

const (
	distros = "distros"
	GOARCH  = "GOARCH"
	GOOS    = "GOOS"
	X86     = "386"
	X86_64  = "amd64"
	DARWIN  = "darwin"
	LINUX   = "linux"
	WINDOWS = "windows"
	bin     = "bin"
	spider  = "spider"
)

func main() {
	compileAcrossPlatforms()
	createPluginDistro()
}

func createPluginDistro() {
	for _, platformEnv := range platformEnvs {
		setEnv(platformEnv)
		fmt.Printf("Creating distro for platform => OS:%s ARCH:%s \n", platformEnv[GOOS], platformEnv[GOARCH])
		createDistro()
	}
	log.Printf("Distributables created in directory => %s \n", filepath.Join(bin, distros))
}

func createDistro() {
	packageName := fmt.Sprintf("%s-%s-%s.%s", spider, getPluginVersion(), getGOOS(), getArch())
	os.Mkdir(filepath.Join(bin, distros), 0755)
	createZipFromUtil(getBinDir(), packageName)
}

func createZipFromUtil(dir, name string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	os.Chdir(dir)
	output, err := executeCommand("zip", "-r", filepath.Join("..", distros, name+".zip"), ".")
	fmt.Println(output)
	if err != nil {
		panic(fmt.Sprintf("Failed to zip: %s", err))
	}
	os.Chdir(wd)
}

func runProcess(command string, arg ...string) {
	cmd := exec.Command(command, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("Execute %v\n", cmd.Args)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func executeCommand(command string, arg ...string) (string, error) {
	cmd := exec.Command(command, arg...)
	bytes, err := cmd.Output()
	return strings.TrimSpace(fmt.Sprintf("%s", bytes)), err
}

func compileGoPackage() {
	runProcess("go", "get", "-d", "./...")
	runProcess("go", "build", "-o", getGaugeExecutablePath(spider))
}

func getGaugeExecutablePath(file string) string {
	return filepath.Join(getBinDir(), getExecutableName(file))
}

func getExecutableName(file string) string {
	if getGOOS() == "windows" {
		return file + ".exe"
	}
	return file
}

func getBinDir() string {
	return filepath.Join(bin, fmt.Sprintf("%s_%s", getGOOS(), getGOARCH()))
}

func getPluginVersion() string {
	return version.Version
}

func setEnv(envVariables map[string]string) {
	for k, v := range envVariables {
		os.Setenv(k, v)
	}
}

var (
	platformEnvs = []map[string]string{
		map[string]string{GOARCH: X86, GOOS: DARWIN, CGO_ENABLED: "0"},
		map[string]string{GOARCH: X86_64, GOOS: DARWIN, CGO_ENABLED: "0"},
		map[string]string{GOARCH: X86, GOOS: LINUX, CGO_ENABLED: "0"},
		map[string]string{GOARCH: X86_64, GOOS: LINUX, CGO_ENABLED: "0"},
		map[string]string{GOARCH: X86, GOOS: WINDOWS, CGO_ENABLED: "0"},
		map[string]string{GOARCH: X86_64, GOOS: WINDOWS, CGO_ENABLED: "0"},
	}
)

func compileAcrossPlatforms() {
	for _, platformEnv := range platformEnvs {
		setEnv(platformEnv)
		fmt.Printf("Compiling for platform => OS:%s ARCH:%s \n", platformEnv[GOOS], platformEnv[GOARCH])
		compileGoPackage()
	}
}

func getArch() string {
	if getGOARCH() == X86 {
		return "x86"
	}
	return "x86_64"
}

func getGOARCH() string {
	goArch := os.Getenv(GOARCH)
	if goArch == "" {
		return runtime.GOARCH
	}
	return goArch
}

func getGOOS() string {
	goOS := os.Getenv(GOOS)
	if goOS == "" {
		return runtime.GOOS
	}
	return goOS
}
