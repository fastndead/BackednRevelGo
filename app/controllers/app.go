package controllers

import "github.com/revel/revel"
import (
	_ "github.com/lib/pq"
	"app/app/lib/auth"
	"app/app/lib/responce"
)


type App struct {
	*revel.Controller
}

func (c *App) Index() revel.Result {
	return c.Render()
}


func (c *App)checkAuthentifcation() revel.Result{
	retVal, err := auth.Auth(c.Controller)
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return retVal
}

func init() {
	revel.InterceptMethod((*App).checkAuthentifcation, revel.AFTER)
	revel.InterceptMethod((*CFlight).Init, revel.BEFORE)
	revel.InterceptMethod((*CFlight).DbClose, revel.AFTER)
	revel.InterceptMethod((*CPlane).Init, revel.BEFORE)
	revel.InterceptMethod((*CPlane).DbClose, revel.AFTER)
	revel.InterceptMethod((*CPilot).Init, revel.BEFORE)
	revel.InterceptMethod((*CPilot).DbClose, revel.AFTER)
}