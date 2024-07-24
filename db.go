package main

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type SiteConfig struct {
	Id               int    `json:"id"`
	Domain           string `json:"domain"`
	IndexTitle       string `json:"index_title"`
	IndexKeywords    string `json:"index_keywords"`
	IndexDescription string `json:"index_description"`
	TemplateName     string `json:"template_name"`
	Routes           string `json:"routes"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

type Article struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	TypeId    int    `json:"type_id"`
	TypeName  string `json:"type_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func initDB() error {
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(127.0.1:3306)/site")
	if err != nil {
		return err
	}
	db.SetConnMaxIdleTime(time.Hour)
	db.SetConnMaxLifetime(24 * time.Hour)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(20)
	return nil
}

func loadConfigs() (*sync.Map, error) {
	stmt, err := db.Prepare(`select * from site_config`)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			panic(err)
		}
	}()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	var result sync.Map
	for rows.Next() {
		s := new(SiteConfig)
		err = rows.Scan(
			&s.Id,
			&s.Domain,
			&s.IndexTitle,
			&s.IndexKeywords,
			&s.IndexDescription,
			&s.TemplateName,
			&s.Routes,
			&s.CreatedAt,
			&s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result.Store(s.Domain, s)
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			panic(err)
		}
	}()

	return &result, nil

}

func GetArticleList(typeId int, page int, size int, order string, direction string) ([]*Article, error) {
	offset := (page - 1) * size
	stmt, err := db.Prepare(`select * from article where type_id=? order by ? ? limit ?,?`)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			panic(err)
		}
	}()
	rows, err := stmt.Query(typeId, order, direction, offset, size)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			panic(err)
		}
	}()
	var articles []*Article
	for rows.Next() {
		a := new(Article)
		err = rows.Scan(&a.Id, &a.Title, &a.Summary, &a.Content, &a.Author, &a.TypeId, &a.TypeName, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}
	return articles, nil

}
func GetArticle(id string) (*Article, error) {
	articelId, err := strconv.Atoi(id)
	if err != nil {
		var times int = 1
		for {
			articelId, _ = GetRandomArticleId()
			if articelId > 0 || times > 4 {
				break
			}
			times += 1
		}
	}
	if articelId <= 0 {
		return nil, errors.New("获取随机ID错误")
	}
	return QueryArticle(articelId)

}
func QueryArticle(articleId int) (*Article, error) {
	stmt, err := db.Prepare(`select * from article where id=?`)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			panic(err)
		}
	}()
	row := stmt.QueryRow(articleId)
	if row.Err() != nil {
		return nil, err
	}
	a := new(Article)
	err = row.Scan(&a.Id, &a.Title, &a.Summary, &a.Content, &a.Author, &a.TypeId, &a.TypeName, &a.UpdatedAt, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return a, nil
}
func GetRandomArticleId() (int, error) {
	stmt, err := db.Prepare("select count(id) from article")
	if err != nil {
		return 0, err
	}

	row := stmt.QueryRow()
	if row.Err() != nil {
		return 0, err
	}
	var count int
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}
	err = stmt.Close()
	if err != nil {
		return 0, err
	}
	order := "order by id "
	switch rand.Intn(3) {
	case 0:
		order = "order by id "
	case 1:
		order = "order by title "
	case 2:
		order = "order by created_at "
	}
	offset := rand.Intn(count)
	s := fmt.Sprintf("select id from article %s limit %d,1", order, offset)
	stmt, err = db.Prepare(s)
	if err != nil {
		return 0, err
	}
	row = stmt.QueryRow()
	if row.Err() != nil {
		return 0, err
	}
	var articleId int
	err = row.Scan(&articleId)
	if err != nil {
		return 0, err
	}
	return articleId, nil
}