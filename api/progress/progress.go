package progress

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/msterzhang/onelist/api/database"
	"github.com/msterzhang/onelist/api/models"
	"gorm.io/gorm"
	"net/url"
	"strconv"
	"strings"
)

func Get(c *gin.Context) {
	db := database.NewDb()
	UserId := c.Request.Header.Get("UserId")
	TvId := c.Query("tv_id")
	SeasonId := c.Query("season_id")
	if UserId == "" {
		c.JSON(400, gin.H{"msg": "未获取到所需字段"})
		return
	}
	err := db.Model(&models.User{}).Where("UserId = ?", UserId).First(models.User{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(400, gin.H{"msg": "该用户不存在"})
		return
	}
	if TvId == "" {
		progress := &models.Progress{}
		err = db.Model(&models.Progress{}).Where("user_id = ?", UserId).Find(progress).Error
		if err != nil {
			c.JSON(200, gin.H{})
			return
		}
		c.JSON(200, gin.H{"data": progress.Data})
		return
	} else {
		progress := &models.ProgressTv{}
		err = db.Model(&models.ProgressTv{}).Where("user_id = ? and tv_id = ? and season_id = ?", UserId, TvId, SeasonId).First(progress).Error
		if err != nil {
			c.JSON(200, gin.H{})
			return
		}
		c.JSON(200, gin.H{"data": progress})
		return
	}
}

func Post(c *gin.Context) {
	db := database.NewDb()
	UserId := c.Request.Header.Get("UserId")
	if UserId == "" {
		c.JSON(400, gin.H{"msg": "未获取到所需字段"})
		return
	}
	err := db.Model(&models.User{}).Where("UserId = ?", UserId).First(models.User{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(400, gin.H{"msg": "该用户不存在"})
		return
	}
	request := &models.Request{}
	err = c.ShouldBindJSON(request)
	progress := &models.ProgressTv{}
	_tv_path, err := url.QueryUnescape(request.Data)
	tv_path_list := strings.Split(_tv_path, "/d")
	tv_path := "/d" + tv_path_list[len(tv_path_list)-1]
	err = db.Model(&models.ProgressTv{}).Where("user_id = ? and tv_path = ?", UserId, tv_path).First(progress).Error
	if err != nil {
		c.JSON(200, gin.H{})
		return
	}
	c.JSON(200, gin.H{"data": progress})
	return
}

func Update(c *gin.Context) {
	db := database.NewDb()
	UserId := c.Request.Header.Get("UserId")
	_TvId := c.Query("tv_id")
	_SeasonId := c.Query("season_id")

	if UserId == "" {
		c.JSON(400, gin.H{"msg": "未获取到所需字段"})
		return
	}
	err := db.Model(&models.User{}).Where("user_id = ?", UserId).First(&models.User{}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(400, gin.H{"msg": "该用户不存在"})
		return
	}
	requestBody := &models.ProgressRequestBody{}

	err = c.ShouldBindJSON(requestBody)
	if err != nil {
		c.JSON(400, gin.H{"msg": "缺少参数"})
		return
	}
	data := models.ProgressRequestBody{}
	data.Data = requestBody.Data

	progress := &models.Progress{}
	dbprogress := &models.Progress{}
	err = db.Model(&models.Progress{}).Where("user_id = ?", UserId).First(progress).Error

	db_data := models.YourModel{}
	_ = json.Unmarshal([]byte(progress.Data), &db_data)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		d, _ := json.Marshal(data.Data)
		progress.Data = string(d)
		progress.UserId = UserId
		progress.CreatedAt = dbprogress.CreatedAt
		err = db.Debug().Model(&models.Progress{}).Create(&progress).Error
		if err != nil {
			c.JSON(500, gin.H{"msg": "记录创建失败"})
			return
		}
		c.JSON(200, gin.H{"msg": "记录创建成功"})
		return
	}
	if _TvId != "" {
		SeasonId, err := strconv.Atoi(_SeasonId)
		TvId, err := strconv.Atoi(_TvId)
		for key, value := range data.Data.ArtPlayerSettings.Times {
			progressTv := &models.ProgressTv{}
			err = db.Model(&models.ProgressTv{}).Where("user_id = ?  and tv_path = ?", UserId, key).First(progressTv).Error
			progressTv.Time = int(value)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				progressTv.UserId = UserId
				progressTv.SeasonId = uint(SeasonId)
				progressTv.TvId = uint(TvId)
				progressTv.TvPath = key
				err = db.Debug().Model(&models.ProgressTv{}).Create(&progressTv).Error
			} else {
				err = db.Model(&models.ProgressTv{}).Where("user_id = ?  and tv_path = ?", UserId, key).Update("time", value).Error
			}

		}
	}

	if db_data.TV == nil {
		db_data.TV = data.Data.TV
	} else {
		for key, value := range data.Data.TV {
			db_data.TV[key] = value
		}
	}
	d, _ := json.Marshal(db_data)
	progress.Data = string(d)

	err = db.Model(&models.Progress{}).Where("id = ?", progress.Id).Select("*").Updates(&progress).Error
	if err != nil {
		c.JSON(500, gin.H{"msg": "记录更新失败"})
		return
	}
	c.JSON(200, gin.H{"msg": "记录更新成功"})
	return
}
