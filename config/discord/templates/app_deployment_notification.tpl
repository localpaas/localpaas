**[{{.ProjectName}}][{{.AppName}}] Deployment {{if .Succeeded}}succeeded{{else}}failed{{end}}**
> Project: `{{.ProjectName}}`
> App: `{{.AppName}}`
> {{if .Method | eq "repo"}}
> Repository: {{.RepoURL}}
> Branch/Ref: `{{.RepoRef}}`
> Commit message: `{{.CommitMsg}}`
> {{else if .Method | eq "image"}}
> Image: `{{.Image}}`
> {{else if .Method | eq "tarball"}}
> Source archive: `{{.SourceArchive}}`
> {{end}}
> Started at: `{{.StartedAt}}`
> Duration: `{{.Duration}}`
> See deployment details: {{.DashboardLink}}
