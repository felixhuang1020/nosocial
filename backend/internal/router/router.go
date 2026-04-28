package router

import (
	"nosocial/internal/bootstrap"
	"nosocial/internal/dao"
	"nosocial/internal/handler"
	adminhandler "nosocial/internal/handler/admin"
	wxhandler "nosocial/internal/handler/wx"
	"nosocial/internal/middleware"
	"nosocial/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(app *bootstrap.App) *gin.Engine {
	if app.Config.App.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	// 信任边界：仅信任显式配置的反向代理 IP/CIDR
	// 未配置时传 nil，Gin 将不解析任何 X-Forwarded-For，ClientIP() 回落到 RemoteAddr
	if len(app.Config.App.TrustedProxies) > 0 {
		_ = r.SetTrustedProxies(app.Config.App.TrustedProxies)
	} else {
		_ = r.SetTrustedProxies(nil)
	}
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// 初始化DAO
	userDAO := dao.NewUserDAO(app.DB)
	shareholderOrderDAO := dao.NewShareholderOrderDAO(app.DB)
	commissionDAO := dao.NewCommissionRecordDAO(app.DB)
	withdrawalDAO := dao.NewWithdrawalDAO(app.DB)
	tarotCardDAO := dao.NewTarotCardDAO(app.DB)
	tarotReadingDAO := dao.NewTarotReadingDAO(app.DB)
	tarotMappingDAO := dao.NewTarotDrinkMappingDAO(app.DB)
	drinkDAO := dao.NewDrinkDAO(app.DB)
	categoryDAO := dao.NewDrinkCategoryDAO(app.DB)
	orderDAO := dao.NewDrinkOrderDAO(app.DB)
	birthdayDAO := dao.NewBirthdayGiftDAO(app.DB)
	reviewDAO := dao.NewReviewDAO(app.DB)
	couponDAO := dao.NewCouponDAO(app.DB)
	bannerDAO := dao.NewBannerDAO(app.DB)
	adminDAO := dao.NewAdminUserDAO(app.DB)
	settingDAO := dao.NewSettingDAO(app.DB)
	paymentNotifyDAO := dao.NewPaymentNotifyLogDAO(app.DB)

	// 初始化Service
	userService := service.NewUserService(userDAO, &app.Config.WX)
	shareholderService := service.NewShareholderService(userDAO, shareholderOrderDAO, commissionDAO, withdrawalDAO, app.WXPay, app.DB)
	tarotService := service.NewTarotService(tarotCardDAO, tarotReadingDAO, tarotMappingDAO, drinkDAO)
	drinkService := service.NewDrinkService(drinkDAO, categoryDAO)
	orderService := service.NewOrderService(orderDAO, drinkDAO, couponDAO, userDAO, app.WXPay)
	reviewService := service.NewReviewService(reviewDAO, couponDAO, userDAO, app.DB)
	couponService := service.NewCouponService(couponDAO, app.DB)
	birthdayService := service.NewBirthdayService(birthdayDAO, userDAO, app.DB)
	bannerService := service.NewBannerService(bannerDAO)
	adminService := service.NewAdminService(adminDAO, userDAO, orderDAO, reviewDAO)
	uploadService := service.NewUploadService()
	systemService := service.NewSystemService(settingDAO)

	// 初始化默认配置
	_ = systemService.InitDefaultSettings()

	// 初始化Handler
	publicHandler := handler.NewPublicHandler(bannerService, drinkService, tarotService, systemService)
	paymentNotifyHandler := handler.NewPaymentNotifyHandler(app.WXPay, orderService, shareholderService, paymentNotifyDAO, app.Logger)
	wxUserHandler := wxhandler.NewUserHandler(userService)
	wxShareholderHandler := wxhandler.NewShareholderHandler(shareholderService)
	wxTarotHandler := wxhandler.NewTarotHandler(tarotService)
	wxBirthdayHandler := wxhandler.NewBirthdayHandler(birthdayService)
	wxReviewHandler := wxhandler.NewReviewHandler(reviewService)
	wxCouponHandler := wxhandler.NewCouponHandler(couponService)
	wxOrderHandler := wxhandler.NewOrderHandler(orderService)
	wxUploadHandler := wxhandler.NewUploadHandler(uploadService)

	adminHandler := adminhandler.NewAdminHandler(adminService)
	userAdminHandler := adminhandler.NewUserAdminHandler(userService)
	shareholderAdminHandler := adminhandler.NewShareholderAdminHandler(shareholderService, userService)
	drinkAdminHandler := adminhandler.NewDrinkAdminHandler(drinkService)
	tarotAdminHandler := adminhandler.NewTarotAdminHandler(tarotService)
	bannerAdminHandler := adminhandler.NewBannerAdminHandler(bannerService)
	reviewAdminHandler := adminhandler.NewReviewAdminHandler(reviewService)
	orderAdminHandler := adminhandler.NewOrderAdminHandler(orderService)
	uploadAdminHandler := adminhandler.NewUploadAdminHandler(uploadService)
	settingAdminHandler := adminhandler.NewSettingAdminHandler(systemService)

	api := r.Group("/api/v1")

	// 公共接口（通用限流）
	pub := api.Group("/public")
	pub.Use(middleware.RateLimitGeneral())
	{
		pub.GET("/banners", publicHandler.ListBanners)
		pub.GET("/drinks", publicHandler.ListDrinks)
		pub.GET("/drinks/:id", publicHandler.GetDrinkDetail)
		pub.GET("/tarot/cards", publicHandler.ListTarotCards)
		pub.GET("/config", publicHandler.GetConfig)
		pub.GET("/shop", publicHandler.GetShopInfo)
	}

	// 微信支付 V3 回调（公开接口，微信服务器调用）
	api.POST("/payment/wxpay/notify", paymentNotifyHandler.Notify)

	// 微信小程序接口
	wx := api.Group("/wx")
	{
		wx.POST("/login", middleware.RateLimitLogin(), wxUserHandler.Login)
		wxAuth := wx.Group("")
		wxAuth.Use(middleware.WXAuth(), middleware.RateLimitGeneral())
		{
			wxAuth.GET("/user", wxUserHandler.GetProfile)
			wxAuth.PUT("/user/birthday", wxUserHandler.SetBirthday)

			wxAuth.POST("/shareholder/apply", wxShareholderHandler.Apply)
			wxAuth.POST("/shareholder/pay", wxShareholderHandler.Pay)
			wxAuth.GET("/shareholder/profile", wxShareholderHandler.Profile)
			wxAuth.GET("/shareholder/earnings", wxShareholderHandler.Earnings)
			wxAuth.GET("/shareholder/team", wxShareholderHandler.Team)
			wxAuth.POST("/shareholder/withdraw", wxShareholderHandler.Withdraw)

			wxAuth.POST("/tarot/divine", wxTarotHandler.Divine)
			wxAuth.GET("/tarot/history", wxTarotHandler.History)

			wxAuth.GET("/birthday/gift", wxBirthdayHandler.CheckGift)
			wxAuth.POST("/birthday/claim", wxBirthdayHandler.Claim)
			wxAuth.GET("/birthday/history", wxBirthdayHandler.History)

			wxAuth.POST("/review/submit", wxReviewHandler.Submit)
			wxAuth.GET("/review/status", wxReviewHandler.Status)

			wxAuth.GET("/coupons", wxCouponHandler.MyCoupons)
			wxAuth.POST("/coupons/use", wxCouponHandler.UseCoupon)

			wxAuth.POST("/orders", wxOrderHandler.Create)
			wxAuth.GET("/orders", wxOrderHandler.List)
			wxAuth.POST("/orders/:id/pay", wxOrderHandler.Pay)
			wxAuth.POST("/orders/:id/cancel", wxOrderHandler.Cancel)

			wxAuth.POST("/free-drink/claim", wxUserHandler.ClaimFreeDrink)

			wxAuth.POST("/upload", wxUploadHandler.Upload)
			wxAuth.GET("/upload/signature", middleware.RateLimitUpload(), wxUploadHandler.GetOSSSignature)
			wxAuth.PUT("/upload/inline", wxUploadHandler.FixObjectInline)
		}
	}

	// 商家管理端接口
	admin := api.Group("/admin")
	{
		admin.POST("/login", middleware.RateLimitLogin(), adminHandler.Login)
		adminAuth := admin.Group("")
		adminAuth.Use(middleware.AdminAuth(), middleware.RateLimitGeneral())
		{
			adminAuth.GET("/dashboard", adminHandler.Dashboard)

			adminAuth.GET("/users", userAdminHandler.List)
			adminAuth.PUT("/users/:id/status", userAdminHandler.UpdateStatus)

			adminAuth.GET("/shareholders", shareholderAdminHandler.List)
			adminAuth.GET("/shareholders/:id/earnings", shareholderAdminHandler.Earnings)
			adminAuth.GET("/shareholders/:id/team", shareholderAdminHandler.Team)

			adminAuth.GET("/drinks", drinkAdminHandler.List)
			adminAuth.POST("/drinks", drinkAdminHandler.Create)
			adminAuth.PUT("/drinks/:id", drinkAdminHandler.Update)
			adminAuth.DELETE("/drinks/:id", drinkAdminHandler.Delete)
			adminAuth.GET("/categories", drinkAdminHandler.Categories)
			adminAuth.POST("/categories", drinkAdminHandler.CreateCategory)
			adminAuth.PUT("/categories/:id", drinkAdminHandler.UpdateCategory)
			adminAuth.DELETE("/categories/:id", drinkAdminHandler.DeleteCategory)

			adminAuth.GET("/tarot/mappings", tarotAdminHandler.Mappings)
			adminAuth.PUT("/tarot/mappings", tarotAdminHandler.UpdateMapping)

			adminAuth.GET("/banners", bannerAdminHandler.List)
			adminAuth.POST("/banners", bannerAdminHandler.Create)
			adminAuth.PUT("/banners/:id", bannerAdminHandler.Update)
			adminAuth.DELETE("/banners/:id", bannerAdminHandler.Delete)

			adminAuth.GET("/reviews", reviewAdminHandler.List)
			adminAuth.POST("/reviews/:id/audit", reviewAdminHandler.Audit)

			adminAuth.GET("/orders", orderAdminHandler.List)
			adminAuth.PUT("/orders/:id/status", orderAdminHandler.UpdateStatus)

			adminAuth.GET("/upload/signature", middleware.RateLimitUpload(), uploadAdminHandler.GetOSSSignature)
			adminAuth.PUT("/upload/inline", uploadAdminHandler.FixObjectInline)

			adminAuth.GET("/settings", settingAdminHandler.GetSettings)
			adminAuth.PUT("/settings", settingAdminHandler.UpdateSettings)
		}
	}

	return r
}
