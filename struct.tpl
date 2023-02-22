package eventos

type {{ Title .Name}} struct {
    {{range .Fields}}
        {{ Title . }}
    {{end}}
}