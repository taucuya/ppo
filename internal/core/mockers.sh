#!/bin/bash

mockgen -source=service/auth/auth.go -destination=mock_structs/auth_mock.go -package=auth_mock
mockgen -source=service/user/user.go -destination=mock_structs/user_mock.go -package=mock_structs
mockgen -source=service/basket/basket.go -destination=mock_structs/basket_mock.go -package=mock_structs
mockgen -source=service/brand/brand.go -destination=mock_structs/brand_mock.go -package=mock_structs
mockgen -source=service/order/order.go -destination=mock_structs/order_mock.go -package=mock_structs
mockgen -source=service/product/product.go -destination=mock_structs/product_mock.go -package=mock_structs
mockgen -source=service/review/review.go -destination=mock_structs/review_mock.go -package=mock_structs
mockgen -source=service/worker/worker.go -destination=mock_structs/worker_mock.go -package=mock_structs