//+build wireinject

package config

import (
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"kpdigisign/activation"
	"kpdigisign/download_document"
	"kpdigisign/registration"
	"kpdigisign/send_document"
	"kpdigisign/sign_document"
	"kpdigisign/user"
)

func InjectUserController(db *gorm.DB) user.Controller {
	wire.Build(
		user.NewUserController,
		user.NewUserMapper,
		user.NewUserService,
		user.NewUserRepository,
	)
	return user.Controller{}
}

func InjectRegistrationController(db *gorm.DB) registration.Controller {
	wire.Build(
		registration.NewRegistrationController,
		registration.NewRegistrationMapper,
		registration.NewRegistrationService,
		registration.NewRegistrationRepository,
	)
	return registration.Controller{}
}

func InjectSendDocumentController(db *gorm.DB) send_document.Controller {
	wire.Build(
		send_document.NewSendDocumentController,
		send_document.NewSendDocumentMapper,
		send_document.NewSendDocumentService,
		send_document.NewSendDocumentRepository,
	)
	return send_document.Controller{}
}

func InjectDownloadDocumentController() download_document.Controller {
	wire.Build(download_document.NewDownloadDocumentController)
	return download_document.Controller{}
}

func InjectActivationController(db *gorm.DB) activation.Controller {
	wire.Build(activation.NewActivationController,
		activation.NewActivationMapper,
		activation.NewActivationService,
		activation.NewActivationRepository,
	)
	return activation.Controller{}
}

func InjectSignDocumentController(db *gorm.DB) sign_document.Controller {
	wire.Build(sign_document.NewSignDocumentController,
		sign_document.NewSignDocumentMapper,
		sign_document.NewSignDocumentService,
		sign_document.NewSignDocumentRepository,
	)
	return sign_document.Controller{}
}
