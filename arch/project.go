package arch

import (
	"github.com/pacur/pacur/constants"
	"github.com/pacur/pacur/utils"
	"path/filepath"
)

type ArchProject struct {
	Root       string
	MirrorRoot string
	BuildRoot  string
	Path       string
	Distro     string
	Release    string
}

func (p *ArchProject) getBuildDir() (path string, err error) {
	path = filepath.Join(p.BuildRoot, p.Distro)

	err = utils.MkdirAll(path)
	if err != nil {
		return
	}

	return
}

func (p *ArchProject) Prep() (err error) {
	buildDir, err := p.getBuildDir()
	if err != nil {
		return
	}

	err = utils.CopyFile("", filepath.Join(p.Path, "..", "sign.key"),
		buildDir, true)
	if err != nil {
		return
	}

	keyPath := filepath.Join(p.Path, "sign.key")
	exists, err := utils.Exists(keyPath)
	if err != nil {
		return
	}

	if exists {
		err = utils.Rsync(keyPath, buildDir)
		if err != nil {
			return
		}
	}

	return
}

func (p *ArchProject) Create() (err error) {
	buildDir, err := p.getBuildDir()
	if err != nil {
		return
	}

	err = utils.Exec("", "docker", "run", "--rm", "-t", "-v",
		buildDir+":/pacur", constants.DockerOrg+p.Distro, "create",
		p.Distro)
	if err != nil {
		return
	}

	err = utils.Rsync(filepath.Join(buildDir, "arch"),
		filepath.Join(p.MirrorRoot, "arch"))
	if err != nil {
		return
	}

	return
}
