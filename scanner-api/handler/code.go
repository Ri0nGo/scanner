package handler

type RespCode int

const (
	RespCodeSuccess RespCode = iota
	RespCodeSystemInternalError
	RespCodeParamsError
)

var CodeMsg = map[RespCode]string{
	RespCodeSuccess:             "Success",
	RespCodeParamsError:         "Params Error",
	RespCodeSystemInternalError: "System Internal Error",
}

func (c RespCode) Msg() string {
	msg, ok := CodeMsg[c]
	if ok {
		return msg
	}
	return CodeMsg[RespCodeSystemInternalError]
}
