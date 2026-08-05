package routers

import (
	"net/http"

	"github.com/ViitoJooj/verkoupe/internal/port/http/controllers"
)

type Controllers struct {
	Address                  *controllers.AddressController
	Auth                     *controllers.AuthController
	Cupom                    *controllers.CupomController
	Organization             *controllers.OrganizationController
	Phone                    *controllers.PhoneController
	Plan                     *controllers.PlanController
	PreparingShippingProduct *controllers.PreparingShippingProductController
	Product                  *controllers.ProductController
	ProductShipped           *controllers.ProductShippedController
	ProductTag               *controllers.ProductTagController
	Rbac                     *controllers.RbacController
	StorageProduct           *controllers.StorageProductController
	Terms                    *controllers.TermsController
	TermsAccepted            *controllers.TermsAcceptedController
	Website                  *controllers.WebsiteController
	WebsiteComponent         *controllers.WebsiteComponentController
}

func Register(mux *http.ServeMux, c *Controllers) {
	RegisterAddressRoutes(mux, c.Address)
	RegisterAuthRoutes(mux, c.Auth)
	RegisterCupomRoutes(mux, c.Cupom)
	RegisterOrganizationRoutes(mux, c.Organization)
	RegisterPhoneRoutes(mux, c.Phone)
	RegisterPlanRoutes(mux, c.Plan)
	RegisterPreparingShippingProductRoutes(mux, c.PreparingShippingProduct)
	RegisterProductRoutes(mux, c.Product)
	RegisterProductShippedRoutes(mux, c.ProductShipped)
	RegisterProductTagRoutes(mux, c.ProductTag)
	RegisterRbacRoutes(mux, c.Rbac)
	RegisterStorageProductRoutes(mux, c.StorageProduct)
	RegisterTermsRoutes(mux, c.Terms)
	RegisterTermsAcceptedRoutes(mux, c.TermsAccepted)
	RegisterWebsiteRoutes(mux, c.Website)
	RegisterWebsiteComponentRoutes(mux, c.WebsiteComponent)
}
