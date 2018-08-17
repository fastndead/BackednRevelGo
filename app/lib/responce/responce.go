package responce

type Responce struct {
	Err  			bool
	ErrorText		string
	Data 			interface{}
}

func Success(Data interface{})Responce{
	return Responce{Err:false, ErrorText:"", Data:Data}
}

func Failed(err error)Responce{
return Responce{true, err.Error(), nil}
}