package provider

import (
	"fmt"
	"strings"

	"github.com/ncw/pwhash/apr1_crypt"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/accesslist"
)

type accessListProvider struct {
	commands accesslist.Commands
}

func newAccessListProvider(commands accesslist.Commands) *accessListProvider {
	return &accessListProvider{
		commands: commands,
	}
}

func (p *accessListProvider) Provide(ctx *Context) ([]File, error) {
	accessLists, err := p.commands.GetAll(ctx.Context)
	if err != nil {
		return nil, err
	}

	outputs := make([]File, 0)
	for _, accessList := range accessLists {
		outputs = append(outputs, p.build(&accessList, ctx.Paths)...)
	}

	return outputs, nil
}

func (p *accessListProvider) build(accessList *accesslist.AccessList, paths *Paths) []File {
	outputs := make([]File, 0)

	if confFile := p.buildConfFile(accessList, paths); confFile != nil {
		outputs = append(outputs, *confFile)
	}

	if htpasswdFile := p.buildHtpasswdFile(accessList); htpasswdFile != nil {
		outputs = append(outputs, *htpasswdFile)
	}

	return outputs
}

func (p *accessListProvider) buildConfFile(
	accessList *accesslist.AccessList,
	paths *Paths,
) *File {
	entriesContents := make([]string, 0)
	for _, entry := range accessList.Entries {
		for _, sourceAddress := range entry.SourceAddress {
			entriesContents = append(
				entriesContents,
				nginxSprintf(
					"%s %s;",
					directiveFragment(toNginxOperation(entry.Outcome)),
					sourceAddress,
				),
			)
		}
	}

	usernamePasswordContents := ""
	if len(accessList.Credentials) > 0 {
		usernamePasswordContents = nginxSprintf(
			`
				auth_basic %s;
				auth_basic_user_file %s;
			`,
			accessList.Realm,
			paths.Config+"access-list-"+accessList.ID.String()+".htpasswd",
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

	contents := nginxSprintf(
		"%s\n%s\n%s all;\n%s\n%s",
		directiveFragment(satisfyContents),
		directiveFragment(strings.Join(entriesContents, "\n")),
		directiveFragment(toNginxOperation(accessList.DefaultOutcome)),
		directiveFragment(usernamePasswordContents),
		directiveFragment(forwardHeadersContents),
	)

	return &File{
		Name:     fmt.Sprintf("access-list-%s.conf", accessList.ID),
		Contents: contents,
	}
}

func (p *accessListProvider) buildHtpasswdFile(accessList *accesslist.AccessList) *File {
	if len(accessList.Credentials) == 0 {
		return nil
	}

	contents := make([]string, 0)
	for _, credential := range accessList.Credentials {
		hash := apr1_crypt.Crypt(credential.Password, apr1_crypt.GenerateSalt(8))
		contents = append(
			contents,
			fmt.Sprintf("%s:%s", sanitizeHtpasswdUsername(credential.Username), hash),
		)
	}

	return &File{
		Name:     fmt.Sprintf("access-list-%s.htpasswd", accessList.ID),
		Contents: strings.Join(contents, "\n"),
	}
}

func toNginxOperation(outcome accesslist.Outcome) string {
	switch outcome {
	case accesslist.AllowOutcome:
		return "allow"
	case accesslist.DenyOutcome:
		return "deny"
	default:
		return ""
	}
}
