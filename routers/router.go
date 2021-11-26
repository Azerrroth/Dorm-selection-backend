package routers

import (
	"go-gin-example/pkg/setting"
	v1 "go-gin-example/routers/api/v1"

	"go-gin-example/middleware/cors"
	"go-gin-example/middleware/jwt"

	// "github.com/gin-contrib/cors"

	// "go-gin-example/pkg/util"

	// "github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())
	r.Use(cors.Cors())
	// r.Use(sessions.Sessions("session", util.Store))

	gin.SetMode(setting.RunMode)

	apiv1 := r.Group("/api/v1")
	apiv1.POST("/login", v1.Login)
	apiv1.POST("/register", v1.Register)

	// apiv1.Use(session.Session())
	apiv1.Use(jwt.JWT())
	{
		apiv1.GET("/token", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "token access",
			})
		})

		dorm := apiv1.Group("/dorm")

		// Dorm route function
		dorm.GET("/user2RoomInfo", v1.GetUser2RoomInfo)
		dorm.GET("/buildings", v1.GetBuildingList)
		dorm.GET("/buildingStatus", v1.GetBuildingStatus)
		dorm.GET("/buildingsStatus", v1.GetBuildingsStatus)
		dorm.POST("/checkOutRoom", v1.CheckOutRoom)

		// apiv1.POST("/dorm", v1.GetDormList)
		// apiv1.GET("/buildings", v1.GetBuildingList)
		// apiv1.GET("/buildingStatus", v1.GetBuildingStatus)
		// apiv1.GET("/buildingsStatus", v1.GetBuildingsStatus)
		// apiv1.GET("/user2RoomInfo", v1.GetUser2RoomInfo)
		// apiv1.POST("/checkOutRoom", v1.CheckOutRoom)

		user := apiv1.Group("/user")

		// user route function
		user.POST("/updateUserProfile", v1.UpdateUserProfile)
		user.GET("/updateCertifyCode", v1.UpdateCertifyCode)

		// apiv1.POST("/updateUserProfile", v1.UpdateUserProfile)
		// apiv1.GET("/updateCertifyCode", v1.UpdateCertifyCode)

		order := apiv1.Group("/order")

		// order route function
		order.POST("/bookOrder", v1.BookOrder)

		// apiv1.POST("/bookOrder", v1.BookOrder)
	}

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})

	return r
}
