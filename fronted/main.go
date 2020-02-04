package main

import (
	"github.com/gin-gonic/gin"
	"seckilling-practice-project/fronted/web/ctrls"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("fronted/web/views/**/*")
	router.Static("/public", "fronted/web/public")
	router.Static("/html", "fronted/web/htmlProductShow")
	userRou := router.Group("/user")
	{
		userRou.GET("/register", ctrls.UserRegisterShowAction)
		userRou.POST("/register", ctrls.UserRegisterAction)
		userRou.GET("/login", ctrls.UserLoginShowAction)
		userRou.POST("/login", ctrls.UserLoginAction)
	}

	productRou := router.Group("/product")
	{
		productRou.GET("/generate_html/:productID", ctrls.GenerateHtml)
	}

	router.Run(":8000")
}
