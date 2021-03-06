package main

import "fmt"
import (
	md52 "crypto/md5"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	m "math/rand"
	"strconv"
	"sync"
	"time"
	"net/http"
	"io/ioutil"
)

const ImageSaveDir = "./images/" // 图片保存目录
var db *sql.DB
var wg sync.WaitGroup

func init() {
	m.Seed(time.Now().Unix()) // 初始化随机种子数
	// runtime.GOMAXPROCS(runtime.NumCPU())
}

// 结构化goods
type Goods struct {
	GoodsDetail struct {
		Alias               string   `json:"alias"`
		Title               string   `json:"title"`
		Price               string   `json:"price"`
		Content             string   `json:"content"`
		StockNum            int      `json:"stock_num"`
		SoldNum             int      `json:"sold_num"`
		AttachmentUrls      []string `json:"attachment_urls"`
		AttachmentThumbUrls []string `json:"thumb_urls"`
	} `json:"goods_detail"`
	ShopDetail struct {
		ShopName    string `json:"team_name"`
		ShopLogo    string `json:"logo"`
		ShopAlias   string `json:"alias"`
		FollowerNum string `json:"follower_num"`
		SellerNum   string `json:"seller_num"`
		GoodsNum    string `json:"goods_num"`
		Contact     struct {
			Mobile string `json:"mobile"`
			Qq     string `json:"qq"`
			Weixin string `json:"weixin"`
		} `json:"contact"`
		Description string `json:"description"`
	} `json:"team_info"`
}

// 保存goods信息
func (g *Goods) save() {
	return
}

// 格式化输出goods信息
func (g *Goods) toString(i int) {
	price, _ := strconv.Atoi(g.GoodsDetail.Price)
	fmt.Printf("[ %d ] [ %s ], [ %s ], [ %d ], [ %s ], [ %s ] \n", i, g.GoodsDetail.Alias, g.GoodsDetail.Title, price/100, g.ShopDetail.ShopName, g.ShopDetail.ShopAlias)
}

func main() {
	// 开启爬虫
	start := time.Now()
	var err error
	// 数据库
	db, err = sql.Open("mysql", "root:123456@tcp(localhost:3306)/spider_youzan?charset=utf8")
	if err != nil {
		log.Printf("%T %+v", err, err)
	}

	db.SetMaxIdleConns(10)                    //连接池中最大空闲连接数
	db.SetMaxOpenConns(150)                   //打开的最大连接数
	db.SetConnMaxLifetime(7200 * time.Second) //连接的最大空闲时间(可选)
	err = db.Ping()
	if err != nil {
		log.Printf("%T %+v", err, err)
	}
	defer db.Close()
	// 运行
	run()

	// 关闭
	end := time.Now()
	fmt.Println("time:", end.Sub(start).Seconds())
}

// 运行
func run() {
	rows, err := db.Query("SELECT `id` as `recordId`, `category`, `goodsAlias`, `json`, `isChanged`, `shopChanged`, `goodsChanged`, `parseTime` FROM `tbl_spider_record`")
	if err != nil {
		log.Printf("%T %+v", err, err)
		return
	}
	defer rows.Close()

	var i int
	for rows.Next() {
		var recordId int      // 记录id
		var category string   // 类目
		var goodsAlias string // 商品别名
		var jsonString string // json
		var isChanged int     // 是否变动
		var shopChanged int   // 店铺是否变动
		var goodsChanged int  // 商品是否变动
		var parseTime int     // 解析时间

		err = rows.Scan(&recordId, &category, &goodsAlias, &jsonString, &isChanged, &shopChanged, &goodsChanged, &parseTime)
		if err != nil {
			log.Printf("%T %+v", err, err)
		}
		// 解析json
		var g Goods
		err = json.Unmarshal([]byte(jsonString), &g)
		if err != nil {
			log.Printf("%T %+v", err, err)
			return
		}

		g.toString(i)
		max_concurrent_count := 10
		tasks := make(chan string, max_concurrent_count)
		wg.Add(max_concurrent_count)
		for gr := 1; gr <= max_concurrent_count; gr++ {
			go fetchImage(tasks, gr)
		}
		// time.Sleep(10 * time.Second)
		for _, imageUrl := range g.GoodsDetail.AttachmentUrls {
			tasks <- imageUrl
		}
		//for range g.GoodsDetail.AttachmentUrls {
		//	fmt.Println(<-ch) // receive from channel ch
		//}
		close(tasks)
		wg.Wait()
		i++
	}
	if err = rows.Close(); err != nil {
		// but what should we do if there's an error?
		log.Printf("%T %+v", err, err)
		return
	}
	fmt.Println("total:", i)
}

// 获取图片
func fetchImage(tasks chan string, gr int) {
	defer wg.Done()

	for {
		task, ok := <-tasks
		if !ok {
			// fmt.Printf("[ %d ] Shutting Down \n", gr)
			return
		}
		// fmt.Printf("[ %s ] [ %d ] Started \n", task, gr)

		// fmt.Printf("[ %s ] [ %s ] [ %d ] Completed \n", time.Now().Format("Jan _2 15:04:05.000000000"), task, gr)
		// ...
		startDownload := time.Now()
		resp, err := http.Get(task)
		if err != nil {
			log.Printf("%T %+v", err, err)
			return
		}
		////body, err := ioutil.ReadAll(resp.Body)
		//fileType := filepath.Ext(task)
		//// 生成md5字符串
		//fileName := GetMd5String(string(m.Int())) + "-" + UniqueId() + fileType
		//dst, err := os.Create(ImageSaveDir + fileName)
		//if err != nil {
		//	log.Printf("%T %+v", err, err)
		//	return
		//}
		////fmt.Println(ImageSaveDir + fileName)
		nbytes,err := io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("%T %+v", err, err)
			return
		}

		//end := time.Since(startDownload).Seconds()
		// fmt.Printf("%.2fs  %7d  %s \n", secs, nbytes, url)
		// ...

		fmt.Printf("[ %s ] [ %s ] [ %7d ] [ %s ] [ %d ] Completed \n", startDownload.Format("15:04:05.000000000"), time.Now().Format("15:04:05.000000000"), nbytes, task, gr)
	}

}

// 生成s的md5
func GetMd5String(s string) string {
	h := md52.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//生成Guid字串
func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}
