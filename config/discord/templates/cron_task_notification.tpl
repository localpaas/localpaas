{{if .Succeeded}}
**[{{.ProjectName}}/{{.AppName}}] Scheduled task succeeded**
{{else}}
**[{{.ProjectName}}/{{.AppName}}] Scheduled task failed**
{{end}}
> {{if .ProjectName | ne ""}}Project: `{{.ProjectName}}`{{end}}
> {{if .AppName | ne ""}}App: `{{.AppName}}`{{end}}
> Cron job: `{{.CronJobName}}`
> Cron job expr: `{{.CronJobExpr}}`
> Created at: `{{.CreatedAt}}`
> Started at: `{{.StartedAt}}`
> Duration: `{{.Duration}}`
> Retries: `{{.Retries}}`
> See task details: {{.DashboardLink}}
