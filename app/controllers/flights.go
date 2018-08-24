package controllers

import (
	"app/app/models/FlightModel"
	"io/ioutil"
	"github.com/revel/revel"
	"app/app/lib/responce"
	"app/app/lib/dbManager"
	"fmt"
)



type CFlight struct {
	*revel.Controller
	provider FlightModel.FlightProvider
}

func(c *CFlight) DbInit()revel.Result{
	db, err :=  dbManager.OpenConnection()
	if err != nil {
		return c.RenderJson(responce.Failed(err))
	}
	c.provider = FlightModel.FlightProvider{Db: db}
	return nil
}

func(c *CFlight)DbClose()revel.Result{
	dbManager.CloseConnection(c.provider.Db)
	return nil
}


func (c *CFlight) GetAll() revel.Result {
	fmt.Println(c.provider.Db)
	returnedValue, err := c.provider.GetAll()
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))

}
func (c *CFlight) Add() revel.Result {
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

func (c *CFlight) Edit(id int) revel.Result {
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
func (c *CFlight) GetById(id int) revel.Result {
	returnedValue, err := c.provider.GetById(id)
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}

func (c *CFlight) Delete(id int) revel.Result {
	returnedValue, err := c.provider.Delete(id)
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}
