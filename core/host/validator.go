package host

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/accesslist"
	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/constants"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"dillmann.com.br/nginx-ignition/core/common/valuerange"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

type validator struct {
	hostRepository      Repository
	integrationCommands integration.Commands
	vpnCommands         vpn.Commands
	accessListCommands  accesslist.Commands
	cacheCommands       cache.Commands
	bindingCommands     binding.Commands
	delegate            *validation.ConsistencyValidator
}

func newValidator(
	hostRepository Repository,
	integrationCommands integration.Commands,
	vpnCommands vpn.Commands,
	accessListCommands accesslist.Commands,
	cacheCommands cache.Commands,
	bindingCommands binding.Commands,
) *validator {
	return &validator{
		hostRepository:      hostRepository,
		integrationCommands: integrationCommands,
		vpnCommands:         vpnCommands,
		accessListCommands:  accessListCommands,
		cacheCommands:       cacheCommands,
		bindingCommands:     bindingCommands,
		delegate:            validation.NewValidator(),
	}
}

const (
	bindingsPath = "bindings"
)

var (
	redirectStatusCodeRange = valuerange.New(300, 399)
	statusCodeRange         = valuerange.New(100, 599)
)

func (v *validator) validate(ctx context.Context, host *Host) error {
	if err := v.validateDefaultFlag(ctx, host); err != nil {
		return err
	}

	if err := v.validateBindings(ctx, host); err != nil {
		return err
	}

	v.validateDomainNames(ctx, host)
	if err := v.validateRoutes(ctx, host); err != nil {
		return err
	}

	if err := v.validateVPNs(ctx, host); err != nil {
		return err
	}

	if err := v.validateAccessList(ctx, host.AccessListID, "accessListId"); err != nil {
		return err
	}

	if err := v.validateCache(ctx, host.CacheID, "cacheId"); err != nil {
		return err
	}

	return v.delegate.Result()
}

func (v *validator) validateDefaultFlag(ctx context.Context, host *Host) error {
	if !host.DefaultServer {
		return nil
	}

	defaultServer, err := v.hostRepository.FindDefault(ctx)
	if err != nil {
		return err
	}

	if defaultServer != nil && host.ID != defaultServer.ID {
		v.delegate.Add("defaultServer", i18n.M(ctx, i18n.K.HostValidationDefaultAlreadyExists))
	}

	if len(host.DomainNames) > 0 {
		v.delegate.Add(
			"domainNames",
			i18n.M(ctx, i18n.K.HostValidationDomainMustBeEmptyForDefault),
		)
	}

	return nil
}

func (v *validator) validateDomainNames(ctx context.Context, host *Host) {
	if len(host.DomainNames) == 0 && !host.DefaultServer {
		v.delegate.Add("domainNames", i18n.M(ctx, i18n.K.CommonValidationAtLeastOneRequired))
	}

	for index, domainName := range host.DomainNames {
		if !constants.TLDPattern.MatchString(domainName) {
			v.delegate.Add(
				fmt.Sprintf("domainNames[%d]", index),
				i18n.M(ctx, i18n.K.CommonValidationInvalidDomainName),
			)
		}
	}
}

