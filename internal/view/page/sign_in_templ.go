// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package vpages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"alvintanoto.id/blog-htmx-templ/internal/dto"
	"alvintanoto.id/blog-htmx-templ/internal/view/component"
)

func SignIn(dto *dto.SignInPageDTO) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = vcomponent.Title().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<body class=\"max-h-[100vh] select-none bg-layout-background text-text text-base overflow-hidden \">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = vcomponent.Header("sign_in", nil).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"h-[calc(100vh-64px)] flex items-center justify-center overflow-y-auto\"><div class=\"w-[480px] h-auto flex-col rounded-md p-4\"><div class=\"text-3xl font-light my-4\">Sign In</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if dto.Error != "" {
			templ_7745c5c3_Err = vcomponent.AlertError(dto.Error).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form method=\"post\"><div class=\"flex flex-row my-3\"><div class=\"min-h-[48px] min-w-[48px] bg-disabled rounded-l-sm flex items-center justify-center\"><svg width=\"24px\" height=\"24px\"><image xlink:href=\"/assets/icons/user.svg\" width=\"24px\" height=\"24px\"></image></svg></div><input type=\"text\" placeholder=\"Username\" name=\"username\" class=\"h-[48px] w-full px-2 py-1 text-base rounded-r-sm outline-none\"></div><div class=\"flex flex-row my-3\"><div class=\"min-h-[48px] min-w-[48px] bg-disabled rounded-l-sm flex items-center justify-center\"><svg width=\"24px\" height=\"24px\"><image xlink:href=\"/assets/icons/password.svg\" width=\"24px\" height=\"24px\"></image></svg></div><input type=\"password\" placeholder=\"Password\" name=\"password\" class=\"h-[48px] w-full px-2 py-1 text-base rounded-r-sm outline-none\"></div><div class=\"mt-3 float-right\"><a href=\"/auth/register\"><span class=\"mx-1 cursor-pointer px-2 py-1 cursor-pointer rounded-md hover:bg-separator hover:text-primary\">Register</span></a> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 = []any{"text-center text-white rounded-md px-2 py-1 cursor-pointer shadow-sm bg-primary hover:bg-primary/80"}
		templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var2...)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ.CSSClasses(templ_7745c5c3_Var2).String()))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">SIGN IN</button></div></form></div></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
