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

func(c *CPlane) Init()revel.Result{//инициализация контроллера
	if auth := c.Request.Header.Get("Authorization"); auth == ""{
		return c.Redirect("/")//если пользователь не авторизирован - переадресация на главную страницу
	}
	err := c.provider.Init()//инициализация провайдера
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return nil
}

func(c *CPlane)DbClose()revel.Result{//закрытие базы
	err := c.provider.Close()
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return nil
}

func (c *CPlane) GetAll() revel.Result {//получение списка самолётов
	returnedValue, err := c.provider.GetAll()//получение данных из провайдера
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))

}
func (c *CPlane) Add() revel.Result {//добавление самолёта
	params, err := ioutil.ReadAll(c.Request.Body)//считывание тела запроса
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	returnedValue, err := c.provider.Add(params)//пеедача данных в провайдер
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}

func (c *CPlane) Edit(id int) revel.Result {//редактирование данных о самолёте
	params, err := ioutil.ReadAll(c.Request.Body)//считывание тела запроса
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	returnedValue, err := c.provider.Edit(id, params)//передача данных в провайдер
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}
func (c *CPlane) GetById(id int) revel.Result {//получение самолёта по id
	returnedValue, err := c.provider.GetById(id)//переача id в провадйер, получение объекта
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}

func (c *CPlane) Delete(id int) revel.Result {//удаление самолёта из базы
	returnedValue, err := c.provider.Delete(id)//передача id в провайдер
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}