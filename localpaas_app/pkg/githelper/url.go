package githelper

import (
	"fmt"

	"github.com/gitsight/go-vcsurl"
)

// git@gitlab.com:commento/docs.git
// git@github.com:go-git/go-git.git
// git clone git@bitbucket.org:mcuadros/discovery-rest.git
func GetSshUrl(v *vcsurl.VCS) string {
	return fmt.Sprintf("git@%s:%s/%s.git", v.Host, v.Username, v.Name)
}

// https://mcuadros@bitbucket.org/mcuadros/discovery-rest.git
// https://gitlab.com/commento/docs.git
// https://github.com/go-git/go-git.git
func GetHttpsUrl(v *vcsurl.VCS) string {
	return fmt.Sprintf("https://%s/%s/%s.git", v.Host, v.Username, v.Name)
}
