package controllers

import "github.com/revel/revel"
import (
	_ "github.com/lib/pq"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.RenderText("fine")
}

func init() {
	revel.InterceptMethod((*CFlight).DbInit, revel.BEFORE)
	revel.InterceptMethod((*CFlight).DbClose, revel.AFTER)
	revel.InterceptMethod((*CPlane).DbInit, revel.BEFORE)
	revel.InterceptMethod((*CPlane).DbClose, revel.AFTER)
	revel.InterceptMethod((*CPilot).DbInit, revel.BEFORE)
	revel.InterceptMethod((*CPilot).DbClose, revel.AFTER)
}