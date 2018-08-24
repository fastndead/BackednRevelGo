package controllers

import (
	"github.com/revel/revel"
	"app/app/lib/responce"
	"io/ioutil"
	"app/app/models/PilotModel"
	"app/app/lib/dbManager"
)

type CPilot struct {
	*revel.Controller
	provider PilotModel.PilotProvider
}

func(c *CPilot) DbInit()revel.Result{
	db, err :=  dbManager.OpenConnection()
	if err != nil {
		return c.RenderJson(responce.Failed(err))
	}
	c.provider = PilotModel.PilotProvider{Db: db}
	return nil
}

func(c *CPilot)DbClose()revel.Result{
	dbManager.CloseConnection(c.provider.Db)
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