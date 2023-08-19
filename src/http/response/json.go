// Package response json
// default response
// code : 200
// status : true
// data : nil
// message : nil
package response

type (
	optFunc func(*res)

	res struct {
		Code    uint16 `json:"code"`
		Status  bool   `json:"status"`
		Data    any    `json:"data"`
		Message any    `json:"message"`
	}

	ResData struct {
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

func SetCode(code uint16) optFunc {
	return func(r *res) {
		r.Code = code
	}
}

func SetStatus(status bool) optFunc {
	return func(r *res) {
		if r.Code == 200 {
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

func (r *res) GetCode() uint16 {
	return r.Code
}

func (r *res) GetData() any {
	return r.Data
}

func (r *res) GetMessage() any {
	return r.Message
}

func Json(opts ...optFunc) *ResData {
	o := defaultResponse()

	for _, fn := range opts {
		fn(&o)
	}

	return &ResData{
		res: o,
	}
}
