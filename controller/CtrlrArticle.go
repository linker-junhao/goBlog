package controller

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"goBlog/container"
	"goBlog/model"
	"log"
	"net/http"
	"strconv"
)

func ArticleByIdGet(writer http.ResponseWriter, request *http.Request, params httprouter.Params, container container.MyContainer) (res MyHttpHandleResult, resErr error) {
	// 验证参数
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		return res, err
	}

	articleM := model.ArticleModel{Container: container}

	article, err := articleM.GetById(id)

	res.Data = article

	if err != nil {
		return res, err
	}

	return res, nil
}

func ArticleByIdDelete(writer http.ResponseWriter, request *http.Request, params httprouter.Params, container container.MyContainer) (res MyHttpHandleResult, resErr error) {
	// 验证参数
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		return res, err
	}

	articleM := model.ArticleModel{Container: container}
	article, err := articleM.DeleteById(id)

	if err != nil {
		return res, err
	}
	res.Data = article

	return res, nil
}

func ArticlePost(writer http.ResponseWriter, request *http.Request, params httprouter.Params, container container.MyContainer) (res MyHttpHandleResult, resErr error) {
	rc := request.Body
	dataArticle := model.Article{}
	if err := json.NewDecoder(rc).Decode(&dataArticle); err != nil {
		return res, err
	}

	articleM := model.ArticleModel{Container: container}
	dataArticle, err := articleM.NewOne(dataArticle)
	if err != nil {
		return res, err
	}
	res.Data = dataArticle
	return res, nil
}

func ArticlePut(writer http.ResponseWriter, request *http.Request, params httprouter.Params, container container.MyContainer) (res MyHttpHandleResult, resErr error) {
	rc := request.Body
	dataArticle := model.Article{}
	if err := json.NewDecoder(rc).Decode(&dataArticle); err != nil {
		return res, err
	}

	articleM := model.ArticleModel{Container: container}
	dataArticle, err := articleM.ModifyOne(dataArticle)
	if err != nil {
		return res, err
	}
	res.Data = dataArticle
	return res, nil
}

func ArticleListGet(writer http.ResponseWriter, request *http.Request, params httprouter.Params, container container.MyContainer) (res MyHttpHandleResult, resErr error) {
	q := request.URL.Query()
	am := model.ArticleModel{Container: container}
	start, err := strconv.ParseInt(q.Get("start"), 10, 64)
	if err != nil {
		log.Println(err)
		return res, errors.New(getQueryParamsErr)
	}

	offset, err := strconv.ParseInt(q.Get("offset"), 10, 64)
	if err != nil {
		log.Println(err)
		return res, errors.New(getQueryParamsErr)
	}

	list, err := am.List(start, offset)
	if err != nil {
		log.Println(err)
		return res, err
	}

	res.Data = list
	return res, nil
}
