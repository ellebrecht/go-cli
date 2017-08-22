// +build windows

package util

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	log "geeny/log"
)

// inline windows batch script to fetch the Git installation path from the registry
const findGit = `@echo off
setlocal enabledelayedexpansion
for /f "skip=2 delims=: tokens=1*" %%a in ('reg query "HKLM\Software\GitForWindows" /v InstallPath 2^> nul') do (
 for /f "tokens=3" %%z in ("%%a") do (
  set GIT=%%z:%%b\
  goto FOUND
 )
)
for %%k in (HKCU HKLM) do (
 for %%w in (\ \Wow6432Node\) do (
  for /f "skip=2 delims=: tokens=1*" %%a in ('reg query "%%k\SOFTWARE%%wMicrosoft\Windows\CurrentVersion\Uninstall\Git_is1" /v InstallLocation 2^> nul') do (
   for /f "tokens=3" %%z in ("%%a") do (
    set GIT=%%z:%%b\
    goto FOUND
   )
  )
 )
)
goto EOF
:FOUND
set PATH=%GIT%bin;%PATH%
echo %GIT%
goto EOF
:EOF
`

func ExecuteCommand(name string, origArgs ...string) ([]byte, error) {
	var args = []string{"-c", name + " " + strings.Join(origArgs, " ")}
	var sh, err = gitBash()
	if err != nil {
		return make([]byte, 0), err
	}
	log.Trace("Using Windows git bash", sh)
	cmd := exec.Command(sh, args...)
	log.Trace("Command exec", *cmd)
	out, err := cmd.Output()
	log.Trace("Command output", string(out))
	if err != nil {
		log.Error("Shell command failed with error: " + err.Error())
	}
	return out, err
}

func gitBash() (string, error) {
	path, err := gitBashFromPath()
	if err == nil {
		return path, err
	}
	return gitBashFromRegistry()
}

func gitBashFromPath() (string, error) {
	out, err := exec.Command("sh.exe", "-c", "git --version").Output()
	if err != nil {
		return "", err
	}
	match, err := regexp.Match("^git[a-zA-Z0-9 ._-]+\\s*$", out)
	if err == nil && match {
		log.Trace("Found git bash in PATH")
		return "sh.exe", err
	}
	return "", errors.New("sh.exe not found in PATH")
}

func gitBashFromRegistry() (string, error) {
	file := CreateTempFile()
	fmt.Fprint(file, findGit) //TODO: should use output.Fprint...
	file.Sync()
	file.Close()
	var oldfilename = file.Name()
	var newfilename = oldfilename + ".bat"
	os.Rename(oldfilename, newfilename)
	log.Trace("Executing", newfilename)
	out, err := exec.Command("cmd", "/C", "start", "/B", newfilename).Output()
	os.Remove(newfilename)
	if err != nil || len(out) < 3 {
		return string(out), err
	}
	log.Trace("Execution finished successfully")
	regex := regexp.MustCompile("^([a-zA-Z0-9\\\\ .:=_-]+)\\s*")
	match := regex.FindSubmatch(out)
	if match == nil || len(match) < 2 {
		return "", errors.New("Problem locating git bash, is git installed?")
	}
	path := string(match[1])
	log.Trace("Found git bash in", path)
	return path + "bin\\sh.exe", err
}
