package controllers

import (
	"github.com/revel/revel"
	"app/app/lib/responce"
	"io/ioutil"
	"app/app/models/PilotModel"
)

type CPilot struct {
	*revel.Controller
	provider PilotModel.PilotProvider
}

func(c *CPilot)Init()revel.Result{
	if auth := c.Request.Header.Get("Authorization"); auth == ""{
		return c.Redirect("/")
	}
	err := c.provider.Init()
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return nil
}

func(c *CPilot)DbClose()revel.Result{
	err := c.provider.Close()
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return nil
}

func (c *CPilot) GetAll() revel.Result {
	returnedValue, err := c.provider.GetAll()
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))

}
func (c *CPilot) Add() revel.Result {
	params, err := ioutil.ReadAll(c.Request.Body)
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	returnedValue, err := c.provider.Add(params)
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}

func (c *CPilot) Edit(id int) revel.Result {
	params, err := ioutil.ReadAll(c.Request.Body)
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	returnedValue, err := c.provider.Edit(id, params)
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}
func (c *CPilot) GetById(id int) revel.Result {
	returnedValue, err := c.provider.GetById(id)
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}

func (c *CPilot) Delete(id int) revel.Result {
	returnedValue, err := c.provider.Delete(id)
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}