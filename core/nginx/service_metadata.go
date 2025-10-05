package nginx

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/core_error"
	"dillmann.com.br/nginx-ignition/core/common/log"
)

func (s *service) getMetadata(ctx context.Context) (*Metadata, error) {
	cmd := exec.CommandContext(ctx, s.processManager.binaryPath, "-V")
	rawOutput, err := cmd.CombinedOutput()
	if err != nil {
		return nil, core_error.New("failed to execute nginx -V: "+err.Error(), false)
	}

	output := string(rawOutput)
	version := s.extractVersion(output)
	buildDetails := s.extractBuildDetails(output)
	tlsSniEnabled := s.extractTlsSniEnabled(output)
	configureArgs := s.extractConfigureArguments(output)
	staticModules := s.extractStaticModules(configureArgs)
	dynamicModules := s.extractDynamicModules(configureArgs)
	modulesPath := s.extractModulesPath(configureArgs)
	fileModules := s.listModuleFiles(modulesPath)
	modules := s.mergeModules(staticModules, dynamicModules, fileModules)

	return &Metadata{
		Version:       version,
		BuildDetails:  buildDetails,
		TLSSNIEnabled: tlsSniEnabled,
		Modules:       modules,
	}, nil
}

func (s *service) extractVersion(output string) string {
	re := regexp.MustCompile(`nginx version: nginx/(.+)`)
	matches := re.FindStringSubmatch(output)

	if len(matches) > 1 {
		return matches[1]
	}

	return "unknown"
}

func (s *service) extractBuildDetails(output string) string {
	regex := regexp.MustCompile(`(?m)^built (.+)$`)
	matches := regex.FindAllStringSubmatch(output, -1)
	buildLines := make([]string, 0, len(matches))

	for _, match := range matches {
		if len(match) > 1 {
			buildLines = append(buildLines, match[1])
		}
	}

	return strings.Join(buildLines, "; ")
}

func (s *service) extractTlsSniEnabled(output string) bool {
	return strings.Contains(output, "TLS SNI support enabled")
}

func (s *service) extractConfigureArguments(output string) string {
	regex := regexp.MustCompile(`configure arguments: (.+)`)
	matches := regex.FindStringSubmatch(output)

	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

func (s *service) extractStaticModules(configureArgs string) []string {
	modules := make([]string, 0)
	modulesSet := make(map[string]bool)
	moduleRegex := regexp.MustCompile(`--with-([a-zA-Z0-9_]+_module)(?:\s|$)`)
	moduleMatches := moduleRegex.FindAllStringSubmatch(configureArgs, -1)

	for _, match := range moduleMatches {
		if len(match) > 1 {
			moduleName := match[1]

			if !modulesSet[moduleName] {
				modulesSet[moduleName] = true
				modules = append(modules, moduleName)
			}
		}
	}

	simpleRegex := regexp.MustCompile(`--with-([a-zA-Z0-9_-]+)(?:\s|$)`)
	simpleMatches := simpleRegex.FindAllStringSubmatch(configureArgs, -1)

	for _, match := range simpleMatches {
		if len(match) > 1 {
			moduleName := match[1]

			if !strings.HasSuffix(moduleName, "_module") && !modulesSet[moduleName] {
				modulesSet[moduleName] = true
				modules = append(modules, moduleName)
			}
		}
	}

	return modules
}

func (s *service) extractDynamicModules(configureArgs string) []string {
	modules := make([]string, 0)
	modulesSet := make(map[string]bool)

	suffixRegex := regexp.MustCompile(`--with-([a-zA-Z0-9_]+_module)=dynamic`)
	suffixMatches := suffixRegex.FindAllStringSubmatch(configureArgs, -1)

	for _, match := range suffixMatches {
		if len(match) > 1 {
			moduleName := match[1]

			if !modulesSet[moduleName] {
				modulesSet[moduleName] = true
				modules = append(modules, moduleName)
			}
		}
	}

	prefixRegex := regexp.MustCompile(`--add-dynamic-module=([^\s]+)`)
	prefixMatches := prefixRegex.FindAllStringSubmatch(configureArgs, -1)

	for _, match := range prefixMatches {
		if len(match) > 1 {
			path := match[1]
			moduleName := filepath.Base(path)
			moduleName = strings.TrimSuffix(moduleName, "/")

			if !modulesSet[moduleName] {
				modulesSet[moduleName] = true
				modules = append(modules, moduleName)
			}
		}
	}

	return modules
}

func (s *service) extractModulesPath(configureArgs string) string {
	regex := regexp.MustCompile(`--modules-path=([^\s]+)`)
	matches := regex.FindStringSubmatch(configureArgs)

	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

func (s *service) listModuleFiles(modulesPath string) []string {
	if modulesPath == "" {
		return []string{}
	}

	files := make([]string, 0)
	entries, err := os.ReadDir(modulesPath)

	if err != nil {
		log.Warnf("Failed to read modules directory %s: %v", modulesPath, err)
		return files
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".so") {
			moduleName := strings.TrimSuffix(entry.Name(), ".so")
			files = append(files, moduleName)
		}
	}

	return files
}

func (s *service) mergeModules(staticModules, dynamicModules, fileModules []string) []string {
	modulesSet := make(map[string]bool)
	modules := make([]string, 0)

	for _, module := range staticModules {
		if !modulesSet[module] {
			modulesSet[module] = true
			modules = append(modules, module)
		}
	}

	for _, module := range dynamicModules {
		if !modulesSet[module] {
			modulesSet[module] = true
			modules = append(modules, module)
		}
	}

	for _, module := range fileModules {
		if !modulesSet[module] {
			modulesSet[module] = true
			modules = append(modules, module)
		}
	}

	return modules
}
