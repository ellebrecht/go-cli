package util

import (
	"os"

	model "geeny/api/model"
)

type MockGitHub struct {
	IsRepoValue bool
	Value       string
	Error       error
	File        *os.File
}

func (g MockGitHub) InitRepo(path string) error {
	return g.Error
}
func (g MockGitHub) IsRepo(path string) bool {
	return g.IsRepoValue
}
func (g MockGitHub) Pull(branch string) error {
	return g.Error
}
func (g MockGitHub) Add(file string, path string) error {
	return g.Error
}
func (g MockGitHub) Commit(message string, path string) error {
	return g.Error
}
func (g MockGitHub) RenameRemote(from string, to string, path string) error {
	return g.Error
}
func (g MockGitHub) GetConfigValue(path string, name string) (*string, error) {
	return &g.Value, g.Error
}
func (g MockGitHub) SetConfigValue(path string, name string, val string) error {
	return g.Error
}
func (g MockGitHub) GetGeenyConfigValue(path string, name string) (*string, error) {
	return &g.Value, g.Error
}
func (g MockGitHub) SetGeenyConfigValue(path string, name string, val string) error {
	return g.Error
}
func (g MockGitHub) CloneRepo(r *model.Repository) error {
	return g.Error
}
func (g MockGitHub) PushRepo(r *model.Repository) error {
	return g.Error
}
func (g MockGitHub) UnSetConfigValue(path string, name string) error {
	return g.Error
}
func (g MockGitHub) SetUpRepoSSH(r *model.Repository) (*os.File, error) {
	return g.File, g.Error
}
func (g MockGitHub) TearDownRepoSSH(r *model.Repository, keyFile *os.File) error {
	return g.Error
}
