package loraoss

import (
	"github.com/go-resty/resty/v2"
	"github.com/lishimeng/go-connector/loraoss/model"
)

type Token struct {
	Jwt string `json:"jwt"`
}

type ConnectorConfig struct {
	Host     string
	UserName string
	Password string
}

type Connector interface {
	Login() (token Token, err error)

	Request() *resty.Request
}

type Gateway interface {
	Create(param model.GatewayForm) (code int, err error)
	Delete(id string) (int, error)
	Edit()
	List()
}

type Application interface {
	Create()
	Delete()
	Edit()
	List()
}

type Device interface {
	Create(device model.DeviceForm) (code int, err error)
	Delete(devEUI string) (int, error)
	Edit(device model.DeviceForm) (code int, err error)
	List(*model.DeviceRequestBuilder) (model.DevicePage, error)

	// OTAA key
	GetOTAAKeys(devEUI string) (keys model.DeviceKeys, code int, err error)
	SetOTAAKeys(keys model.DeviceKeys) (code int, err error)
	UpdateOTAAKeys(keys model.DeviceKeys) (code int, err error)
}

type DeviceProfile interface {
	List() (dps model.DeviceProfilePage, err error)
}
