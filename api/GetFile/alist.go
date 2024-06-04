package GetFile

import (
	"github.com/gin-gonic/gin"
	"github.com/msterzhang/onelist/api/database"
	"github.com/msterzhang/onelist/api/models"
	"github.com/msterzhang/onelist/plugins/alist"
	"net/url"
	"strings"
)

func Get(c *gin.Context) {

	db := database.NewDb()

	id := c.Query("id")
	gallery_type := c.Query("gallery_type")
	gallery := models.Gallery{}

	if gallery_type == "tv" {
		thetvDb := models.TheTv{}
		err := db.Model(&models.TheTv{}).Where("id = ?", id).First(&thetvDb).Error
		if err != nil {
			c.JSON(200, gin.H{"code": 201, "msg": "没有查询到资源!"})
			return
		}
		err = db.Model(&models.Gallery{}).Where("gallery_uid = ?", thetvDb.GalleryUid).First(&gallery).Error
		if err != nil {
			c.JSON(200, gin.H{"code": 201, "msg": "没有查询到对应影视库!"})
			return
		}
	} else {
		thetvDb := models.TheMovie{}
		err := db.Model(&models.TheMovie{}).Where("id = ?", id).First(&thetvDb).Error
		if err != nil {
			c.JSON(200, gin.H{"code": 201, "msg": "没有查询到资源!"})
			return
		}
		err = db.Model(&models.Gallery{}).Where("gallery_uid = ?", thetvDb.GalleryUid).First(&gallery).Error
		if err != nil {
			c.JSON(200, gin.H{"code": 201, "msg": "没有查询到对应影视库!"})
			return
		}
	}

	request := &models.Request{}
	_ = c.ShouldBindJSON(request)
	_tv_path, _ := url.QueryUnescape(request.Data)
	tv_path_list := strings.Split(_tv_path, "/d")
	tv_path := tv_path_list[len(tv_path_list)-1]
	tv_path = strings.Split(tv_path, "?")[0]

	file_data, _ := alist.AlistFileUrl(gallery, tv_path)
	c.JSON(200, gin.H{"code": 200, "data": file_data})
	return

}
