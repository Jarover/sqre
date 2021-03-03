package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"sqre/models"
	"sqre/utils"
	"time"
)
var (
	version = "0.0.3"
)

func redirect(c *gin.Context) {
	par := c.Param("par")
	var link = models.LinkTrek{}

	err := models.GetDB().First(&link, "Short = ?", par).Error
	if err != nil {
		log.Println(err)
		log.Println("bad url : " + par)
		c.JSON(http.StatusOK, gin.H{

			"error": par,
		})
	}
	log.Println("redirect to : " + link.Remote)
	t := time.Now()
	t.Format("2006-01-02 15:04:05")
	var click = models.Click{
		Link_id: link.ID,
		Created: t.Format("2006-01-02 15:04:05"),
		Referrer: "",
		User_agent: "+",
		User_agent_source: "",
	}
	result:= models.GetDB().Create(&click)
	if result.Error != nil {
		log.Println(result.Error)
	}
	c.Redirect(http.StatusMovedPermanently, link.Remote)
	c.Abort()

}

func sqltask(c *gin.Context) {

	var link = models.LinkTrek{}
	var link2 = models.LinkTrek{}
	models.GetDB().First(&link, "Short = ?", "V35")

	err := models.GetDB().Table("linktrek_link").Where("id = ?", 33).Find(&link2).Error
	if err != nil {
		log.Println(err)

	}

	c.JSON(http.StatusOK, gin.H{
		"ok":      "ok",
		"remote":  link.Remote,
		"remote2": link2.Remote,
	})
}

func startPage(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"version": version,
	})
}



func main() {

	dir := utils.GetDir()

	logpath := dir + "/sqre_app.log"

	fmt.Println(logpath)
	l := &lumberjack.Logger{
		Filename:   logpath,
		MaxSize:    500, // megabytes
		MaxBackups: 10,
		MaxAge:     1,     //days
		Compress:   false, // disabled by default
	}
	log.SetOutput(l)

	log.Println("Version : " + version)
	r := gin.Default()
	r.GET("/", startPage)

	r.GET("/:par", redirect)
	r.GET("/:par/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/:par/sqltask", sqltask)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
