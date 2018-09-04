package responce

type Responce struct {
	Err  			bool //наличие ошибки
	ErrorText		string//текст ошибки
	Data 			interface{}//данные
}

func Success(Data interface{})Responce{//метода, в случае отсутствия ошибки
	return Responce{Err:false, ErrorText:"", Data:Data}
}

func Failed(err error)Responce{//метод, в случае наличия ошибки
return Responce{true, err.Error(), nil}
}