package controllers

import "github.com/revel/revel"
import (
	_ "github.com/lib/pq"
	"app/app/lib/auth"
	"fmt"
)


type App struct {
	*revel.Controller
}

func (c *App) Index() revel.Result {

	return c.Render()
}


func (c *App)checkAuthentifcation() revel.Result{
	fmt.Println("hell")
	return auth.Auth(c.Controller)
}

func (c *App)LogOut()revel.Result{
	auth.LogOut(c.Controller)
	fmt.Println("HERE!")
	return c.Index()
}


func init() {
	revel.InterceptMethod((*App).checkAuthentifcation, revel.BEFORE)
	revel.InterceptMethod((*CFlight).Init, revel.BEFORE)
	revel.InterceptMethod((*CFlight).DbClose, revel.AFTER)
	revel.InterceptMethod((*CPlane).Init, revel.BEFORE)
	revel.InterceptMethod((*CPlane).DbClose, revel.AFTER)
	revel.InterceptMethod((*CPilot).Init, revel.BEFORE)
	revel.InterceptMethod((*CPilot).DbClose, revel.AFTER)
}