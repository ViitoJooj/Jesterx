package routers

import (
	"database/sql"

	"github.com/ViitoJooj/verkoupe/internal/domain/usecases"
	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
	"github.com/ViitoJooj/verkoupe/internal/port/persistence/repositories"
)

func NewControllers(db *sql.DB) *Controllers {
	addressRepo := repositories.NewAddressRepository(db)
	cupomRepo := repositories.NewCupomRepository(db)
	orgRepo := repositories.NewOrganizationRepository(db)
	phoneRepo := repositories.NewPhoneRepository(db)
	planRepo := repositories.NewPlanRepository(db)
	preparingShippingProductRepo := repositories.NewPreparingShippingProductRepository(db)
	productRepo := repositories.NewProductRepository(db)
	productShippedRepo := repositories.NewProductShippedRepository(db)
	productTagRepo := repositories.NewProductTagRepository(db)
	rbacRepo := repositories.NewRbacRepository(db)
	storageProductRepo := repositories.NewStorageProductRepository(db)
	termsRepo := repositories.NewTermsRepository(db)
	termsAcceptedRepo := repositories.NewTermsAcceptedRepository(db)
	userRepo := repositories.NewUserRepository(db)
	websiteRepo := repositories.NewWebsiteRepository(db)
	websiteComponentRepo := repositories.NewWebsiteComponentRepository(db)

	return &Controllers{
		Address:                  controllers.NewAddressController(usecases.NewCreateAddressUseCase(db, addressRepo), addressRepo),
		Auth:                     controllers.NewAuthController(usecases.NewRegisterUserUseCase(db, userRepo)),
		Cupom:                    controllers.NewCupomController(usecases.NewCreateCupomUseCase(db, cupomRepo)),
		Organization:             controllers.NewOrganizationController(usecases.NewCreateOrganizationUseCase(db, orgRepo), orgRepo),
		Phone:                    controllers.NewPhoneController(usecases.NewCreatePhoneUseCase(db, phoneRepo), phoneRepo),
		Plan:                     controllers.NewPlanController(usecases.NewCreatePlanUseCase(db, planRepo)),
		PreparingShippingProduct: controllers.NewPreparingShippingProductController(usecases.NewCreatePreparingShippingProductUseCase(db, preparingShippingProductRepo)),
		Product:                  controllers.NewProductController(usecases.NewCreateProductUseCase(db, productRepo)),
		ProductShipped:           controllers.NewProductShippedController(usecases.NewCreateProductShippedUseCase(db, productShippedRepo)),
		ProductTag:               controllers.NewProductTagController(usecases.NewCreateProductTagUseCase(db, productTagRepo)),
		Rbac:                     controllers.NewRbacController(usecases.NewCreateRbacUseCase(db, rbacRepo), rbacRepo),
		StorageProduct:           controllers.NewStorageProductController(usecases.NewCreateStorageProductUseCase(db, storageProductRepo)),
		Terms:                    controllers.NewTermsController(usecases.NewCreateTermsUseCase(db, termsRepo), termsRepo),
		TermsAccepted:            controllers.NewTermsAcceptedController(usecases.NewCreateTermsAcceptedUseCase(db, termsAcceptedRepo), termsAcceptedRepo),
		Website:                  controllers.NewWebsiteController(usecases.NewCreateWebsiteUseCase(db, websiteRepo)),
		WebsiteComponent:         controllers.NewWebsiteComponentController(usecases.NewCreateWebsiteComponentUseCase(db, websiteComponentRepo)),
	}
}
