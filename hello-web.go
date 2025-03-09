package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	testAsciiJSON()
	// testHTML1()
	// testHTML2()
	// testHTML3()
	// testHTTP2ServerPush()
	// testJSONP()
	// testMultipartBind()
	// testMultipartForm()
	// testPureJSON()
	// testQueryAndPostForm()
	// testSecureJSON()
	// testXmlJSONYAMLPROTOBUF() TODO
	// testSecurityHeader()
	// testBindHTMLCheckBox() TODO
	// testBindUrl()
	// testBindFormDataToCustomStruct()
	// testBindQueryStrOrFormData()
	// testQueryStringParam()
	// testReadDataFromReader()
	// testCustomRouteLogFormat()
	// testBingReqBodyToDifferentStruct()
	// testStaticFileServer()
	// testRouteParam()
	// testRouteGroup()
	// testUploadSingleFile()
	// testUploadMulFile()
	// testCookie()
	// testUseBasicAuth()
	// testMapStringOrFormParam()
	// testAutoCertManager()
	// testOnlyBindUrlQryStr()
	// testRedirect()
	// testCustomHttpConfig()
	// testCustomLog()
	// 未通过 testCustomValidator()

}

// func testCustomValidator() {
// 	// curl "localhost:8085/bookable?check_in=2018-04-16&check_out=2018-04-17"
// 	// {"message":"Booking dates are valid!"}
// 	// curl "localhost:8085/bookable?check_in=2018-03-08&check_out=2018-03-09"
// 	// {"error":"Key: 'Booking.CheckIn' Error:Field validation for 'CheckIn' failed on the 'bookabledate' tag"}
// 	type Booking struct {
// 		CheckIn  time.Time `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
// 		CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn,bookabledate" time_format:"2008-01-02"`
// 	}

// 	var bookableDate validator.Func = func(fl validator.FieldLevel) bool {
// 		date, ok := fl.Field().Interface().(time.Time)
// 		if ok {
// 			today := time.Now()
// 			if today.After(date) {
// 				return false
// 			}
// 		}
// 		return true
// 	}

// 	router := gin.Default()

// 	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
// 		v.RegisterValidation("bookabledate", bookableDate)
// 	}

// 	router.GET("/bookable", func(ctx *gin.Context) {
// 		var b Booking
// 		if err := ctx.ShouldBindWith(&b, binding.Query); err == nil {
// 			ctx.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
// 		} else {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		}
// 	})

// 	router.Run(":8080")

// }

// func testCustomLog() {
// 	router := gin.New()
// 	// LoggerWithFormatter 中间件会写入日志到 gin.DefaultWriter
// 	// 默认 gin.DefaultWriter = os.Stdout
// 	router.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
// 		// 你的自定义格式
// 		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
// 			params.ClientIP,
// 			params.TimeStamp.Format(time.RFC1123),
// 			params.Method,
// 			params.Path,
// 			params.Request.Proto,
// 			params.StatusCode,
// 			params.Latency,
// 			params.Request.UserAgent(),
// 			params.ErrorMessage,
// 		)
// 	}))
// 	router.Use(gin.Recovery())
// 	router.GET("/ping", func(ctx *gin.Context) {
// 		ctx.String(200, "pong")
// 	})

// 	router.Run(":8080")
// }

// func testCustomHttpConfig() {
// 	router := gin.Default()

// 	s := &http.Server{
// 		Addr:           ":8888",
// 		Handler:        router,
// 		ReadTimeout:    10 * time.Second,
// 		WriteTimeout:   10 * time.Second,
// 		MaxHeaderBytes: 1 << 20,
// 	}
// 	s.ListenAndServe()
// }

// func testRedirect() {
// 	r := gin.Default()

// 	r.GET("/base", func(ctx *gin.Context) {
// 		ctx.JSON(http.StatusOK, gin.H{"hello": "world"})
// 	})

// 	// 外部重定向
// 	r.GET("/baidu", func(ctx *gin.Context) {
// 		ctx.Redirect(http.StatusMovedPermanently, "http://www.baidu.com")
// 	})

// 	// 内部重定向
// 	r.POST("/hello", func(ctx *gin.Context) {
// 		ctx.Redirect(http.StatusFound, "/base")
// 	})

// 	// 路由重定向
// 	r.GET("/router", func(ctx *gin.Context) {
// 		ctx.Request.URL.Path = "/base"
// 		r.HandleContext(ctx)
// 	})

// 	r.Run(":8080")
// }

