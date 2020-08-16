package permissioncontroller

import "sso/config/env"

type PermissionController struct {
	env *env.Env
}

func NewPermissionController(env *env.Env) *PermissionController {
	return &PermissionController{env: env}
}

func (p *PermissionController) Index()  {

}

func (p *PermissionController) Store()  {

}

func (p *PermissionController) Show()  {

}

func (p *PermissionController) Update()  {

}

func (p *PermissionController) Destroy()  {

}
