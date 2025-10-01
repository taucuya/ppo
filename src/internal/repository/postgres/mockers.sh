#!/bin/bash

mockgen -source=reps/auth/auth_interface.go -destination=mocs/auth_mock.go -package=mocks
mockgen -source=reps/user/user_interface.go -destination=mocks/user_mock.go -package=mocks
mockgen -source=reps/basket/basket_interface.go -destination=mocks/basket_mock.go -package=mocks
mockgen -source=reps/brand/brand_interface.go -destination=mocks/brand_mock.go -package=mocks
mockgen -source=reps/favourites/favourites_interface.go -destination=mocks/favourites_mock.go -package=mocks
mockgen -source=reps/order/order_interface.go -destination=mocks/order_mock.go -package=mocks
mockgen -source=reps/product/product_interface.go -destination=mocks/product_mock.go -package=mocks
mockgen -source=reps/review/review_interface.go -destination=mocks/review_mock.go -package=mocks
mockgen -source=reps/worker/worker_interface.go -destination=mocks/worker_mock.go -package=mocks