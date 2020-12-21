package main
import 
(
	"encoding/json"
	"database/sql"
	"fmt"
	"strings"
	"os"
	"io/ioutil"
  
	_ "github.com/go-sql-driver/mysql"
  )
const 
(
	name = "root"
	password="19980912"
	hostname = "127.0.0.1"
	port ="3306"
	dbName   = "ympt"
)

func initDB() {
	path := strings.Join([]string{name, ":", password, "@tcp(", hostname, ":", port, ")/", dbName, "?charset=utf8"}, "")
	fmt.Println(path)
	DB, _ := sql.Open("mysql", path)
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	if err := DB.Ping(); err != nil {
		fmt.Println("opon database fail:", err)
		return
	}
	fmt.Println("connnect success")
}
 
type Video struct {
	URL           string `json:"url"`
	Published     string `json:"published"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Author        string `json:"author"`
	FavoriteCount int `json:"favoriteCount"`
	ViewCount     int `json:"viewCount"`
	Res           int `json:"res"`
	Duration      int `json:"duration"`
	Category      string `json:"category"`
}

func main() {
	initDB()
	jsonFile, err := os.Open("a.json")
	if err != nil {
	  fmt.Println(err)
	}
	fmt.Println("Successfully Opened ytmp0.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	fmt.Println(string(byteValue))
	var video Video
	json.Unmarshal(byteValue, &video)
	fmt.Println(video)
  }