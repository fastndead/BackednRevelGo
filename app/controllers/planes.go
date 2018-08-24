package controllers

import (
	"github.com/revel/revel"
	"io/ioutil"
	"app/app/lib/responce"
	"app/app/models/PlaneModel"
)

type CPlane struct {
	*revel.Controller
	provider PlaneModel.PlaneProvider
}

func(c *CPlane) DbInit()revel.Result{
	err := c.provider.Init()
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return nil
}

func(c *CPlane)DbClose()revel.Result{
	err := c.provider.Close()
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return nil
}

func (c *CPlane) GetAll() revel.Result {
	returnedValue, err := c.provider.GetAll()
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))

}
func (c *CPlane) Add() revel.Result {
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

func (c *CPlane) Edit(id int) revel.Result {
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
func (c *CPlane) GetById(id int) revel.Result {
	returnedValue, err := c.provider.GetById(id)
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}

func (c *CPlane) Delete(id int) revel.Result {
	returnedValue, err := c.provider.Delete(id)
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}