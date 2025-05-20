package controller

import (
	"github.com/taucuya/ppo/internal/core/service/auth"
	"github.com/taucuya/ppo/internal/core/service/basket"
	"github.com/taucuya/ppo/internal/core/service/brand"
	"github.com/taucuya/ppo/internal/core/service/favourites"
	"github.com/taucuya/ppo/internal/core/service/order"
	"github.com/taucuya/ppo/internal/core/service/product"
	"github.com/taucuya/ppo/internal/core/service/review"
	"github.com/taucuya/ppo/internal/core/service/user"
	"github.com/taucuya/ppo/internal/core/service/worker"
)

type Controller struct {
	AuthServise       auth.Service
	BasketService     basket.Service
	BrandService      brand.Service
	FavouritesService favourites.Service
	OrderService      order.Service
	ProductService    product.Service
	ReviewService     review.Service
	UserService       user.Service
	WorkerService     worker.Service
}

func New(a auth.Service, ba basket.Service,
	br brand.Service, f favourites.Service, o order.Service, p product.Service,
	r review.Service, u user.Service, w worker.Service) *Controller {
	return &Controller{AuthServise: a, BasketService: ba, BrandService: br,
		OrderService: o, ProductService: p, ReviewService: r, UserService: u, WorkerService: w}
}
