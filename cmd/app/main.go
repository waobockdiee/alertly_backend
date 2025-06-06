package main

import (
	"alertly/internal/activate"
	"alertly/internal/auth"
	"alertly/internal/comments"
	"alertly/internal/database"
	"alertly/internal/editprofile"
	"alertly/internal/getcategories"
	"alertly/internal/getclusterby"
	"alertly/internal/getclustersbylocation"
	"alertly/internal/getincidentsasreels"
	"alertly/internal/getsubcategoriesbycategoryid"
	"alertly/internal/middleware"
	"alertly/internal/myplaces"
	"alertly/internal/newincident"
	"alertly/internal/profile"
	"alertly/internal/saveclusteraccount"
	"alertly/internal/signup"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {

	// ############ PRODUCTION ENV #############/
	// gin.SetMode(gin.ReleaseMode)
	// #########################################/

	err := godotenv.Load()
	if err != nil {
		log.Println("No se pudo cargar el archivo .env, se usarán las variables de entorno del sistema")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	IpServer := os.Getenv("IP_SERVER")
	serverPort := os.Getenv("SERVER_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	database.InitDB(dsn)
	defer database.DB.Close()

	router := gin.Default()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	uploadsPath := filepath.Join(wd, "uploads")
	log.Println("Serving uploads from:", uploadsPath)
	router.Static("/uploads", "./uploads")

	router.POST("/account/signup", signup.RegisterUserHandler)
	router.POST("/account/signin", auth.SignIn)
	router.POST("/account/activate", activate.ActivateAccount)

	api := router.Group("/api")
	api.Use(middleware.TokenAuthMiddleware())
	api.GET("/account/validate", auth.ValidateSession)
	router.GET("/category/get_all", getcategories.GetCategories)
	router.GET("/category/getsubcategoriesbycategoryid/:id", getsubcategoriesbycategoryid.GetSubcategoriesByCategoryId)
	api.POST("/incident/create", newincident.Create)
	api.GET("/cluster/getbyid/:incl_id", getclusterby.View)
	router.GET("/cluster/getbylocation/:min_latitude/:max_latitude/:min_longitude/:max_longitude/:from_date/:to_date/:insu_id", getclustersbylocation.Get)
	api.GET("/cluster/getasreel/:min_latitude/:max_latitude/:min_longitude/:max_longitude", getincidentsasreels.GetReel)

	router.GET("/account/myplaces/get/:account_id", myplaces.Get)
	api.GET("/account/get_my_profile", editprofile.GetMyProfile)
	api.POST("/account/edit_fullname", editprofile.UpdateFullName)
	api.POST("/account/edit/nickname", editprofile.UpdateNickname)
	api.POST("/account/edit/birthdate", editprofile.UpdateBirthDate)
	api.POST("/account/edit/email", editprofile.UpdateEmail)
	api.POST("/account/edit/password", editprofile.UpdatePassword)
	api.POST("/account/check_password", auth.CheckPasswordMatch)
	api.POST("/account/myplaces/add", myplaces.Add)
	api.GET("/account/myplaces/get", myplaces.GetByAccountId)
	api.GET("/account/myplaces/get_by_id/:afl_id", myplaces.GetById)
	api.POST("/account/myplaces/update", myplaces.Update)
	api.POST("/account/myplaces/full_update", myplaces.FullUpdate)
	api.GET("/account/myplaces/delete/:afl_id", myplaces.Delete)
	api.GET("/account/profile/get_by_id/:account_id", profile.GetById)
	api.GET("/account/cluster/toggle_save/:incl_id", saveclusteraccount.ToggleSaveClusterAccount)
	api.POST("/cluster/send_comment", comments.SaveClusterComment)
	api.GET("/saved/get_my_list", saveclusteraccount.GetMyList)
	api.GET("/saved/delete/:acs_id", saveclusteraccount.DeleteFollowIncident)
	api.POST("/account/report/:account_id", profile.ReportAccount)

	log.Printf("Servidor iniciado en :%s", serverPort)

	router.Run(IpServer + ":" + serverPort)
	addr := IpServer + ":" + serverPort
	log.Printf("Servidor iniciado en %s", addr)

	// Configuramos HTTP/2 en modo h2c (sin TLS) para entorno de desarrollo
	h2s := &http2.Server{}
	srv := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(router, h2s),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