// func testOnlyBindUrlQryStr() {
// 	type Person struct {
// 		Name    string `form:"name"`
// 		Address string `form:"address"`
// 	}

// 	route := gin.Default()
// 	route.Any("testing", func(ctx *gin.Context) {
// 		var person Person
// 		if ctx.ShouldBindQuery(&person) == nil {
// 			log.Println("==============only Bind By Query String==============")
// 			log.Println(person.Name)
// 			log.Println(person.Address)
// 		}
// 		ctx.String(200, "Success")
// 	})
// 	route.Run(":8080")
// }

// func testAutoCertManager() {
// 	r := gin.Default()
// 	// Ping handler
// 	r.GET("/ping", func(ctx *gin.Context) {
// 		ctx.String(200, "pong")
// 	})

// 	m := autocert.Manager{
// 		Prompt:     autocert.AcceptTOS,
// 		HostPolicy: autocert.HostWhitelist("example1.com", "example2.com"),
// 		Cache:      autocert.DirCache("/var/www/.cache"),
// 	}

// 	log.Fatal(autotls.RunWithManager(r, &m))
// }

// func testMapStringOrFormParam() {
// 	router := gin.Default()

// 	router.POST("/post", func(ctx *gin.Context) {
// 		ids := ctx.QueryMap("ids")
// 		names := ctx.PostFormMap("names")

// 		fmt.Printf("ids: %v; names: %v", ids, names)
// 	})

// 	router.Run(":8080")
// }

// func testUseBasicAuth() {
// 	var secrets = gin.H{"foo": gin.H{"email": "foo@bar.com", "phone": "123456"},
// 		"austin": gin.H{"email": "austin@example.com", "phone": "666"},
// 		"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
// 	}

// 	r := gin.Default()

// 	// 路由使用 gin.BasicAuth() 中间件
// 	// gin.Accounts 是map[string]string的一种快捷方式

// 	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
// 		"foo":    "bar",
// 		"austin": "1234",
// 		"lena":   "hello2",
// 		"manu":   "4321",
// 	}))

// 	// /admin/secrets 端点
// 	// 触发"localhost:8080/admin/secrets"
// 	authorized.GET("/secrets", func(ctx *gin.Context) {
// 		// 获取用户,它是由 BasicAuth 中间件设置的
// 		user := ctx.MustGet(gin.AuthUserKey).(string)
// 		if secret, ok := secrets[user]; ok {
// 			ctx.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
// 		} else {
// 			ctx.JSON(http.StatusOK, gin.H{"user": user, "secret": "No Secret :("})
// 		}
// 	})
// 	r.Run(":8080")
// }

// func testCookie() {
// 	router := gin.Default()

// 	router.GET("/cookie", func(ctx *gin.Context) {
// 		cookie, err := ctx.Cookie("gin_cookie")

// 		if err != nil {
// 			cookie = "NotSet"
// 			ctx.SetCookie("gin_cookie", "test", 3000, "/", "localhost", false, true)
// 		}
// 		fmt.Printf("Cookie value: %s\n", cookie)
// 	})
// 	router.Run()
// }

// func testUploadMulFile() {
// 	router := gin.Default()
// 	// 为multi form 设置较低的内存限制, 默认是32MIB
// 	router.MaxMultipartMemory = 8 << 20 // 8 MIB
// 	router.POST("/upload", func(ctx *gin.Context) {
// 		// martipart form
// 		form, _ := ctx.MultipartForm()
// 		files := form.File["upload[]"]

// 		for _, file := range files {
// 			log.Println(file.Filename)

// 			// 上传文件到制定目录
// 			dis := file.Filename
// 			ctx.SaveUploadedFile(file, dis)
// 		}
// 		ctx.String(http.StatusOK, fmt.Sprintf("%d files uploaded", len(files)))
// 	})
// 	router.Run(":8080")
// }

// func testUploadSingleFile() {
//	router := gin.Default()
//	// 为mulpart form设置较低的内存设置默认是32 Mib
//	router.MaxMultipartMemory = 8 << 20
//	router.POST("/upload", func(ctx *gin.Context) {
//		file, _ := ctx.FormFile("file")
//		log.Println(file.Filename)
//
//		dst := "./" + file.Filename
//		ctx.SaveUploadedFile(file, dst)
//		ctx.String(http.StatusOK, fmt.Sprintf("'%s' upload!", file.Filename))
//	})
//	router.Run(":8080")
//}

