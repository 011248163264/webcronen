package controllers

type HelpController struct {
	BaseController
}

func (this *HelpController) Index() {

	this.Data["pageTitle"] = "Using help"
	this.display()
}
