// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.513
package view

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "alvintanoto.id/blog-htmx-templ/internal/dto"

func RegisterPage(dto *dto.RegisterPageDTO) templ.Component {
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
		templ_7745c5c3_Err = headerComponent("sign-in", nil).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<body class=\"max-h-[100vh] select-none bg-grey text-base\"><div class=\"min-h-[calc(100vh-64px)] my-4 flex items-center justify-center\"><div class=\"w-[480px] h-auto  flex-col  border border-grey-darker rounded-md p-4 mb-4\"><div class=\"text-3xl font-light mb-2\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var2 := `Register`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var2)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><form method=\"post\" action=\"/api/register\"><div class=\"flex flex-row my-3\"><div class=\"min-h-[48px] min-w-[48px] bg-grey-darker rounded-l-sm flex items-center justify-center\"><svg width=\"24px\" height=\"24px\"><image xlink:href=\"/assets/icons/user.svg\" width=\"24px\" height=\"24px\"></image></svg></div><input type=\"text\" placeholder=\"Username\" name=\"username\" class=\"h-[48px] w-full px-2 py-1 text-base rounded-r-sm outline-none\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(dto.RegisterFieldDTO.Username.Value))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, errValue := range dto.RegisterFieldDTO.Username.Errors {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"text-base text-error-text my-1\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(errValue)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `view/register_page.templ`, Line: 30, Col: 18}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex flex-row my-3\"><div class=\"min-h-[48px] min-w-[48px] bg-grey-darker rounded-l-sm flex items-center justify-center\"><svg width=\"24px\" height=\"24px\"><image xlink:href=\"/assets/icons/email.svg\" width=\"24px\" height=\"24px\"></image></svg></div><input type=\"email\" placeholder=\"Email\" name=\"email\" class=\"h-[48px] w-full px-2 py-1 text-base rounded-r-sm outline-none\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(dto.RegisterFieldDTO.Email.Value))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, errValue := range dto.RegisterFieldDTO.Email.Errors {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"text-base text-error-text my-1\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 string
			templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(errValue)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `view/register_page.templ`, Line: 49, Col: 18}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex flex-row my-3\"><div class=\"min-h-[48px] min-w-[48px] bg-grey-darker rounded-l-sm flex items-center justify-center\"><svg width=\"24px\" height=\"24px\"><image xlink:href=\"/assets/icons/password.svg\" width=\"24px\" height=\"24px\"></image></svg></div><input type=\"password\" placeholder=\"Password\" name=\"password\" class=\"h-[48px] w-full px-2 py-1 text-base rounded-r-sm outline-none\"></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, errValue := range dto.RegisterFieldDTO.PasswordErrors {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"text-base text-error-text my-1\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(errValue)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `view/register_page.templ`, Line: 67, Col: 18}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex flex-row my-3\"><div class=\"min-h-[48px] min-w-[48px] bg-grey-darker rounded-l-sm flex items-center justify-center\"><svg width=\"24px\" height=\"24px\"><image xlink:href=\"/assets/icons/password.svg\" width=\"24px\" height=\"24px\"></image></svg></div><input type=\"password\" placeholder=\"Confirm Password\" name=\"confirm_password\" class=\"h-[48px] w-full px-2 py-1 text-base rounded-r-sm outline-none\"></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, errValue := range dto.RegisterFieldDTO.ConfirmPasswordErrors {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"text-base text-error-text my-1\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var6 string
			templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(errValue)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `view/register_page.templ`, Line: 85, Col: 18}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"mt-3 float-right\"><a href=\"/sign-in\"><span class=\"mx-2 cursor-pointer hover:text-primary\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var7 := `Sign In`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var7)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></a> <button class=\"bg-primary text-white font-light px-2 py-1 rounded-sm hover:bg-primary/90\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var8 := `REGISTER`
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var8)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</button></div></form></div></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