// 	v1 := router.Group("/v1")
// 	{
// 		v1.GET("/login", func(ctx *gin.Context) {
// 			ctx.String(http.StatusOK, "V1_LOGIN")
// 		})
// 		v1.GET("/submit", func(ctx *gin.Context) {
// 			ctx.String(http.StatusOK, "V1_SUBMIT")
// 		})
// 		v1.GET("/read", func(ctx *gin.Context) {
// 			ctx.String(http.StatusOK, "V1_READ")
// 		})
// 	}

// 	// 简单的路由组: v2
// 	v2 := router.Group("/v2")
// 	{
// 		v2.GET("/login", func(ctx *gin.Context) {
// 			ctx.String(http.StatusOK, "V2_LOGIN")
// 		})
// 		v2.GET("/submit", func(ctx *gin.Context) {
// 			ctx.String(http.StatusOK, "V2_SUBMIT")
// 		})
// 		v2.GET("/read", func(ctx *gin.Context) {
// 			ctx.String(http.StatusOK, "V2_READ")
// 		})
// 	}

// 	router.Run(":8080")
// }

// func testRouteParam() {
// 	router := gin.Default()

// 	// 此 handler 将匹配 /user/john 但不会匹配 /user/ 或者 /user
// 	router.GET("/user/:name", func(c *gin.Context) {
// 		name := c.Param("name")
// 		c.String(http.StatusOK, "Hello %s", name)
// 	})

// 	// 此 handler 将匹配 /user/john/ 和 /user/john/send
// 	// 如果没有其他路由匹配 /user/john，它将重定向到 /user/john/
// 	router.GET("/user/:name/*action", func(c *gin.Context) {
// 		name := c.Param("name")
// 		action := c.Param("action")
// 		message := name + " is " + action
// 		c.String(http.StatusOK, message)
// 	})

// 	router.Run(":8080")
// }

// func testStaticFileServer() {
// 	router := gin.Default()
// 	router.Static("/templates", "./templates")
// 	router.StaticFS("/more_static", http.Dir("/opt"))
// 	router.StaticFile("/testdata", "./testdata")

// 	// 监听并在 0.0.0.0:8080 上启动服务
// 	router.Run(":8080")
// }

// func testBingReqBodyToDifferentStruct() {
// 	r := gin.Default()

// 	type formA struct {
// 		Foo string `json:"foo" xml:"foo" binding:"required"`
// 	}

// 	type formB struct {
// 		Bar string `json:"bar" xml:"bar" binding:"required"`
// 	}

// 	r.POST("/oneBind", func(c *gin.Context) {
// 		objA := formA{}
// 		objB := formB{}
// 		// c.ShouldBind 使用了 c.Request.Body，不可重用。
// 		if errA := c.ShouldBind(&objA); errA == nil {
// 			c.String(http.StatusOK, `the body should be formA`)
// 			// 因为现在 c.Request.Body 是 EOF，所以这里会报错。
// 		} else if errB := c.ShouldBind(&objB); errB == nil {
// 			c.String(http.StatusOK, `the body should be formB`)
// 		} else {
// 		}
// 	})

// 	r.GET("/mulBind", func(c *gin.Context) {
// 		objA := formA{}
// 		objB := formB{}
// 		// 读取 c.Request.Body 并将结果存入上下文。
// 		if errA := c.ShouldBindBodyWith(&objA, binding.JSON); errA == nil {
// 			c.String(http.StatusOK, `the body should be formA`)
// 			// 这时, 复用存储在上下文中的 body。
// 		} else if errB := c.ShouldBindBodyWith(&objB, binding.JSON); errB == nil {
// 			c.String(http.StatusOK, `the body should be formB JSON`)
// 			// 可以接受其他格式
// 		} else if errB2 := c.ShouldBindBodyWith(&objB, binding.XML); errB2 == nil {
// 			c.String(http.StatusOK, `the body should be formB XML`)
// 		} else {
// 		}
// 	})

// 	// 监听并在 0.0.0.0:8080 上启动服务
// 	r.Run()
// }

// func testCustomRouteLogFormat() {
// 	r := gin.Default()
// 	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
// 		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
// 	}

// 	r.POST("/foo", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, "foo")
// 	})

// 	r.GET("/bar", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, "bar")
// 	})

// 	r.GET("/status", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, "ok")
// 	})

// 	// 监听并在 0.0.0.0:8080 上启动服务
// 	r.Run()
// }

// func testReadDataFromReader() {
// 	router := gin.Default()
// 	router.GET("/someDataFromReader", func(c *gin.Context) {
// 		response, err := http.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png")
// 		if err != nil || response.StatusCode != http.StatusOK {
// 			c.Status(http.StatusServiceUnavailable)
// 			return
// 		}

