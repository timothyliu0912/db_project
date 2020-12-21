package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"html/template"
	"io/ioutil"
	//"io"
	"encoding/hex"
	"net/http"
	//"os"
	"crypto/md5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/olivere/elastic/v7"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
    //"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"

	"strings"
	"log"
)

const (
	name     = "root"
	password = "19980912"
	hostname = "127.0.0.1"
	port     = "3306"
	dbName   = "ympt"
)

type JsonResponse struct {
	Data interface{} `json:"data"`
}
type Video struct {
	Url     string `json:"url" bson:"url"`
	Vid     string `json:"vid" bson:"vid"`
	Publish string `json:"published" bson:"published"`
	Title   string `json:"title" bson:"title"`
	Content string `json:"content" bson:"content"`
	Author  string `"json:"author" bson:"author"`
	Favcnt  int `json:"favoritecount" bson:"favoritecount"`
	Viewcnt int `json:"viewcount" bson:"viewcount"`
	Res     int `json:"res" bson:"res"`
	Dur     int `json:"duration" bson:"duration"`
	Cate    string `json:"category" bson:"category"`
}
type DB_video struct {
	Id    string
	Video Video
}
type Videoslice struct {
	Videos []Video
}
var (
    video   Video
    indexName = "youtube"
    servers   = []string{"http://localhost:9200/"}
)
var ytmp = make(map[string]*Video)

