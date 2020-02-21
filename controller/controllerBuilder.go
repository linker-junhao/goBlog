package controller

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"goBlog/container"
	"log"
	"net/http"
)

const (
	MC_OK = 200
	MC_ERROR = 100
)

type httpHandleResult struct {
	Code int
	Data []interface{}
	Msg string
}

type myHttpHandle func(http.ResponseWriter, *http.Request, httprouter.Params, container.MyContainer) (httpHandleResult, error)

type myMiddleware func(writer http.ResponseWriter, request *http.Request, params httprouter.Params, c container.MyContainer, next myHttpHandle) myHttpHandle

type CtrlHandler struct {
	container    container.MyContainer
	handlerFunc  myHttpHandle
	routerHandle httprouter.Handle
	middleware   []myMiddleware
}

func New(c container.MyContainer) *CtrlHandler {
	return &CtrlHandler{container: c}
}

func (ch *CtrlHandler)SetHandlerFunc(hf myHttpHandle) *CtrlHandler {
	ch.handlerFunc = hf
	return ch
}

func (ch *CtrlHandler)AddMiddleware(mw myMiddleware) *CtrlHandler {
	ch.middleware = append(ch.middleware, mw)
	return ch
}

func (ch *CtrlHandler)Handler() httprouter.Handle {
	if ch.routerHandle == nil {
		ch.routerHandle = func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
			var res httpHandleResult
			var err error
			for _, mw := range ch.middleware {
				ch.handlerFunc = mw(writer, request, params, ch.container, ch.handlerFunc)
			}
			res, err = ch.handlerFunc(writer, request, params, ch.container)
			if err != nil {
				res.Code = MC_OK
				res.Msg = err.Error()
			}

			b, jsonErr := json.Marshal(res)
			if jsonErr != nil {
				log.Printf("an json encode error happend: %s", err.Error())
			}
			writer.Header().Set("content-type", "application/json")
			_, err = writer.Write(b)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return ch.routerHandle
}

func HandlerBuild(container container.MyContainer, handler myHttpHandle) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		var res httpHandleResult
		var err error
		res, err = handler(writer, request, params, container)
		if err != nil {
			res.Code = MC_OK
			res.Msg = err.Error()
		}

		b, jsonErr := json.Marshal(res)
		if jsonErr != nil {
			log.Printf("an json encode error happend: %s", err.Error())
		}
		writer.Header().Set("content-type", "application/json")
		_, err = writer.Write(b)
		if err != nil {
			log.Println(err)
		}
	}
}