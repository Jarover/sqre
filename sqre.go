package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"encoding/json"

	"github.com/Jarover/sqre/models"
	"github.com/Jarover/sqre/readconfig"
	"github.com/Jarover/sqre/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

func redirect(c *gin.Context) {
	par := c.Param("par")

	switch par {
	case "news":

		out := news(par)
		c.JSON(http.StatusOK, out)
		c.Abort()
		return

	}

	var link = models.LinkTrek{}

	err := models.GetDB().First(&link, "Short = ?", par).Error
	if err != nil {
		log.Println(err)
		log.Println("bad url : " + par)
		c.JSON(http.StatusOK, gin.H{

			"error": par,
		})
		return
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
	if len(link.Remote) > 5 {
		c.Redirect(http.StatusMovedPermanently, link.Remote)
		c.Abort()
	} else {

		switch link.Cat_id {
		//contact
		case 3:
			c.HTML(http.StatusOK, "contact.html", gin.H{
				"title":  "Контакт",
				"anonce": link.Text,
				"name":   link.Name,
			})
			break
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
			url := "https://yandex.ru/maps/?pt=" + utils.FloatToString(obj.Lon) + "," + utils.FloatToString(obj.Lat) + "&z=14&l=map"
			log.Println(url)
			c.Redirect(http.StatusMovedPermanently, url)
			c.Abort()
			break
		default:
			c.JSON(http.StatusOK, gin.H{
				"version": readconfig.Version.VersionStr(),
				"url":     par,
				"typeid":  link.Cat_id,
				"obj":     link.Objid,
			})
		}

	}

}

func suffix(c *gin.Context) {
	par := c.Param("par")
	suf := c.Param("suf")
	out := gin.H{
		"version": readconfig.Version.VersionStr(),
		"short":   par,
	}
	switch suf {
	case "info":
		out = info(par)

	case "u":
		fallthrough
	case "m":
		var link = models.LinkTrek{}

		err := models.GetDB().First(&link, "Short = ?", par).Error

		if err != nil {
			log.Println(err)
			log.Println("bad url : " + par)
			out["error"] = par

		}
		switch link.Cat_id {

		// Object
		case 4:
			var obj = models.Gobject{}
			err := models.GetDB().First(&obj, "id = ?", link.Objid).Error
			if err != nil {
				log.Println(err)
				log.Println("bad object cat id : " + strconv.Itoa(link.ID))
				out["error"] = "bad object cat id : " + strconv.Itoa(link.ID)
				break
			}
			attribute := models.Attribute{}
			bytes := []byte(obj.Attributes)
			json.Unmarshal(bytes, &attribute)
			var urls []models.FieldRow
			var redirectM, redirectU string
			for _, v := range attribute.Urls {
				if v.Suffix == "m" {
					redirectM = v.Name
				}
				if v.Suffix == "u" {
					redirectU = v.Name
				}

				urls = append(urls, models.FieldRow{Name: v.Name, Info: v.Info})
			}
			out["urls"] = urls

			if suf == "m" && len(redirectM) > 0 {
				c.Redirect(http.StatusMovedPermanently, redirectM)
				c.Abort()
			}
			if suf == "u" && len(redirectU) > 0 {
				c.Redirect(http.StatusMovedPermanently, redirectU)
				c.Abort()
			}
		}

	default:
		out = gin.H{
			"suffix":  suf,
			"short":   par,
			"version": readconfig.Version.VersionStr(),
		}
	}

	c.JSON(http.StatusOK, out)
}

func news(par string) gin.H {
	out := gin.H{
		"version": readconfig.Version.VersionStr(),
	}
	var news []models.News
	if err := models.GetDB().Where("published = ?", true).Limit(5).Order("-id").Find(&news).Error; err != nil {
		log.Println(err)
	}
	out["items"] = news
	return out
}

func info(par string) gin.H {

	var link = models.LinkTrek{}

	err := models.GetDB().First(&link, "Short = ?", par).Error
	out := gin.H{
		"version": readconfig.Version.VersionStr(),
		"url":     par,
	}
	if err != nil {
		log.Println(err)
		log.Println("bad url : " + par)
		out["error"] = par
		return out

	}

	out["typeid"] = link.Cat_id

	switch link.Cat_id {
	// Object
	case 2:
		out["anonce"] = link.Text
		out["name"] = link.Name
		out["url"] = link.Remote
		break
	// Contact
	case 3:
		out["anonce"] = link.Text
		out["name"] = link.Name
		out["url"] = ""

		break
	// Object
	case 4:
		var obj = models.Gobject{}
		err := models.GetDB().First(&obj, "id = ?", link.Objid).Error
		if err != nil {
			log.Println(err)
			log.Println("bad object cat id : " + strconv.Itoa(link.ID))

		}
		out["id"] = obj.ID
		out["name"] = obj.Name
		out["anonce"] = obj.Anonce
		out["address"] = obj.Address
		out["desc"] = obj.Desc
		out["catid"] = obj.Cat_id
		out["lat"] = obj.Lat
		out["lon"] = obj.Lon

		if len(link.Remote) > 0 {
			out["url"] = link.Remote
		} else {
			out["url"] = "https://yandex.ru/maps/?pt=" + utils.FloatToString(obj.Lon) + "," + utils.FloatToString(obj.Lat) + "&z=14&l=map"

		}

		var emails, phones, urls, photos, audios, tracks, routes, videos []models.FieldRow
		var images []models.Upload
		if err := models.GetDB().Where("gobject_id = ?", obj.ID).Order("id").Find(&images).Error; err != nil {
			log.Println(err)
		}

		attribute := models.Attribute{}
		bytes := []byte(obj.Attributes)
		json.Unmarshal(bytes, &attribute)
		var firstChar string
		for _, v := range images {
			if len(v.Suffix) > 0 {
				firstChar = v.Suffix[:1]
			} else {
				firstChar = ""
			}
			if firstChar == "i" {
				photos = append(photos, models.FieldRow{Name: "/media/" + v.Ufile, Info: v.Name})
			}
			if firstChar == "a" {
				audios = append(audios, models.FieldRow{Name: "/media/" + v.Ufile, Info: v.Name})
			}
			if firstChar == "r" {
				routes = append(routes, models.FieldRow{Name: "/media/" + v.Ufile, Info: v.Name})
			}
			if firstChar == "t" {
				tracks = append(tracks, models.FieldRow{Name: "/media/" + v.Ufile, Info: v.Name})
			}

		}

		if len(obj.Email) > 0 {
			emails = append(emails, models.FieldRow{Name: obj.Email, Info: "E-mail"})
		}

		for _, v := range attribute.Phones {
			if len(v.Suffix) > 0 {
				firstChar = v.Suffix[:1]
			} else {
				firstChar = ""
			}
			if firstChar == "e" {
				emails = append(emails, models.FieldRow{Name: v.Name, Info: v.Info})
			}
			if firstChar == "p" || v.Suffix == "" {
				phones = append(phones, models.FieldRow{Name: v.Name, Info: v.Info})
			}

		}

		for _, v := range attribute.Urls {
			if len(v.Suffix) > 0 {
				firstChar = v.Suffix[:1]
			} else {
				firstChar = ""
			}
			if firstChar == "u" || firstChar == "" || firstChar == "m" {
				urls = append(urls, models.FieldRow{Name: v.Name, Info: v.Info})
			}
			if firstChar == "r" {
				routes = append(routes, models.FieldRow{Name: v.Name, Info: v.Info})
			}
			if firstChar == "t" {
				tracks = append(tracks, models.FieldRow{Name: v.Name, Info: v.Info})
			}
			if firstChar == "a" {
				audios = append(audios, models.FieldRow{Name: v.Name, Info: v.Info})
			}
			if firstChar == "v" {
				videos = append(videos, models.FieldRow{Name: v.Name, Info: v.Info})
			}
			if v.Suffix == "e" {
				emails = append(emails, models.FieldRow{Name: v.Name, Info: v.Info})
			}

		}
		out["emails"] = emails
		out["phones"] = phones
		out["photos"] = photos
		out["audios"] = audios
		out["tracks"] = tracks
		out["routes"] = routes
		out["videos"] = videos
		out["urls"] = urls
		//out["attributes"] = attribute

		break
	}

	return out
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
	flag.StringVar(&configFlag.ConfigFile, "f", readconfig.GetEnv("CONFIGFILE", utils.GetBaseFile()+".json"), "config file")
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
	r.LoadHTMLGlob(dir + "/templates/*")
	r.GET("/", startPage)

	r.GET("/:par", redirect)
	/*
		r.GET("/:par/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		r.GET("/:par/info", info)
		r.GET("/:par/sqltask", sqlTask)
	*/

	r.GET("/:par/:suf", suffix)

	r.Run(":" + strconv.FormatUint(uint64(Config.Port), 10)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
