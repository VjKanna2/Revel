package app

import (
	"fmt"
	"net/http"

	_ "github.com/revel/modules"
	"github.com/revel/revel"
)

var (
	AppVersion string
	BuildTime  string
)

func VerifyAuthToken(c *revel.Controller) revel.Result {
	return nil
}

func init() {

	fmt.Println("App Starts.")

	revel.InterceptFunc(VerifyAuthToken, revel.BEFORE, revel.AllControllers)

	revel.Filters = []revel.Filter{
		revel.PanicFilter,
		AllowedMethodsFilter,
		revel.RouterFilter,
		revel.FilterConfiguringFilter,
		revel.ParamsFilter,
		revel.SessionFilter,
		revel.FlashFilter,
		revel.ValidationFilter,
		revel.I18nFilter,
		HeaderFilter,
		revel.InterceptorFilter,
		revel.CompressFilter,
		revel.BeforeAfterFilter,
		revel.ActionInvoker,
	}

}

var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:])
}

var AllowedMethodsFilter = func(c *revel.Controller, fc []revel.Filter) {
	switch c.Request.Method {
	case http.MethodGet, http.MethodPost, http.MethodOptions:
		fc[0](c, fc[1:])
	default:
		c.Response.Status = http.StatusMethodNotAllowed
		c.Response.Out.Header().Set("Allow", "GET, POST, OPTIONS")
		c.Result = c.RenderJSON(map[string]string{
			"status":  "error",
			"message": "Method not allowed",
		})
	}
}
