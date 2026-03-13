**{{if .ProjectName | ne ""}}[{{.ProjectName}}]{{if .AppName | ne ""}}[{{.AppName}}]{{end}}{{else}}[System]{{end}} Healthcheck {{if .Succeeded}}succeeded{{else}}failed{{end}}**
{{if .ProjectName | ne ""}}> Project: `{{.ProjectName}}`{{end}}
{{if .AppName | ne ""}}> App: `{{.AppName}}`{{end}}
> Name: `{{.HealthcheckName}}`
> Type: `{{.HealthcheckType}}`
> Started at: `{{.StartedAt}}`
> Duration: `{{.Duration}}`
> Retries: `{{.Retries}}`
> {{if not .Succeeded}}Expect: `{{.Expect}}`{{end}}
> {{if not .Succeeded}}Actual: `{{.Actual}}`{{end}}
> See task details: {{.DashboardLink}}
