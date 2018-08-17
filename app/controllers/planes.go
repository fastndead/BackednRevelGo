package controllers

import (
	"github.com/revel/revel"
	"io/ioutil"
	"app/app/lib/responce"
	"app/app/models/PlaneModel"
)
//import "strconv"

type CPlane struct {
	*revel.Controller
}

func (c CPlane) GetAll() revel.Result {
	returnedValue, err := PlaneModel.GetAll()
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))

}
func (c CPlane) Add() revel.Result {
	params, err := ioutil.ReadAll(c.Request.Body)
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	returnedValue, err := PlaneModel.Add(params)
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}

func (c CPlane) Edit(id int) revel.Result {
	params, err := ioutil.ReadAll(c.Request.Body)
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	returnedValue, err := PlaneModel.Edit(id, params)
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}
func (c CPlane) GetById(id int) revel.Result {
	returnedValue, err := PlaneModel.GetById(id)
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}

func (c CPlane) Delete(id int) revel.Result {
	returnedValue, err := PlaneModel.Delete(id)
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}