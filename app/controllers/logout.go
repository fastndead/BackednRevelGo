package controllers

import "app/app/lib/auth"
import (
	"github.com/revel/revel"
	"errors"
)

type CLogOut struct {
	*revel.Controller
}


func (c *CLogOut)LogOut() revel.Result{
	auth.LogOut(c.Controller)
	return c.RenderError(errors.New("401: You're not authorized"))
}
