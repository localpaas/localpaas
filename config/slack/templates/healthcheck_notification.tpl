{{if .Succeeded}}
*[{{.ProjectName}}/{{.AppName}}] Healthcheck succeeded*
{{else}}
*[{{.ProjectName}}/{{.AppName}}] Healthcheck failed*
{{end}}
> {{if .ProjectName | ne ""}}Project: `{{.ProjectName}}`{{end}}
> {{if .AppName | ne ""}}App: `{{.AppName}}`{{end}}
> Healthcheck name: `{{.HealthcheckName}}`
> Healthcheck type: `{{.HealthcheckType}}`
> Started at: `{{.StartedAt}}`
> Duration: `{{.Duration}}`
> Retries: `{{.Retries}}`
> See task details: {{.DashboardLink}}
