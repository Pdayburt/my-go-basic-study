package web

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/gin-gonic/gin"
	"my-go-basic-study/webook/internal/domain"
	"my-go-basic-study/webook/internal/service"
	"net/http"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc     service.ArticleService
	intrSvc service.InteractiveService
}

func NewArticleHandler(svc service.ArticleService, intrSvc service.InteractiveService) *ArticleHandler {
	return &ArticleHandler{
		svc:     svc,
		intrSvc: intrSvc,
	}
}

func (a *ArticleHandler) RegisterRouters(server *gin.Engine) {
	articleGroup := server.Group("/articles")
	articleGroup.POST("/edit", a.Edit)
	articleGroup.POST("/publish", a.Publish)
	articleGroup.POST("/withdraw", a.withdraw)
}

func (a *ArticleHandler) withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	err := a.svc.Withdraw(ctx, domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: 1,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result[int64]{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result[int64]{
		Msg:  "OK",
		Data: req.Id,
	})
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	id, err := a.svc.Save(ctx, domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: 1,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result[int64]{
			Code: 5,
			Msg:  "系统错误",
		})
		log.Info("article 保存失败:", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, Result[int64]{
		Msg:  "OK",
		Data: id,
	})

}

func (a *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	id, err := a.svc.Publish(ctx, domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: 1,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result[int64]{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result[int64]{
		Msg:  "OK",
		Data: id,
	})
}

type ArticleReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
