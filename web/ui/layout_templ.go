// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package ui

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Layout(title string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<html lang=\"en\"><head><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1, shrink-to-fit=no\"><meta name=\"theme-color\" content=\"#000000\"><script>\n\t\t\tfunction copyToClipboard() {\n\t\t\t\tvar copyText = document.getElementById(\"sealedSecretYaml\");\n\t\t\t\tcopyText.select();\n\t\t\t\tcopyText.setSelectionRange(0, 99999);\n\t\t\t\tdocument.execCommand(\"copy\");\n\t\t\t}\n\t\t</script><style>\n\t\t\t\t.token.number,\n\t\t\t\t.token.tag {\n\t\t\t\t  all: inherit;\n\t\t\t\t  color: hsl(14, 58%, 55%);\n\t\t\t\t}\n\t\t\t\t.loading-indicator {\n        \tdisplay:none;\n    \t\t}\n    \t\t.htmx-request .loading-indicator {\n        \tdisplay:inline;\n    \t\t}\n    \t\t.htmx-request.loading-indicator {\n        \tdisplay:inline;\n    \t\t}\n\t\t\t</style><link rel=\"stylesheet\" href=\"https://cdn.jsdelivr.net/npm/bulma@0.9.4/css/bulma.min.css\"><title>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(title)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/ui/layout.templ`, Line: 34, Col: 17}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</title><script src=\"https://unpkg.com/htmx.org@2.0.1\" integrity=\"sha384-QWGpdj554B4ETpJJC9z+ZHJcA/i59TyjxEPXiiUgN2WmTyV5OEZWCD6gQhgkdpB/\" crossorigin=\"anonymous\"></script></head><body><div id=\"content\" class=\"container p-5 content\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ_7745c5c3_Var1.Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><script>\n\t\t\tdocument.addEventListener(\"htmx:beforeRequest\", function(event) {\n\t\t\t\tconsole.log(event);\n\t\t\t\t\tif (event.detail.pathInfo.requestPath === \"/sealed-secret\") {\n\t\t\t\t\t\tdocument.getElementById(\"encryptButton\").style.display = \"none\";\n\t\t\t\t\t\tconst element = document.querySelector(\".message\");\n\t\t\t\t\t\tif (element) {\n    \t\t\t\t\telement.style.display = \"none\";\n\t\t\t\t\t\t}\n\t\t\t\t\t};\n\t\t\t});\n\t\t\tdocument.addEventListener(\"htmx:afterRequest\", function(event) {\n\t\t\t\tconsole.log(event);\n\t\t\t\t\tif (event.detail.pathInfo.requestPath === \"/sealed-secret\") {\n\t\t\t\t\t\tdocument.getElementById(\"encryptButton\").style.display = \"block\";\n\t\t\t\t\t}\n\t\t\t});\n\t\t\t</script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
