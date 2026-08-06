package main

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/domain/usecases"
	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
	"github.com/ViitoJooj/verkoupe/internal/port/http/middleware"
	"github.com/ViitoJooj/verkoupe/internal/port/http/routers"
	"github.com/ViitoJooj/verkoupe/internal/port/persistence/repositories"
	"github.com/ViitoJooj/verkoupe/pkg/dotenv"
	"github.com/ViitoJooj/verkoupe/pkg/logger"
	"github.com/ViitoJooj/verkoupe/pkg/postgresql"
	"github.com/ViitoJooj/verkoupe/pkg/server"
)

func main() {
	cfg, err := dotenv.Conn()
	if err != nil {
		logger.Fatal(err).Print()
	}

	db, err := postgresql.Conn(cfg.PostgreSQL)
	if err != nil {
		logger.Fatal(err).Print()
	}

	err = postgresql.Migrations(db)
	if err != nil {
		logger.Fatal(err).Print()
	}

	mux := http.NewServeMux()

	corsMiddleware := middleware.CORSMiddleware(cfg.Application.ViewUrl)
	authMiddleware := middleware.AuthMiddleware(cfg.Security.PasetoSecretKey)

	addressRepository := repositories.NewAddressRepository(db)
	createAddressUseCase := usecases.NewCreateAddressUseCase(addressRepository)
	addressController := controllers.NewAddressController(createAddressUseCase)
	routers.RegisterAddressRoutes(mux, addressController, corsMiddleware, authMiddleware)

	userRepository := repositories.NewUserRepository(db)
	authUseCase := usecases.NewAuthUseCase(userRepository)
	authController := controllers.NewAuthController(authUseCase)
	routers.RegisterAuthRoutes(mux, authController, corsMiddleware, authMiddleware)

	cupomRepository := repositories.NewCupomRepository(db)
	createCupomUseCase := usecases.NewCreateCupomUseCase(cupomRepository)
	cupomController := controllers.NewCupomController(createCupomUseCase)
	routers.RegisterCupomRoutes(mux, cupomController, corsMiddleware, authMiddleware)

	organizationRepository := repositories.NewOrganizationRepository(db)
	organizationUseCase := usecases.NewCreateOrganizationUseCase(organizationRepository)
	organizationController := controllers.NewOrganizationController(organizationUseCase)
	routers.RegisterOrganizationRoutes(mux, organizationController, corsMiddleware, authMiddleware)

	phoneRepository := repositories.NewPhoneRepository(db)
	createPhoneUseCase := usecases.NewCreatePhoneUseCase(phoneRepository)
	phoneController := controllers.NewPhoneController(createPhoneUseCase)
	routers.RegisterPhoneRoutes(mux, phoneController, corsMiddleware, authMiddleware)

	planRepository := repositories.NewPlanRepository(db)
	createPlanUseCase := usecases.NewCreatePlanUseCase(planRepository)
	planController := controllers.NewPlanController(createPlanUseCase)
	routers.RegisterPlanRoutes(mux, planController, corsMiddleware, authMiddleware)

	preparingShippingProductRepository := repositories.NewPreparingShippingProductRepository(db)
	createPreparingShippingProductUseCase := usecases.NewCreatePreparingShippingProductUseCase(preparingShippingProductRepository)
	preparingShippingProductController := controllers.NewPreparingShippingProductController(createPreparingShippingProductUseCase)
	routers.RegisterPreparingShippingProductRoutes(mux, preparingShippingProductController, corsMiddleware, authMiddleware)

	productRepository := repositories.NewProductRepository(db)
	createProductUseCase := usecases.NewCreateProductUseCase(productRepository)
	productController := controllers.NewProductController(createProductUseCase)
	routers.RegisterProductRoutes(mux, productController, corsMiddleware, authMiddleware)

	productShippedRepository := repositories.NewProductShippedRepository(db)
	createProductShippedUseCase := usecases.NewCreateProductShippedUseCase(productShippedRepository)
	productShippedController := controllers.NewProductShippedController(createProductShippedUseCase)
	routers.RegisterProductShippedRoutes(mux, productShippedController, corsMiddleware, authMiddleware)

	productTagRepository := repositories.NewProductTagRepository(db)
	createProductTagUseCase := usecases.NewCreateProductTagUseCase(productTagRepository)
	productTagController := controllers.NewProductTagController(createProductTagUseCase)
	routers.RegisterProductTagRoutes(mux, productTagController, corsMiddleware, authMiddleware)

	rbacRepository := repositories.NewRbacRepository(db)
	createRbacUseCase := usecases.NewCreateRbacUseCase(rbacRepository)
	rbacController := controllers.NewRbacController(createRbacUseCase)
	routers.RegisterRbacRoutes(mux, rbacController, corsMiddleware, authMiddleware)

	storageProductRepository := repositories.NewStorageProductRepository(db)
	createStorageProductUseCase := usecases.NewCreateStorageProductUseCase(storageProductRepository)
	storageProductController := controllers.NewStorageProductController(createStorageProductUseCase)
	routers.RegisterStorageProductRoutes(mux, storageProductController, corsMiddleware, authMiddleware)

	termsRepository := repositories.NewTermsRepository(db)
	createTermsUseCase := usecases.NewCreateTermsUseCase(termsRepository)
	termsController := controllers.NewTermsController(createTermsUseCase)
	routers.RegisterTermsRoutes(mux, termsController, corsMiddleware, authMiddleware)

	termsAcceptedRepository := repositories.NewTermsAcceptedRepository(db)
	createTermsAcceptedUseCase := usecases.NewCreateTermsAcceptedUseCase(termsAcceptedRepository)
	termsAcceptedController := controllers.NewTermsAcceptedController(createTermsAcceptedUseCase)
	routers.RegisterTermsAcceptedRoutes(mux, termsAcceptedController, corsMiddleware, authMiddleware)

	websiteRepository := repositories.NewWebsiteRepository(db)
	createWebsiteUseCase := usecases.NewCreateWebsiteUseCase(websiteRepository)
	websiteController := controllers.NewWebsiteController(createWebsiteUseCase)
	routers.RegisterWebsiteRoutes(mux, websiteController, corsMiddleware, authMiddleware)

	websiteComponentRepository := repositories.NewWebsiteComponentRepository(db)
	createWebsiteComponentUseCase := usecases.NewCreateWebsiteComponentUseCase(websiteComponentRepository)
	websiteComponentController := controllers.NewWebsiteComponentController(createWebsiteComponentUseCase)
	routers.RegisterWebsiteComponentRoutes(mux, websiteComponentController, corsMiddleware, authMiddleware)

	server.Start(cfg.Application.Port, mux)
}