func db_con() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "19980912"
	dbName := "ympt"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("connnect success")
	return db
}
func es_open() (client *elastic.Client,  ctx context.Context){
	ctx = context.Background()
	client, err := elastic.NewClient(elastic.SetURL(servers...))
    if err != nil {
        // Handle error
        panic(err)
	}
	return client, ctx
}
func mongo_open()(client *mongo.Client, ctx context.Context){

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	return client, ctx
}
func search(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.ParseForm()
	q := r.FormValue("query")
	t := r.FormValue("type")
	fmt.Println(t, p)
	var videos []Video
	cmp := strings.Compare("ElasticSearch",t)
	if (cmp == 0){
		client,ctx:=es_open()
		termQuery := elastic.NewTermQuery("content", q)
		searchResult, err := client.Search().
        Index("youtube").
        Query(termQuery).
        Sort("viewCount", true). // 按id升序排序
        From(0).Size(50). // 拿前10個結果
        Pretty(true).
		Do(ctx) // 執行
		if err != nil {
			panic(err)
		}
		fmt.Println("success")
		total := searchResult.TotalHits()
		fmt.Printf("Found %d subjects\n", total)
		if total > 0 {
			for _, item := range searchResult.Each(reflect.TypeOf(video)) {
				fmt.Println(item)
				if t, ok := item.(Video); ok {
					//fmt.Printf("Found: youtube(id=%d, title=%s)\n", t.Vid, t.Title)
					videos = append(videos,t)
				}
			}
	
		} else {
			fmt.Println("Not found!")
		}
	} else{
		client,ctx:=mongo_open()
		collection := client.Database("youtube").Collection("videos")
		options := options.Find()
		options.SetSort(bson.D{{"viewcount", -1}})
		options.SetLimit(50)

		query := bson.M{
			"$text": bson.M{
				"$search":q,
				},
			}
		cur, err := collection.Find(ctx,query,options)
		if err != nil { 
			log.Fatal(err) 
		}
		for cur.Next(ctx) {
			var result Video
			err := cur.Decode(&result)
			fmt.Println(cur)
   			if err != nil { 
				log.Fatal(err) 
			}
			videos = append(videos, result)
   // do something with result....
		}
		if err := cur.Err(); err != nil {
  			log.Fatal(err)
		}
	}
	response := &JsonResponse{Data: &videos}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	DB := db_con()
	start := time.Now()
	rows, er := DB.Query("SELECT * FROM videos ORDER BY viewcont DESC limit 50;")
	end := time.Now()
    fmt.Println("insert total time:",end.Sub(start).Seconds())
	if er != nil {
		log.Fatalln(er)
	}
	var db_videos []DB_video
	for rows.Next() {
		var db_video DB_video
		er = rows.Scan(&db_video.Id, &db_video.Video.Url, &db_video.Video.Vid, &db_video.Video.Publish, &db_video.Video.Title, &db_video.Video.Content, &db_video.Video.Author, &db_video.Video.Favcnt, &db_video.Video.Viewcnt, &db_video.Video.Res, &db_video.Video.Dur, &db_video.Video.Cate)
		if er != nil {
			log.Fatalln(er)
		}
		db_videos = append(db_videos, DB_video{
			Id:    db_video.Id,
			Video: db_video.Video,
		})
	}
	rows.Close()
	defer DB.Close()

	data := struct {
		DB_videos []DB_video
	}{
		DB_videos: db_videos,
	}
	var tmpl = template.Must(template.ParseFiles("view/layout.html", "view/index.html", "view/head.html", "view/nav.html"))
	err := tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
func check() {
	DB := db_con()
	rows, er := DB.Query("SELECT count(url) FROM videos")
	if er != nil {
		log.Fatalln(er)
	}
	var cnt int
	for rows.Next() {
		er = rows.Scan(&cnt)
	}
	if cnt == 0 {
		fmt.Println("empty")
		start1 := time.Now()
		byteValue, err :=ioutil.ReadFile("b.json")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Successfully Opened a.json")
		fmt.Println("readAll spend : ", time.Now().Sub(start1))

		var v Videoslice
		json.Unmarshal(byteValue, &v)

		start := time.Now()
		tx,_ := DB.Begin()
		execstring := "insert ignore into videos(id,url,vid,published,title, content, author, favcont, viewcont, res, duration, cate) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		//_, err := DB.Prepare("insert ignore into videos(id,url,published,title, content, author, favcont, viewcont, res, duration, cate) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		for k, i := range v.Videos{
			fmt.Println(i.Vid)
			d := []byte(i.Url)
			b := md5.Sum(d)
			str := hex.EncodeToString(b[:])
			_, err := tx.Exec(execstring,str,i.Url,i.Vid,i.Publish,i.Title,i.Content,i.Author,i.Favcnt,i.Viewcnt,i.Res,i.Dur,i.Cate)
			if err != nil {
				fmt.Println(err)
				return
			}
				fmt.Println("success 10 thousand rows",k)
			// var onedata struct{Id string;Video Video}
			// onedata.Id = str
			// onedata.Video = i
			// in_db = onedata
			//stmt.Exec(in_db.Id,in_db.Video.Url,in_db.Video.Publish,in_db.Video.Title,in_db.Video.Content,in_db.Video.Author,in_db.Video.Favcnt,in_db.Video.Viewcnt,in_db.Video.Res,in_db.Video.Dur,in_db.Video.Cate)
		}
		end := time.Now()
    	fmt.Println("insert total time:",end.Sub(start).Seconds())
		tx.Commit()
	} else {
		fmt.Println("not empty")
	}
	defer DB.Close()
}
func update(w http.ResponseWriter, r *http.Request, p httprouter.Params){
	DB := db_con()
	start := time.Now()
	rows, er := DB.Query("SELECT * FROM videos ORDER BY viewcont DESC limit 50;")
	end := time.Now()
    fmt.Println("insert total time:",end.Sub(start).Seconds())
	if er != nil {
		log.Fatalln(er)
	}
	var db_videos []DB_video
	for rows.Next() {
		var db_video DB_video
		er = rows.Scan(&db_video.Id, &db_video.Video.Url, &db_video.Video.Vid, &db_video.Video.Publish, &db_video.Video.Title, &db_video.Video.Content, &db_video.Video.Author, &db_video.Video.Favcnt, &db_video.Video.Viewcnt, &db_video.Video.Res, &db_video.Video.Dur, &db_video.Video.Cate)
		if er != nil {
			log.Fatalln(er)
		}
		db_videos = append(db_videos, DB_video{
			Id:    db_video.Id,
			Video: db_video.Video,
		})
	}
	rows.Close()
	defer DB.Close()
	data := struct {
		DB_videos []DB_video
	}{
		DB_videos: db_videos,
	}
	var tmpl = template.Must(template.ParseFiles("view/layout.html", "view/up_record.html", "view/head.html", "view/update.html"))
	err := tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
func insert(w http.ResponseWriter, r *http.Request, p httprouter.Params){
	var tmpl = template.Must(template.ParseFiles("view/layout.html", "view/index.html", "view/head.html", "view/insert.html"))
	err := tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
func show(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.ParseForm()
	id := r.FormValue("query")
	fmt.Println(id)
	DB := db_con()
	var db_video DB_video
	row := DB.QueryRow("SELECT * FROM videos WHERE id = ?",id)
	err := row.Scan(&db_video.Id, &db_video.Video.Url, &db_video.Video.Vid, &db_video.Video.Publish, &db_video.Video.Title, &db_video.Video.Content, &db_video.Video.Author, &db_video.Video.Favcnt, &db_video.Video.Viewcnt, &db_video.Video.Res, &db_video.Video.Dur, &db_video.Video.Cate)
	if err != nil {
		log.Fatal(err)
	}

	//response := "here"
	response := &JsonResponse{Data: &db_video}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}
func main() {
	check()
	mux := httprouter.New()
	mux.GET("/", index)
	mux.GET("/search", search)
	mux.GET("/update", update)
	mux.GET("/insert", insert)
	mux.GET("/show_all", show)
	mux.ServeFiles("/js/*filepath", http.Dir("public/js"))
	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
