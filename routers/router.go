package routers

import (
	"catapi/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/random-cat", &controllers.MainController{}, "get:GetRandomCats")
	beego.Router("/api/breeds", &controllers.MainController{}, "get:GetCatBreeds")
	beego.Router("/api/votes", &controllers.MainController{}, "post:RecordVote")
	beego.Router("/api/votes", &controllers.MainController{}, "get:GetVotes")
	beego.Router("/api/favorites", &controllers.MainController{}, "post:AddFavorite")
	beego.Router("/api/favorites", &controllers.MainController{}, "get:GetFavorites")
	beego.Router("/api/favorites/:favoriteId", &controllers.MainController{}, "delete:DeleteFavorite")
	beego.Router("/api/config", &controllers.MainController{}, "get:GetConfig")

}
