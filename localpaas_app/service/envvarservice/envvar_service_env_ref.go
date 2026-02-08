package envvarservice

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/localpaas/localpaas/localpaas_app/entity"
)

var (
	reEnvOrSecretRef = regexp.MustCompile(`\$\{(secrets\.)?(\w)+\}`)
	// reSecretRef = regexp.MustCompile(`\$\{secrets\.(\w)+\}`)
)

func (s *envVarService) processRefs(
	env *EnvVar,
	envStore map[string]*entity.EnvVar,
	secretStore map[string]*entity.Secret,
) {
	processRefs(env, envStore, secretStore, make(map[string]struct{}))
}

func processRefs(
	env *EnvVar,
	envStore map[string]*entity.EnvVar,
	secretStore map[string]*entity.Secret,
	visitMap map[string]struct{},
) {
	if env.IsLiteral {
		return
	}

	replFunc := func(match string) string {
		envName, isSecret := parseEnvName(match) // env form: ${NAME} or ${secrets.NAME}
		if isSecret {
			refSecret, exists := secretStore[envName]
			if !exists {
				env.Errors = append(env.Errors, fmt.Sprintf("secret '%s' not found", envName))
				return match
			}
			value, err := refSecret.Value.GetPlain()
			if err != nil {
				env.Errors = append(env.Errors, fmt.Sprintf("failed to parse secret '%s'", envName))
				return match
			}
			return value
		}

		// Prevent infinite loop due to circular references
		if _, exists := visitMap[envName]; exists {
			env.Errors = append(env.Errors, fmt.Sprintf("circular references detected at '%s'", envName))
			return match
		}
		visitMap[envName] = struct{}{}

		val, exists := envStore[envName]
		if !exists {
			env.Errors = append(env.Errors, fmt.Sprintf("env '%s' not found", envName))
			return match
		}
		refEnv := &EnvVar{EnvVar: val}
		processRefs(refEnv, envStore, secretStore, visitMap)
		if len(refEnv.Errors) > 0 {
			env.Errors = append(env.Errors, refEnv.Errors...)
			return match
		}
		return refEnv.Value
	}

	env.Value = reEnvOrSecretRef.ReplaceAllStringFunc(env.Value, replFunc)
}

func parseEnvName(match string) (string, bool) {
	envName := match[2 : len(match)-1]
	if strings.HasPrefix(envName, "secrets.") {
		return strings.TrimPrefix(envName, "secrets."), true
	}
	return envName, false
}
