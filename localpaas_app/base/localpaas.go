package base

const (
	LocalpaasAppServiceName = "localpaas_app"
	LocalpaasAppKey         = "localpaas_app"

	LocalpaasWorkerServiceName = "localpaas_worker"
	LocalpaasWorkerKey         = "localpaas_worker"

	LocalpaasDbServiceName = "localpaas_db"
	LocalpaasDbAppKey      = "localpaas_db"

	LocalpaasCacheServiceName = "localpaas_redis"
	LocalpaasCacheAppKey      = "localpaas_redis"

	LocalpaasTraefikServiceName = "localpaas_traefik"
	LocalpaasTraefikAppKey      = "localpaas_traefik"

	LocalpaasUpdaterServiceName = "localpaas_updater"
	LocalpaasUpdaterKey         = "localpaas_updater"

	LocalpaasDockerProxyServiceName = "localpaas_docker_proxy"
	LocalpaasDockerProxyKey         = "localpaas_docker_proxy"

	LocalpaasAgentServiceName = "localpaas_agent"
	LocalpaasAgentKey         = "localpaas_agent"
)

const (
	LocalpaasProjectName = "LocalPaaS"
	LocalpaasProjectKey  = "localpaas"
)

var (
	UnallowedProjectKeys = []string{LocalpaasProjectKey}
)

const (
	NetworkGlobalRouting  = "localpaas_net"
	NetworkDockerProxy    = "localpaas_docker_proxy_net"
	NetworkLocalpaasLocal = "localpaas_local_net"
)
