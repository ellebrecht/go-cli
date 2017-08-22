package data

import (
	"errors"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	model "geeny/api/model"
	log "geeny/log"
	util "geeny/util"
)

const (
	gitConfig = ".git/config"
	geeny     = "geeny"
)

// GitHub deals with git commands, e.g. clone, push. Currently wrapped around the `git` command
type GitHub struct {
}

// GitHubInterface describes the GitHub functions
type GitHubInterface interface {
	InitRepo(path string) error
	IsRepo(path string) bool
	Pull(branch string) error
	Add(file string, path string) error
	Commit(message string, path string) error
	RenameRemote(from string, to string, path string) error
	GetConfigValue(path string, name string) (*string, error)
	SetConfigValue(path string, name string, val string) error
	GetGeenyConfigValue(path string, name string) (*string, error)
	SetGeenyConfigValue(path string, name string, val string) error
	UnSetConfigValue(path string, name string) error
	SetUpRepoSSH(r *model.Repository) (*os.File, error)
	TearDownRepoSSH(r *model.Repository, keyFile *os.File) error
	CloneRepo(r *model.Repository) error
	PushRepo(r *model.Repository) error
}

// SSH return an ssh repo url
func SSH(company string, name string) string {
	return "git@github.com:" + company + "/" + name + ".git"
}

// HTTPS return an https repo url
func HTTPS(company string, name string) string {
	return "https://github.com/" + company + "/" + name + ".git"
}

// InitRepo initialises an empty repo
func (g *GitHub) InitRepo(path string) error {
	log.Tracef("trying initialise git repo at path %s", path)
	_, err := util.ExecuteCommand("git", "--git-dir", path, "init")
	if err != nil {
		return errors.New("there was a problem initialising an empty git repo at " + path + ", with error: " + err.Error())
	}
	return nil
}

// IsRepo returns true if it's a repo
func (g *GitHub) IsRepo(path string) bool {
	log.Tracef("checking if repo %s", path)
	_, err := os.Stat(path)
	if err != nil {
		log.Error("got err", err)
	}
	return err == nil
}

// GetConfigValue gets a value from the .git/config file
func (g *GitHub) GetConfigValue(path string, name string) (*string, error) {
	log.Tracef("trying to get config value %s at path %s", name, path)
	out, err := util.ExecuteCommand("git", "config", "-f", path+"/"+gitConfig, "--get", name)
	if err != nil {
		return nil, errors.New("there was a problem getting a value from .git/config. are you in the repo?: " + err.Error())
	}
	str := strings.Replace(string(out), "\n", "", -1) // for some reason there's a newline...
	return &str, nil
}

// SetConfigValue sets a value in the .git/config file
func (g *GitHub) SetConfigValue(path string, name string, val string) error {
	log.Tracef("trying to set config value %s at path %s", name, path)
	var err error
	if runtime.GOOS == "windows" {
		_, err = util.ExecuteCommand("git", "config", "-f", path+"/"+gitConfig, "--replace-all", name, "\""+val+"\"")
	} else {
		_, err = util.ExecuteCommand("git", "config", "-f", path+"/"+gitConfig, "--replace-all", name, val)
	}
	if err != nil {
		return errors.New("there was a problem setting a value in .git/config. are you in the repo?: " + err.Error())
	}
	return nil
}

// UnSetConfigValue unsets a value from the .git/config file
func (g *GitHub) UnSetConfigValue(path string, name string) error {
	log.Tracef("trying to unset config value %s at path %s", name, path)
	var err error
	if runtime.GOOS == "windows" {
		_, err = util.ExecuteCommand("git", "config", "-f", path+"/"+gitConfig, "--unset-all", "\""+name+"\"")
	} else {
		_, err = util.ExecuteCommand("git", "config", "-f", path+"/"+gitConfig, "--unset-all", name)
	}
	if err != nil {
		return errors.New("there was a problem unsetting a value from .git/config. are you in the repo?: " + err.Error())
	}
	return nil
}

// GetGeenyConfigValue gets a geeny value from the .git/config file
func (g *GitHub) GetGeenyConfigValue(path string, name string) (*string, error) {
	return g.GetConfigValue(path, geeny+"."+name)
}

// SetGeenyConfigValue sets a geeny value in the .git/config file
func (g *GitHub) SetGeenyConfigValue(path string, name string, val string) error {
	return g.SetConfigValue(path, geeny+"."+name, val)
}

