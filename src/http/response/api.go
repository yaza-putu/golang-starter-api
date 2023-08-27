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
		Code    int  `json:"code"`
		Status  bool `json:"status"`
		Data    any  `json:"data"`
		Message any  `json:"message"`
	}

	DataApi struct {
		res
	}
)

func defaultResponse() res {
	return res{
		Code:    200,
		Status:  true,
		Data:    nil,
		Message: nil,
	}
}

func SetCode(code int) optFunc {
	return func(r *res) {
		if code != 200 {
			r.Status = false
		}
		r.Code = code
	}
}

func SetStatus(status bool) optFunc {
	return func(r *res) {
		if r.Status == false {
			r.Code = 500
		}
		r.Status = status
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

func (r *res) GetStatus() bool {
	return r.Status
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

func Api(opts ...optFunc) DataApi {
	o := defaultResponse()

	for _, fn := range opts {
		fn(&o)
	}

	return DataApi{
		res: o,
	}
}
