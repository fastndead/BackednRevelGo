package controllers

import (
	"app/app/models/FlightModel"
	"io/ioutil"
	"github.com/revel/revel"
	"app/app/lib/responce"
)



type CFlight struct {
	*revel.Controller
	provider FlightModel.FlightProvider
}

func(c *CFlight) Init()revel.Result{//инициализация контроллера
	if auth := c.Request.Header.Get("Authorization"); auth == ""{//проверка авторизации
		return c.Redirect("/")//если не авторизирован - перенаправление на главную страницу с проверкой авторизации
	}
	err := c.provider.Init()//инициализация провайдера
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return nil
}

func(c *CFlight)DbClose()revel.Result{//закрытие базы данных провайдера
	err := c.provider.Close()
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return nil
}


func (c *CFlight) GetAll() revel.Result {//получение списка рейсов
	returnedValue, err := c.provider.GetAll()//запрос данных из провайдера
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))

}
func (c *CFlight) Add() revel.Result {//добавление рейса
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

func (c *CFlight) Edit(id int) revel.Result {//редактирование рейса
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
func (c *CFlight) GetById(id int) revel.Result {//получение рейса по id
	returnedValue, err := c.provider.GetById(id)//получение данных из провайдера
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}

func (c *CFlight) Delete(id int) revel.Result {//удаление рейса
	returnedValue, err := c.provider.Delete(id)// передача id в провайдер
	if err != nil{
		return c.RenderJson(responce.Failed(err))
	}
	return c.RenderJson(responce.Success(returnedValue))
}
