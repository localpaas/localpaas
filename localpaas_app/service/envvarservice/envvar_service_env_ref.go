package envvarservice

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/localpaas/localpaas/localpaas_app/entity"
)

var (
	reEnvRef = regexp.MustCompile(`\$\{(secrets\.)?(\w)+\}`)
)

func (s *envVarService) ProcessEnvVarRefs(env *EnvVar, secretStores []map[string]*entity.Secret,
	envStores []map[string]string) {
	s.processEnvVarRefs(env, secretStores, envStores, make(map[string]struct{}))
}

func (s *envVarService) processEnvVarRefs(env *EnvVar, secretStores []map[string]*entity.Secret,
	envStores []map[string]string, visitMap map[string]struct{}) {
	replFunc := func(match string) string {
		envName, isSecret := parseEnvName(match) // env form: ${NAME} or ${secrets.NAME}
		if isSecret {
			for _, store := range secretStores {
				secret, exists := store[envName]
				if !exists {
					continue
				}
				value, err := secret.Value.GetPlain()
				if err != nil {
					env.Error += fmt.Sprintf("failed to parse secret '%s'\n", envName)
					return match
				}
				return value
			}
			env.Error += fmt.Sprintf("secret '%s' not found\n", envName)
			return match
		}

		// Prevent infinite loop due to circular references
		if _, exists := visitMap[envName]; exists {
			env.Error += fmt.Sprintf("circular references detected at '%s'\n", envName)
			return match
		}
		visitMap[envName] = struct{}{}

		for _, store := range envStores {
			val, exists := store[envName]
			if !exists {
				continue
			}
			refEnv := &EnvVar{Key: envName, Value: val}
			s.processEnvVarRefs(refEnv, secretStores, envStores, visitMap)
			if refEnv.Error != "" {
				env.Error += refEnv.Error
				return match
			}
			return refEnv.Value
		}

		env.Error += fmt.Sprintf("env '%s' not found\n", envName)
		return match
	}

	env.Value = reEnvRef.ReplaceAllStringFunc(env.Value, replFunc)
}

func parseEnvName(match string) (string, bool) {
	envName := match[2 : len(match)-1]
	if strings.HasPrefix(envName, "secrets.") {
		return strings.TrimPrefix(envName, "secrets."), true
	}
	return envName, false
}
