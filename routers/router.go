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

}
