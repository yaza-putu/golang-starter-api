// Package response api
// default response
// code : 200
// status : true
// data : nil
// message : nil
package response

type (
	optFunc func(*res)

	res struct {
		Code    int `json:"code,omitempty"`
		Data    any `json:"data,omitempty"`
		Message any `json:"message,omitempty"`
		Errors  any `json:"errors,omitempty"`
	}

	DataApi struct {
		res
	}
)

func defaultResponse() res {
	return res{
		Code:    200,
		Data:    nil,
		Message: nil,
		Errors:  nil,
	}
}

func SetCode(code int) optFunc {
	return func(r *res) {
		r.Code = code
	}
}

func SetData(data any) optFunc {
	return func(r *res) {
		r.Data = data
	}
}

func SetMessage(msg any) optFunc {
	return func(r *res) {
		r.Message = msg
	}
}

func SetError(e any) optFunc {
	return func(r *res) {
		r.Errors = e
	}
}

func (r *res) GetCode() int {
	return r.Code
}

func (r *res) GetData() any {
	return r.Data
}

func (r *res) GetMessage() any {
	return r.Message
}

func (r *res) GetError() any {
	return r.Errors
}

func Api(opts ...optFunc) DataApi {
	o := defaultResponse()

	for _, fn := range opts {
		fn(&o)
	}

	return DataApi{
		res: o,
	}
}

func TimeOut() DataApi {
	return Api(SetCode(408), SetMessage("Request timeout or canceled by user"))
}

func BadRequest(err error) DataApi {
	return DataApi{
		res: res{
			Code:    400,
			Data:    nil,
			Message: "bad request",
			Errors:  err,
		},
	}
}
