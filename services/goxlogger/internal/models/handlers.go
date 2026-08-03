package models

type BaseHandler struct {
	app AppInterface
}

func (h *BaseHandler) SetApp(app AppInterface) {
	h.app = app
}

func (h *BaseHandler) GetApp() AppInterface {
	return h.app
}
