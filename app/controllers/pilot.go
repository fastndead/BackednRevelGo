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

func(c *CPilot)Init()revel.Result{//инициализация контроллера
	if auth := c.Request.Header.Get("Authorization"); auth == ""{
		return c.Redirect("/")//если пользователь не авторизован - переадресация на главную страницу
	}
	err := c.provider.Init()//инициализация провайдера
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return nil
}

func(c *CPilot)DbClose()revel.Result{//закрытие базы
	err := c.provider.Close()
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return nil
}

func (c *CPilot) GetAll() revel.Result {//получение списка всех пилотов
	returnedValue, err := c.provider.GetAll()//получение данных из провайдера
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))

}
func (c *CPilot) Add() revel.Result {//добавление пилота
	params, err := ioutil.ReadAll(c.Request.Body)//считывание тела запроса
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	returnedValue, err := c.provider.Add(params)//передача данных в провайдер
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}

func (c *CPilot) Edit(id int) revel.Result {//редактирование пилота
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
func (c *CPilot) GetById(id int) revel.Result {//получение пилота по id
	returnedValue, err := c.provider.GetById(id)//получение данных из провайдера
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}

func (c *CPilot) Delete(id int) revel.Result {//удаление пилота
	returnedValue, err := c.provider.Delete(id)//передача id в провайдер
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}