// SetUpRepoSSH prepares a repo for ssh
func (g *GitHub) SetUpRepoSSH(r *model.Repository) (*os.File, error) {
	log.Tracef("preparing repo %s for ssh", r.Name)
	tmpFile := util.CreateTempFile()
	err := g.prepareSSH(r, tmpFile)
	if err != nil {
		return nil, err
	}
	filePathUnix := strings.Replace(tmpFile.Name(), "\\", "/", -1)
	log.Tracef("check tempFile.name vs. replace slashes: %s %s", tmpFile.Name(), filePathUnix)
	err = g.SetConfigValue(".", "core.sshCommand", "ssh -i "+filePathUnix+" -o IdentitiesOnly=yes")
	if err != nil {
		return nil, err
	}
	return tmpFile, nil
}

// TearDownRepoSSH deconfigures the ssh repo
func (g *GitHub) TearDownRepoSSH(r *model.Repository, keyFile *os.File) error {
	log.Tracef("removing core.sshCommand from %s .git/config, and removing file %s", r.Name, keyFile.Name())
	err := g.UnSetConfigValue(".", "core.sshCommand")
	if err != nil {
		return err
	}
	return os.Remove(keyFile.Name())
}

// CloneRepo clones a github repo
func (g *GitHub) CloneRepo(r *model.Repository) error {
	if r.IsSSH() {
		return g.cloneRepoSSH(r)
	}
	if r.IsHTTPS() {
		return g.cloneRepoHTTPS(r)
	}
	return errors.New("repo " + r.Name + " has bad url " + r.URL)
}

// Pull updates from project
func (g *GitHub) Pull(branch string) error {
	log.Tracef("pulling origin %s", branch)
	_, err := util.ExecuteCommand("git", "pull", "geeny", branch, "--allow-unrelated-histories")
	if err != nil {
		return errors.New("there was a problem pulling files. are you in the repo?: " + err.Error())
	}
	return nil
}

// Add file to project
func (g *GitHub) Add(file string, path string) error {
	log.Tracef("adding file %s", file)
	_, err := util.ExecuteCommand("git", "--git-dir="+path+"/.git/", "--work-tree="+path, "add", file)
	if err != nil {
		return errors.New("there was a problem adding files. are you in the repo?: " + err.Error())
	}
	return nil
}

// Commit with a message
func (g *GitHub) Commit(message string, path string) error {
	if runtime.GOOS == "windows" {
		return g.commit("--git-dir="+path+"/.git/", "--work-tree="+path, "commit", "-m", "\""+message+"\"")
	}
	return g.commit("--git-dir="+path+"/.git/", "--work-tree="+path, "commit", "-m", message)
}

// RenameRemote renames a remote
func (g *GitHub) RenameRemote(from string, to string, path string) error {
	log.Tracef("renaming remote %s to %s", from, to)
	_, err := util.ExecuteCommand("git", "--git-dir="+path+"/.git/", "--work-tree="+path, "remote", "rename", from, to)
	if err != nil {
		return errors.New("there was a problem renaming the remote. are you in the repo?: " + err.Error())
	}
	return nil
}

// PushRepo pushes the github repo using ssh credentials
func (g *GitHub) PushRepo(r *model.Repository) error {
	log.Tracef("trying to push %s", r.URL)
	_, err := util.ExecuteCommand("git", "push", "geeny", "master", "-f")
	if err != nil {
		return errors.New("there was a problem pushing your repo. are you in the repo? " + err.Error())
	}
	return nil
}

// - private

func (g *GitHub) commit(args ...string) error {
	log.Tracef("try to commit: %s", args)
	_, err := util.ExecuteCommand("git", args...)
	if err != nil {
		return errors.New("there was a problem committing - maybe there are no changes?: " + err.Error())
	}
	return nil
}

func (g *GitHub) prepareSSH(r *model.Repository, tmpFile *os.File) error {
	log.Tracef("saving ssh key %s temporarily to: %s", r.PrivateKey, tmpFile.Name())
	return ioutil.WriteFile(tmpFile.Name(), []byte(r.PrivateKey), 0600)
}

func (g *GitHub) cloneRepoSSH(r *model.Repository) error {
	log.Tracef("trying to ssh clone %s", r.Name)
	tmpFile := util.CreateTempFile()
	err := g.prepareSSH(r, tmpFile)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	_, err = util.ExecuteCommand("ssh-agent", "bash", "-c", "ssh-add "+tmpFile.Name()+"; GIT_SSH_COMMAND=\"ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no\" git clone "+r.URL)
	if err != nil {
		return errors.New("there was a problem cloning over ssh. are you in the repo? and do you have your github ssh keys correctly setup?: " + err.Error())
	}
	return nil
}

func (g *GitHub) cloneRepoHTTPS(r *model.Repository) error {
	log.Tracef("trying to http clone %s %s", r.Name, r.URL)
	_, err := util.ExecuteCommand("git", "clone", "--depth", "1", r.URL)
	if err != nil {
		return errors.New("there was an problem cloning over http: " + err.Error())
	}
	return nil
}
