{{if .Succeeded}}
**[{{.ProjectName}}/{{.AppName}}] Healthcheck succeeded**
{{else}}
**[{{.ProjectName}}/{{.AppName}}] Healthcheck failed**
{{end}}
> {{if .ProjectName | ne ""}}Project: `{{.ProjectName}}`{{end}}
> {{if .AppName | ne ""}}App: `{{.AppName}}`{{end}}
> Name: `{{.HealthcheckName}}`
> Type: `{{.HealthcheckType}}`
> Started at: `{{.StartedAt}}`
> Duration: `{{.Duration}}`
> Retries: `{{.Retries}}`
> {{if not .Succeeded}}Expect: `{{.Expect}}`{{end}}
> {{if not .Succeeded}}Actual: `{{.Actual}}`{{end}}
> See task details: {{.DashboardLink}}
