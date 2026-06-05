<b>{{if .ProjectName}}[{{.ProjectName}}]{{if .AppName}}[{{.AppName}}]{{end}}{{else}}[System]{{end}} Scheduled task {{if .Succeeded}}✅ Succeeded{{else}}❌ Failed{{end}}</b>

{{if .ProjectName}}<b>• Project:</b> {{.ProjectName}}
{{end}}{{if .AppName}}<b>• App:</b> {{.AppName}}
{{end}}<b>• Scheduled Job:</b> <code>{{.SchedJobName}}</code>
<b>• Schedule:</b> <code>{{.Schedule}}</code>
<b>• Retries:</b> <code>{{.Retries}}</code>
<b>• Created At:</b> <code>{{.CreatedAt}}</code>
<b>• Started At:</b> <code>{{.StartedAt}}</code>
<b>• Duration:</b> <code>{{.Duration}}</code>

🔗 <a href="{{.DashboardLink}}">Go to Dashboard</a>
