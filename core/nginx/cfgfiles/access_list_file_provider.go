package cfgfiles

import (
	"fmt"
	"strings"

	"github.com/ncw/pwhash/apr1_crypt"

	"dillmann.com.br/nginx-ignition/core/access_list"
)

type accessListFileProvider struct {
	accessListRepository access_list.Repository
}

func newAccessListFileProvider(accessListRepository access_list.Repository) *accessListFileProvider {
	return &accessListFileProvider{accessListRepository: accessListRepository}
}

func (p *accessListFileProvider) provide(ctx *providerContext) ([]File, error) {
	accessLists, err := p.accessListRepository.FindAll(ctx.context)
	if err != nil {
		return nil, err
	}

	var outputs []File
	for _, accessList := range accessLists {
		outputs = append(outputs, p.build(accessList, ctx.paths)...)
	}

	return outputs, nil
}

func (p *accessListFileProvider) build(accessList *access_list.AccessList, paths *Paths) []File {
	var outputs []File

	if confFile := p.buildConfFile(accessList, paths); confFile != nil {
		outputs = append(outputs, *confFile)
	}

	if htpasswdFile := p.buildHtpasswdFile(accessList); htpasswdFile != nil {
		outputs = append(outputs, *htpasswdFile)
	}

	return outputs
}

func (p *accessListFileProvider) buildConfFile(accessList *access_list.AccessList, paths *Paths) *File {
	var entriesContents []string
	for _, entry := range accessList.Entries {
		for _, sourceAddress := range entry.SourceAddress {
			if sourceAddress == nil {
				continue
			}

			entriesContents = append(
				entriesContents,
				fmt.Sprintf("%s %s;", toNginxOperation(entry.Outcome), *sourceAddress),
			)
		}
	}

	usernamePasswordContents := ""
	if len(accessList.Credentials) > 0 {
		usernamePasswordContents = fmt.Sprintf(
			`
				auth_basic "%s"; 
				auth_basic_user_file %saccess-list-%s.htpasswd;
			`,
			accessList.Realm, paths.Config, accessList.ID,
		)
	}

	satisfyContents := "satisfy any;"
	if len(accessList.Credentials) > 0 && len(accessList.Entries) > 0 {
		if accessList.SatisfyAll {
			satisfyContents = "satisfy all;"
		} else {
			satisfyContents = "satisfy any;"
		}
	}

	forwardHeadersContents := ""
	if !accessList.ForwardAuthenticationHeader {
		forwardHeadersContents = `proxy_set_header Authorization "";`
	}

	contents := fmt.Sprintf(
		"%s\n%s\n%s all;\n%s\n%s",
		satisfyContents,
		strings.Join(entriesContents, "\n"),
		toNginxOperation(accessList.DefaultOutcome),
		usernamePasswordContents,
		forwardHeadersContents,
	)

	return &File{
		Name:     fmt.Sprintf("access-list-%s.conf", accessList.ID),
		Contents: contents,
	}
}

func (p *accessListFileProvider) buildHtpasswdFile(accessList *access_list.AccessList) *File {
	if len(accessList.Credentials) == 0 {
		return nil
	}

	var contents []string
	for _, credential := range accessList.Credentials {
		hash := apr1_crypt.Crypt(credential.Password, apr1_crypt.GenerateSalt(8))
		contents = append(contents, fmt.Sprintf("%s:%s", credential.Username, hash))
	}

	return &File{
		Name:     fmt.Sprintf("access-list-%s.htpasswd", accessList.ID),
		Contents: strings.Join(contents, "\n"),
	}
}

func toNginxOperation(outcome access_list.Outcome) string {
	switch outcome {
	case access_list.AllowOutcome:
		return "allow"
	case access_list.DenyOutcome:
		return "deny"
	default:
		return ""
	}
}
