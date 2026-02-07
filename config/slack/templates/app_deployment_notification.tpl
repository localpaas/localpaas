{{if .Succeeded}}
*[{{.ProjectName}}/{{.AppName}}] Deployment succeeded*
{{else}}
*[{{.ProjectName}}/{{.AppName}}] Deployment failed*
{{end}}
> Project: `{{.ProjectName}}`
> App: `{{.AppName}}`
> {{if .Method | eq "repo"}}
> Repository: `{{.RepoURL}}`
> Branch/Ref: `{{.RepoRef}}`
> Commit Message: `{{.CommitMsg}}`
> {{else if .Method | eq "image"}}
> Image: `{{.Image}}`
> {{else if .Method | eq "tarball"}}
> Source Archive: `{{.SourceArchive}}`
> {{end}}
> Started at: `{{.StartedAt}}`
> Duration: `{{.Duration}}`
> See deployment details: {{.DashboardLink}}