func (v *validator) validateBindings(ctx context.Context, host *Host) error {
	if host.UseGlobalBindings && len(host.Bindings) > 0 {
		v.delegate.Add(
			bindingsPath,
			i18n.M(ctx, i18n.K.HostValidationBindingsMustBeEmptyForGlobal),
		)
	}

	if !host.UseGlobalBindings {
		if len(host.Bindings) == 0 {
			v.delegate.Add(bindingsPath, i18n.M(ctx, i18n.K.CommonValidationAtLeastOneRequired))
		}

		for index, b := range host.Bindings {
			if err := v.bindingCommands.Validate(
				ctx,
				bindingsPath,
				index,
				&b,
				v.delegate,
			); err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *validator) validateRoutes(ctx context.Context, host *Host) error {
	if len(host.Routes) == 0 {
		v.delegate.Add("routes", i18n.M(ctx, i18n.K.CommonValidationAtLeastOneRequired))
	}

	priorityMap := make(map[int]int)
	for _, route := range host.Routes {
		priorityMap[route.Priority]++
	}

	for priority, count := range priorityMap {
		if count > 1 {
			v.delegate.Add(
				"routes",
				i18n.M(ctx, i18n.K.HostValidationDuplicatedRoutePriority).V("priority", priority),
			)
		}
	}

	distinctPaths := make(map[string]bool)
	for index, route := range host.Routes {
		if err := v.validateRoute(ctx, &route, index, &distinctPaths); err != nil {
			return err
		}
	}

	return nil
}

func (v *validator) validateRoute(
	ctx context.Context,
	route *Route,
	index int,
	distinctPaths *map[string]bool,
) error {
	if (*distinctPaths)[route.SourcePath] {
		v.delegate.Add(
			buildIndexedRoutePath(index, "sourcePath"),
			i18n.M(ctx, i18n.K.HostValidationDuplicatedSourcePath),
		)
	} else {
		(*distinctPaths)[route.SourcePath] = true
	}

	if err := v.validateAccessList(
		ctx,
		route.AccessListID,
		buildIndexedRoutePath(index, "accessListId"),
	); err != nil {
		return err
	}

	if err := v.validateCache(
		ctx,
		route.CacheID,
		buildIndexedRoutePath(index, "cacheId"),
	); err != nil {
		return err
	}

	switch route.Type {
	case ProxyRouteType:
		v.validateProxyRoute(ctx, route, index)
	case RedirectRouteType:
		v.validateRedirectRoute(ctx, route, index)
	case StaticResponseRouteType:
		v.validateStaticResponseRoute(ctx, route, index)
	case IntegrationRouteType:
		return v.validateIntegrationRoute(ctx, route, index)
	case ExecuteCodeRouteType:
		v.validateExecuteCodeRoute(ctx, route, index)
	case StaticFilesRouteType:
		v.validateStaticFilesRoute(ctx, route, index)
	default:
		v.delegate.Add(
			buildIndexedRoutePath(index, "type"),
			i18n.M(ctx, i18n.K.CommonValidationInvalidValue),
		)
	}

	return nil
}

func (v *validator) validateStaticFilesRoute(ctx context.Context, route *Route, index int) {
	targetURIField := buildIndexedRoutePath(index, "targetUri")
	if route.TargetURI == nil || strings.TrimSpace(*route.TargetURI) == "" {
		v.delegate.Add(
			targetURIField,
			i18n.M(ctx, i18n.K.HostValidationTargetUriRequired).V("type", "directory"),
		)
		return
	}

	if !strings.HasPrefix(*route.TargetURI, "/") {
		v.delegate.Add(targetURIField, i18n.M(ctx, i18n.K.CommonValidationStartsWithSlashRequired))
	}
}

func (v *validator) validateProxyRoute(ctx context.Context, route *Route, index int) {
	targetURIField := buildIndexedRoutePath(index, "targetUri")
	if route.TargetURI == nil || strings.TrimSpace(*route.TargetURI) == "" {
		v.delegate.Add(
			targetURIField,
			i18n.M(ctx, i18n.K.HostValidationTargetUriRequired).V("type", "proxy"),
		)
	} else {
		if _, err := url.Parse(*route.TargetURI); err != nil {
			v.delegate.Add(targetURIField, i18n.M(ctx, i18n.K.CommonValidationInvalidUrl))
		}
	}
}

func (v *validator) validateRedirectRoute(ctx context.Context, route *Route, index int) {
	targetURIField := buildIndexedRoutePath(index, "targetUri")
	if route.TargetURI == nil || strings.TrimSpace(*route.TargetURI) == "" {
		v.delegate.Add(
			targetURIField,
			i18n.M(ctx, i18n.K.HostValidationTargetUriRequired).V("type", "redirect"),
		)
	} else {
		if _, err := url.ParseRequestURI(*route.TargetURI); err != nil {
			v.delegate.Add(targetURIField, i18n.M(ctx, i18n.K.CommonValidationInvalidUri))
		}
	}

	if route.RedirectCode == nil || !redirectStatusCodeRange.Contains(*route.RedirectCode) {
		v.delegate.Add(
			buildIndexedRoutePath(index, "redirectCode"),
			i18n.M(ctx, i18n.K.CommonValidationBetweenValues).
				V("min", redirectStatusCodeRange.Min).
				V("max", redirectStatusCodeRange.Max),
		)
	}
}

func (v *validator) validateStaticResponseRoute(ctx context.Context, route *Route, index int) {
	if route.Response == nil {
		v.delegate.Add(
			buildIndexedRoutePath(index, "response"),
			i18n.M(ctx, i18n.K.HostValidationStaticResponseRequired),
		)
		return
	}

	if !statusCodeRange.Contains(route.Response.StatusCode) {
		v.delegate.Add(
			buildIndexedRoutePath(index, "response.statusCode"),
			i18n.M(ctx, i18n.K.CommonValidationBetweenValues).
				V("min", statusCodeRange.Min).
				V("max", statusCodeRange.Max),
		)
	}
}

func (v *validator) validateIntegrationRoute(ctx context.Context, route *Route, index int) error {
	requiredMessage := i18n.M(ctx, i18n.K.HostValidationIntegrationRequired)

	if route.Integration == nil {
		v.delegate.Add(buildIndexedRoutePath(index, "integration"), requiredMessage)
		return nil
	}

	exists, err := v.integrationCommands.Exists(ctx, route.Integration.IntegrationID)
	if err != nil {
		return err
	}

	if !*exists {
		v.delegate.Add(buildIndexedRoutePath(index, "integration.integrationId"), requiredMessage)
	}

	if strings.TrimSpace(route.Integration.OptionID) == "" {
		v.delegate.Add(buildIndexedRoutePath(index, "integration.optionId"), requiredMessage)
	}

	return nil
}

func (v *validator) validateExecuteCodeRoute(ctx context.Context, route *Route, index int) {
	requiredMessage := i18n.M(ctx, i18n.K.HostValidationSourceCodeRequired)

	if route.SourceCode == nil {
		v.delegate.Add(buildIndexedRoutePath(index, "sourceCode"), requiredMessage)
		return
	}

	if strings.TrimSpace(route.SourceCode.Contents) == "" {
		v.delegate.Add(buildIndexedRoutePath(index, "sourceCode.code"), requiredMessage)
	}

	if route.SourceCode.Language != JavascriptCodeLanguage &&
		route.SourceCode.Language != LuaCodeLanguage {
		v.delegate.Add(
			buildIndexedRoutePath(index, "sourceCode.language"),
			i18n.M(ctx, i18n.K.CommonValidationInvalidValue),
		)
	}

	if route.SourceCode.Language == JavascriptCodeLanguage &&
		(route.SourceCode.MainFunction == nil || strings.TrimSpace(*route.SourceCode.MainFunction) == "") {
		v.delegate.Add(
			buildIndexedRoutePath(index, "sourceCode.mainFunction"),
			i18n.M(ctx, i18n.K.HostValidationJsMainFunctionRequired),
		)
	}
}

func (v *validator) validateVPNs(ctx context.Context, host *Host) error {
	vpnNameUsage := make(map[uuid.UUID]map[string]int)

	for index, value := range host.VPNs {
		basePath := fmt.Sprintf("vpns[%d]", index)
		vpnIDPath := basePath + ".vpnId"
		namePath := basePath + ".name"

		if strings.TrimSpace(value.Name) == "" {
			v.delegate.Add(namePath, i18n.M(ctx, i18n.K.CommonValidationValueMissing))
		}

		if value.VPNID == uuid.Nil {
			v.delegate.Add(vpnIDPath, i18n.M(ctx, i18n.K.CommonValidationValueMissing))
			continue
		}

		if vpnNameUsage[value.VPNID] == nil {
			vpnNameUsage[value.VPNID] = make(map[string]int)
		}

		if vpnNameUsage[value.VPNID][value.Name] > 0 {
			v.delegate.Add(namePath, i18n.M(ctx, i18n.K.HostValidationDuplicatedVpnName))
		}

		vpnNameUsage[value.VPNID][value.Name]++

		vpnData, err := v.vpnCommands.Get(ctx, value.VPNID)
		if err != nil {
			return err
		}

		if vpnData == nil {
			v.delegate.Add(vpnIDPath, i18n.M(ctx, i18n.K.HostValidationVpnNotFound))
			continue
		}

		if !vpnData.Enabled {
			v.delegate.Add(vpnIDPath, i18n.M(ctx, i18n.K.HostValidationVpnDisabled))
		}
	}

	return nil
}

func buildIndexedRoutePath(index int, childPath string) string {
	return fmt.Sprintf("routes[%d].%s", index, childPath)
}

func (v *validator) validateAccessList(
	ctx context.Context,
	accessListID *uuid.UUID,
	path string,
) error {
	if accessListID == nil {
		return nil
	}

	exists, err := v.accessListCommands.Exists(ctx, *accessListID)
	if err != nil {
		return err
	}

	if !exists {
		v.delegate.Add(path, i18n.M(ctx, i18n.K.HostValidationAccessListNotFound))
	}

	return nil
}

func (v *validator) validateCache(ctx context.Context, cacheID *uuid.UUID, path string) error {
	if cacheID == nil {
		return nil
	}

	exists, err := v.cacheCommands.Exists(ctx, *cacheID)
	if err != nil {
		return err
	}

	if !exists {
		v.delegate.Add(path, i18n.M(ctx, i18n.K.HostValidationCacheNotFound))
	}

	return nil
}
