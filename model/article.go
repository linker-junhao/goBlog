package model

import (
	"errors"
	"goBlog/container"
	"goBlog/mysqlORM"
)

type ArticleModel struct {
	Container container.MyContainer
	orm       *mysqlORM.MysqlORM
	isInit    bool
}

type Article struct {
	Id        int64
	Title     string
	Content   string
	CreatedAt string
	UpdatedAt string
}

func (a *ArticleModel) initRun() {
	if !a.isInit {
		cfg := mysqlORM.NewORMConfig("article", a.Container.MysqlDB)
		orm := mysqlORM.NewMysqlORM(cfg)
		a.orm = &orm
		a.isInit = true
	}
}

func (a *ArticleModel) GetById(id int64) (article Article, err error) {
	a.initRun()
	whereArticle := Article{Id: id}
	ret, err := a.orm.Select().SetSelectWhereFields("and", []string{"id"}).
		Where(whereArticle).Commit()
	if err != nil {
		return article, err
	}

	if ret == nil {
		return article, errors.New("the item you find didnt exists")
	}
	re := ret[0].(Article)

	return re, nil
}

func (a *ArticleModel) DeleteById(id int64) (article Article, err error) {
	a.initRun()
	where := Article{Id: id}
	ret, err := a.orm.Delete().SetDeleteWhereFields("and", []string{"id"}).Where(where).Commit()
	if err != nil {
		return where, err
	}
	return ret.(Article), nil
}

func (a *ArticleModel) NewOne(article Article) (ar Article, err error) {
	a.initRun()
	ret, err := a.orm.Insert().SetInsertFields([]string{"title", "content"}).Commit(article)
	if err != nil {
		return article, err
	}
	return ret.(Article), nil
}

func (a *ArticleModel) ModifyOne(article Article) (ar Article, err error) {
	a.initRun()
	ret, err := a.orm.Update().SetUpdateFields([]string{"title", "content"}).
		SetUpdateWhereFields("and", []string{"id"}).Where(article).Commit(article)
	if err != nil {
		return article, err
	}
	return ret.(Article), nil
}

func (a *ArticleModel) List(start int64, offset int64) (ar []Article, err error) {
	a.initRun()
	ret, err := a.orm.Select().SetSelectFields([]string{"id", "title", "created_at", "updated_at"}).Limit(start, offset).Where(Article{}).Commit()
	if err != nil {
		return nil, err
	}
	for i, _ := range ret {
		ar = append(ar, ret[i].(Article))
	}

	return ar, nil
}
