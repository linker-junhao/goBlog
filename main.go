package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"goBlog/container"
	ctlr "goBlog/controller"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile)

	// create mysql con
	db, err := sql.Open("mysql",
		"root:root12345@tcp(127.0.0.1:32768)/myblog")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// container
	c := container.MyContainer{MysqlDB: db}

	router := httprouter.New()
	articleByIdGetHandle := ctlr.New(c).SetHandlerFunc(ctlr.ArticleByIdGet).Handler()
	articleByIdDeleteHandle := ctlr.New(c).SetHandlerFunc(ctlr.ArticleByIdDelete).Handler()
	articlePostHandle := ctlr.New(c).SetHandlerFunc(ctlr.ArticlePost).Handler()
	articlePutHandle := ctlr.New(c).SetHandlerFunc(ctlr.ArticlePut).Handler()
	articleListGetHandle := ctlr.New(c).SetHandlerFunc(ctlr.ArticleListGet).Handler()

	router.GET("/article/:id", articleByIdGetHandle)
	router.DELETE("/article/:id", articleByIdDeleteHandle)
	router.POST("/article", articlePostHandle)
	router.PUT("/article", articlePutHandle)
	router.GET("/article", articleListGetHandle)

	http.ListenAndServe(":8080", router)
}
