package param

type ReqCreateDeviceV1 struct {
	Name   string `json:"name"`
	Serial string `json:"serial" validate:"required"`
}