// 		reader := response.Body
// 		contentLength := response.ContentLength
// 		contentType := response.Header.Get("Content-Type")

// 		extraHeaders := map[string]string{
// 			"Content-Disposition": `attachment; filename="gopher.png"`,
// 		}

// 		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
// 	})
// 	router.Run(":8080")
// }

// func testQueryStringParam() {
// 	router := gin.Default()

// 	// 使用现有的基础请求对象解析查询字符串参数。
// 	// 示例 URL： /welcome?firstname=Jane&lastname=Doe
// 	router.GET("/welcome", func(c *gin.Context) {
// 		firstname := c.DefaultQuery("firstname", "Guest")
// 		lastname := c.Query("lastname") // c.Request.URL.Query().Get("lastname") 的一种快捷方式

// 		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
// 	})
// 	router.Run(":8080")
// }

// func testBindQueryStrOrFormData() {
// 	// $ curl -X GET "localhost:8085/testing?name=appleboy&address=xyz&birthday=1992-03-15"
// 	type Person struct {
// 		Name     string    `form:"name"`
// 		Address  string    `form:"address"`
// 		Birthday time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
// 	}

// 	route := gin.Default()
// 	route.GET("/testing", func(c *gin.Context) {
// 		var person Person
// 		// 如果是 `GET` 请求，只使用 `Form` 绑定引擎（`query`）。
// 		// 如果是 `POST` 请求，首先检查 `content-type` 是否为 `JSON` 或 `XML`，然后再使用 `Form`（`form-data`）。
// 		// 查看更多：https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L88
// 		if c.ShouldBind(&person) == nil {
// 			log.Println(person.Name)
// 			log.Println(person.Address)
// 			log.Println(person.Birthday)
// 		}

// 		c.String(200, "Success")
// 	})
// 	route.Run(":8085")
// }

// func testBindFormDataToCustomStruct() {
// 	// $ curl "http://localhost:8080/getb?field_a=hello&field_b=world"
// 	// {"a":{"FieldA":"hello"},"b":"world"}
// 	// $ curl "http://localhost:8080/getc?field_a=hello&field_c=world"
// 	// {"a":{"FieldA":"hello"},"c":"world"}
// 	// $ curl "http://localhost:8080/getd?field_x=hello&field_d=world"
// 	// {"d":"world","x":{"FieldX":"hello"}}
// 	r := gin.Default()
// 	r.GET("/getb", GetDataB)
// 	r.GET("/getc", GetDataC)
// 	r.GET("/getd", GetDataD)

// 	r.Run()
// }

// type StructA struct {
// 	FieldA string `form:"field_a"`
// }

// type StructB struct {
// 	NestedStruct StructA
// 	FieldB       string `form:"field_b"`
// }

// type StructC struct {
// 	NestedStructPointer *StructA
// 	FieldC              string `form:"field_c"`
// }

// type StructD struct {
// 	NestedAnonyStruct struct {
// 		FieldX string `form:"field_x"`
// 	}
// 	FieldD string `form:"field_d"`
// }

// func GetDataB(c *gin.Context) {
// 	var b StructB
// 	c.Bind(&b)
// 	c.JSON(200, gin.H{
// 		"a": b.NestedStruct,
// 		"b": b.FieldB,
// 	})
// }

// func GetDataC(c *gin.Context) {
// 	var b StructC
// 	c.Bind(&b)
// 	c.JSON(200, gin.H{
// 		"a": b.NestedStructPointer,
// 		"c": b.FieldC,
// 	})
// }

// func GetDataD(c *gin.Context) {
// 	var b StructD
// 	c.Bind(&b)
// 	c.JSON(200, gin.H{
// 		"x": b.NestedAnonyStruct,
// 		"d": b.FieldD,
// 	})
// }

// func testBindUrl() {
// 	// $ curl -v localhost:8088/thinkerou/987fbc97-4bed-5078-9f07-9141ba07c9f3
// 	// $ curl -v localhost:8088/thinkerou/not-uuid
// 	type Person struct {
// 		ID   string `uri:"id" binding:"required,uuid"`
// 		Name string `uri:"name" binding:"required"`
// 	}

