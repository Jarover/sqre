package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Jarover/sqre/models"
	"github.com/Jarover/sqre/readconfig"
	"github.com/Jarover/sqre/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

func FloatToString(inputNum float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}

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
		Link_id:           link.ID,
		Created:           t.Format("2006-01-02 15:04:05"),
		Referrer:          "",
		User_agent:        "+",
		User_agent_source: "",
	}
	result := models.GetDB().Create(&click)
	if result.Error != nil {
		log.Println(result.Error)
	}
	if len(link.Remote) > 0 {
		c.Redirect(http.StatusMovedPermanently, link.Remote)
		c.Abort()
	} else {

		switch link.Cat_id {
		// Object
		case 4:
			var obj = models.Gobject{}
			err := models.GetDB().First(&obj, "id = ?", link.Objid).Error
			if err != nil {
				log.Println(err)
				log.Println("bad object cat id : " + strconv.Itoa(link.ID))
				c.JSON(http.StatusOK, gin.H{

					"errorGobject": par,
				})
			}
			url := "https://yandex.ru/maps/?pt=" + FloatToString(obj.Lon) + "," + FloatToString(obj.Lat) + "&z=14&l=map"
			log.Println(url)
			c.Redirect(http.StatusMovedPermanently, url)
			c.Abort()
			break
		default:
			c.JSON(http.StatusOK, gin.H{
				"version": readconfig.Version.VersionStr(),
				"url":     par,
				"typeId":  link.Cat_id,
				"obj":     link.Objid,
			})
		}

	}

}

func info(c *gin.Context) {

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

	out := gin.H{
		"version": readconfig.Version.VersionStr(),
		"url":     par,
		"typeId":  link.Cat_id,
	}

	switch link.Cat_id {
	// Object
	case 4:
		var obj = models.Gobject{}
		err := models.GetDB().First(&obj, "id = ?", link.Objid).Error
		if err != nil {
			log.Println(err)
			log.Println("bad object cat id : " + strconv.Itoa(link.ID))

		}
		out["obj"] = obj
		if len(link.Remote) > 0 {
			out["url"] = link.Remote
		} else {
			out["url"] = "https://yandex.ru/maps/?pt=" + FloatToString(obj.Lon) + "," + FloatToString(obj.Lat) + "&z=14&l=map"

		}
		break
	}

	c.JSON(http.StatusOK, out)
}

func sqlTask(c *gin.Context) {

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
		"version": readconfig.Version.VersionStr(),
		"data":    readconfig.Version.BuildTime,
	})
}

// Читаем флаги и окружение
func readFlag(configFlag *readconfig.Flag) {
	flag.StringVar(&configFlag.ConfigFile, "f", readconfig.GetEnv("CONFIGFILE", readconfig.GetDefaultConfigFile()), "config file")
	//flag.StringVar(&configFlag.Host, "h", readconfig.GetEnv("HOST", ""), "host")
	flag.UintVar(&configFlag.Port, "p", uint(readconfig.GetEnvInt("PORT", 0)), "port")
	flag.Parse()

}

func main() {

	dir := utils.GetDir()
	err := readconfig.Version.ReadVersionFile(dir + "/version.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(readconfig.Version)
	var configFlag readconfig.Flag
	readFlag(&configFlag)

	fmt.Println(configFlag)
	fmt.Println(dir + "/" + configFlag.ConfigFile)
	Config, err := readconfig.ReadConfig(dir + "/" + configFlag.ConfigFile)
	if configFlag.Port != 0 {
		fmt.Println(Config)
		Config.Port = configFlag.Port
	}

	if err != nil {
		panic(err)
	}

	fmt.Println(Config)

	logPath := dir + "/sqre_app.log"

	fmt.Println(logPath)
	l := &lumberjack.Logger{ //nolint:typecheck
		Filename:   logPath,
		MaxSize:    500, // megabytes
		MaxBackups: 10,
		MaxAge:     1,     //days
		Compress:   false, // disabled by default
	}
	log.SetOutput(l)

	r := gin.Default()
	r.GET("/", startPage)

	r.GET("/:par", redirect)
	r.GET("/:par/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/:par/info", info)
	r.GET("/:par/sqltask", sqlTask)

	r.Run(":" + strconv.FormatUint(uint64(Config.Port), 10)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
