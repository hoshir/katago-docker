package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var DefaultDockerImage = "hoshir/katago-v1.2-cuda10.0-linux-x64"
var MountDir = "/data"

var (
	configFile = flag.String("config", "", "A path to the config file")
	image      = flag.String("image", DefaultDockerImage, "A docker image name")
	modelFile  = flag.String("model", "", "A path to the model file")
)

func getPathInDocker(base, path string) (string, error) {
	if strings.Contains(path, "..") || strings.HasPrefix(path, "/") {
		return "", fmt.Errorf("path should be relative from the current directory")
	}

	if strings.HasPrefix(path, "./") {
		path = path[2:]
	}

	return fmt.Sprintf("%s/%s", MountDir, path), nil
}

func buildDockerCommand() string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	modelFilePath, err := getPathInDocker(currentDir, *modelFile)
	if err != nil {
		log.Fatalln("invalid model path.", err)
	}

	configFilePath, err := getPathInDocker(currentDir, *configFile)
	if err != nil {
		log.Fatalln("invalid config path.", err)
	}

	dockerArgs := fmt.Sprintf("run -it --rm --runtime nvidia -v %s:%s", currentDir, MountDir)
	katagoArg := fmt.Sprintf("gtp -model %s -config %s -override-version 0.17", modelFilePath, configFilePath)

	return fmt.Sprintf("%s %s %s", dockerArgs, *image, katagoArg)
}

func checkArgs() {
	flag.Parse()

	// Check arguments
	if *modelFile == "" {
		log.Fatalln("-model should not be empty")
	}

	if *configFile == "" {
		log.Fatalln("-config should not be empty")
	}
}

func main() {
	checkArgs()

	dockerArgs := strings.Split(buildDockerCommand(), " ")
	log.Println(dockerArgs)
	cmd := exec.Command("docker", dockerArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Start()
	if err != nil {
		log.Fatal("Failed to start docker", err)
	}

	cmd.Wait()
}
