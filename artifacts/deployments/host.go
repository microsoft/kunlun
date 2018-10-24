package deployments

import yaml "gopkg.in/yaml.v2"

type Host struct {
	Alias         string
	Host          string // if different from the alias you wish to give to it.
	User          string // the user to run this job.
	SSHCommonArgs string
	Extras        yaml.MapSlice
}
