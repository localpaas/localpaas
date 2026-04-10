package base

const (
	LocalpaasAppServiceName = "localpaas_app"
	LocalpaasAppKey         = "localpaas_app"

	LocalpaasDbServiceName = "localpaas_db"
	LocalpaasDbAppKey      = "localpaas_db"

	LocalpaasCacheServiceName = "localpaas_redis"
	LocalpaasCacheAppKey      = "localpaas_redis"

	LocalpaasTraefikServiceName = "localpaas_traefik"
	LocalpaasTraefikAppKey      = "localpaas_traefik"
)

const (
	LocalpaasProjectName = "LocalPaaS"
	LocalpaasProjectKey  = "localpaas"
)

var (
	UnallowedProjectKeys = []string{LocalpaasProjectKey}
)

const (
	NetworkGlobalRouting = "localpaas_net"
)