// 	route := gin.Default()
// 	route.GET("/:name/:id", func(c *gin.Context) {
// 		var person Person
// 		if err := c.ShouldBindUri(&person); err != nil {
// 			c.JSON(400, gin.H{"msg": err.Error()})
// 			return
// 		}
// 		c.JSON(200, gin.H{"name": person.Name, "uuid": person.ID})
// 	})
// 	route.Run(":8088")
// }

// func testBindHTMLCheckBox() {
// 	// https://gin-gonic.com/zh-cn/docs/examples/bind-html-checkbox/
// }

// func testSecurityHeader() {
// 	// 检查页眉
// 	// curl localhost:8080/ping -I
// 	// 检查主机标头注入
// 	// curl localhost:8080/ping -I -H "Host:neti.ee"
// 	r := gin.Default()

// 	expectedHost := "localhost:8080"

// 	// Setup Security Headers
// 	r.Use(func(c *gin.Context) {
// 		if c.Request.Host != expectedHost {
// 			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
// 			return
// 		}
// 		c.Header("X-Frame-Options", "DENY")
// 		c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
// 		c.Header("X-XSS-Protection", "1; mode=block")
// 		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
// 		c.Header("Referrer-Policy", "strict-origin")
// 		c.Header("X-Content-Type-Options", "nosniff")
// 		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
// 		c.Next()
// 	})

// 	r.GET("/ping", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"message": "pong",
// 		})
// 	})

// 	r.Run() // listen and serve on 0.0.0.0:8080
// }

// func testXmlJSONYAMLPROTOBUF() {
// 	r := gin.Default()

// 	// gin.H 是 map[string]interface{} 的一种快捷方式
// 	r.GET("/someJSON", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
// 	})

// 	r.GET("/moreJSON", func(c *gin.Context) {
// 		// 你也可以使用一个结构体
// 		var msg struct {
// 			Name    string `json:"user"`
// 			Message string
// 			Number  int
// 		}
// 		msg.Name = "Lena"
// 		msg.Message = "hey"
// 		msg.Number = 123
// 		// 注意 msg.Name 在 JSON 中变成了 "user"
// 		// 将输出：{"user": "Lena", "Message": "hey", "Number": 123}
// 		c.JSON(http.StatusOK, msg)
// 	})

// 	r.GET("/someXML", func(c *gin.Context) {
// 		c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
// 	})

// 	r.GET("/someYAML", func(c *gin.Context) {
// 		c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
// 	})

// 	r.GET("/someProtoBuf", func(c *gin.Context) {
// 		reps := []int64{int64(1), int64(2)}
// 		label := "test"
// 		// protobuf 的具体定义写在 testdata/protoexample 文件中。
// 		data := &protoexample.Test{
// 			Label: &label,
// 			Reps:  reps,
// 		}
// 		// 请注意，数据在响应中变为二进制数据
// 		// 将输出被 protoexample.Test protobuf 序列化了的数据
// 		c.ProtoBuf(http.StatusOK, data)
// 	})

// 	// 监听并在 0.0.0.0:8080 上启动服务
// 	r.Run(":8080")
// }

// func testSecureJSON() {
// 	r := gin.Default()

// 	// 你也可以使用自己的 SecureJSON 前缀
// 	// r.SecureJsonPrefix(")]}',\n")

// 	r.GET("/someJSON", func(c *gin.Context) {
// 		names := []string{"lena", "austin", "foo"}

// 		// 将输出：while(1);["lena","austin","foo"]
// 		c.SecureJSON(http.StatusOK, names)
// 	})

// 	// 监听并在 0.0.0.0:8080 上启动服务
// 	r.Run(":8080")
// }

// func testQueryAndPostForm() {
// 	router := gin.Default()

// 	router.POST("/post", func(c *gin.Context) {

// 		id := c.Query("id")
// 		page := c.DefaultQuery("page", "0")
// 		name := c.PostForm("name")
// 		message := c.PostForm("message")

// 		fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
// 	})
// 	router.Run(":8080")
// }

// func testPureJSON() {
// 	r := gin.Default()

// 	// 提供 unicode 实体
// 	r.GET("/json", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"html": "<b>Hello, world!</b>",
// 		})
// 	})

// 	// 提供字面字符
// 	r.GET("/purejson", func(c *gin.Context) {
// 		c.PureJSON(200, gin.H{
// 			"html": "<b>Hello, world!</b>",
// 		})
// 	})

// 	// 监听并在 0.0.0.0:8080 上启动服务
// 	r.Run(":8080")
// }

// func testMultipartForm() {
// 	router := gin.Default()

// 	router.POST("/form_post", func(c *gin.Context) {
// 		message := c.PostForm("message")
// 		nick := c.DefaultPostForm("nick", "anonymous")

