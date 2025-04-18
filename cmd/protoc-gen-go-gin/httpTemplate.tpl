{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}

{{- range .MethodSets}}
const Operation{{$svrType}}{{.OriginalName}} = "/{{$svrName}}/{{.OriginalName}}"
{{- end}}

type I{{.ServiceType}}Server interface {
{{- range .MethodSets}}
	{{- if ne .Comment ""}}
	{{.Comment}}
	{{- end}}
	{{.Name}}(*gin.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- end}}
}

func Register{{.ServiceType}}Server(r gin.IRoutes, srv I{{.ServiceType}}Server) {
	{{- range .Methods}}
	r.{{.Method}}("{{.Path}}", _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv))
	{{- end}}
}

{{range .Methods}}
func _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv I{{$svrType}}Server) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var in {{.Request}}
		{{- if .HasBody}}
		if err := ctx.ShouldBindJSON(&in{{.Body}}); err != nil {
		    resp.FailWithMessage(ctx, err.Error())
			return
		}
		{{- end}}
		if err := ctx.BindQuery(&in); err != nil {
		    resp.FailWithMessage(ctx, err.Error())
			return
		}
		{{- if .HasVars}}
		if err := ctx.ShouldBindUri(&in); err != nil {
			resp.FailWithMessage(ctx, err.Error())
            return
		}
		{{- end}}
		// http.SetOperation(ctx, Operation{{$svrType}}{{.OriginalName}})
		reply, err := srv.{{.Name}}(ctx, &in)
		if err != nil {
		    resp.FailWithMessage(ctx, err.Error())
			return
		}
		resp.OkWithData(ctx, reply{{.ResponseBody}})
	}
}
{{end}}