// 		c.JSON(200, gin.H{
// 			"status":  "posted",
// 			"message": message,
// 			"nick":    nick,
// 		})
// 	})
// 	router.Run(":8080")
// }

// func testMultipartBind() {
// 	// test: $ curl -v --form user=user --form password=password http://localhost:8080/login
// 	type LoginForm struct {
// 		User     string `form:"user" binding:"required"`
// 		Password string `form:"password" binding:"required"`
// 	}
// 	router := gin.Default()
// 	router.POST("/login", func(c *gin.Context) {
// 		// 你可以使用显式绑定声明绑定 multipart form：
// 		// c.ShouldBindWith(&form, binding.Form)
// 		// 或者简单地使用 ShouldBind 方法自动绑定：
// 		var form LoginForm
// 		// 在这种情况下，将自动选择合适的绑定
// 		if c.ShouldBind(&form) == nil {
// 			if form.User == "user" && form.Password == "password" {
// 				c.JSON(200, gin.H{"status": "you are logged in"})
// 			} else {
// 				c.JSON(401, gin.H{"status": "unauthorized"})
// 			}
// 		}
// 	})
// 	router.Run(":8080")
// }

// func testJSONP() {
// 	r := gin.Default()

// 	r.GET("/JSONP", func(c *gin.Context) {
// 		data := map[string]interface{}{
// 			"foo": "bar",
// 		}

// 		// /JSONP?callback=x
// 		// 将输出：x({\"foo\":\"bar\"})
// 		c.JSONP(http.StatusOK, data)
// 	})

// 	// 监听并在 0.0.0.0:8080 上启动服务
// 	r.Run(":8080")
// }

// func testHTTP2ServerPush() {
// 	var html = template.Must(template.New("https").Parse(`
// 	<html>
// 	<head>
// 	<title>Https Test</title>
// 	<script src="/assets/app.js"></script>
// 	</head>
// 	<body>
// 	<h1 style="color:red;">Welcome, Ginner!</h1>
// 	</body>
// 	</html>
// 	`))

// 	r := gin.Default()
// 	r.Static("/assets", "./assets")
// 	r.SetHTMLTemplate(html)

// 	r.GET("/", func(c *gin.Context) {
// 		if pusher := c.Writer.Pusher(); pusher != nil {
// 			// 使用 pusher.Push() 做服务器推送
// 			if err := pusher.Push("/assets/app.js", nil); err != nil {
// 				log.Printf("Failed to push: %v", err)
// 			}
// 		}
// 		c.HTML(200, "https", gin.H{
// 			"status": "success",
// 		})
// 	})

// 	// 监听并在 https://127.0.0.1:8080 上启动服务
// 	r.RunTLS(":8080", "./testdata/server.pem", "./testdata/server.key")
// }

// func testHTML3() {
// 	router := gin.Default()
// 	router.Delims("{[{", "}]}")
// 	router.SetFuncMap(template.FuncMap{
// 		"formatAsDate": formatAsDate,
// 	})
// 	router.LoadHTMLFiles("./testdata/template/raw.tmpl")

// 	router.GET("/raw", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "raw.tmpl", map[string]interface{}{
// 			"now": time.Date(2017, 07, 01, 0, 0, 0, 0, time.UTC),
// 		})
// 	})

// 	router.Run(":8080")
// }

// func formatAsDate(t time.Time) string {
// 	year, month, day := t.Date()
// 	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
// }

// func testHTML2() {
// 	router := gin.Default()
// 	router.LoadHTMLGlob("templates/**/*")
// 	router.GET("/posts/index", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
// 			"title": "Posts",
// 		})
// 	})
// 	router.GET("/users/index", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
// 			"title": "Users",
// 		})
// 	})
// 	router.Run(":8080")
// }

// func testHTML1() {
// 	router := gin.Default()
// 	router.LoadHTMLGlob("templates/*")
// 	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
// 	router.GET("/index", func(c *gin.Context) {
// 		c.HTML(http.StatusOK, "index.tmpl", gin.H{
// 			"title": "Main website",
// 		})
// 	})
// 	router.Run(":8080")
// }

func testAsciiJSON() {
	r := gin.Default()

	r.GET("/someJSON", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO语言",
			"tag":  "<br>",
		}

		// 输出 : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
		c.AsciiJSON(http.StatusOK, data)
	})

	// 监听并在 0.0.0.0:8080 上启动服务
	r.Run(":8080")
}
