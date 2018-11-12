// novelSpider project main.go
//用于爬取http://www.yunlaige.com的小说

package main

import (
	"fmt"
	//	"io/ioutil"
	//	"net/http"
	//	"bytes"
	//	"html"
	"io"
	"os"
	"regexp"
	"strings"
	//	"golang.org/x/net/html"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	//	"golang.org/x/text/encoding/simplifiedchinese"
	//	"golang.org/x/text/transform"
)

func main() {
	fmt.Println("天诛八尺,还我公图")
	if len(os.Args) < 2 {
		fmt.Println("请输入需要爬取的小说的首页URL,只能爬取yunlaige.com,例如" + os.Args[0] + " http://www.yunlaige.com/html/18/18327/8493460.html")
		return
	}
	inputUrl := os.Args[1]
	pat := `\d+\.html`
	re := regexp.MustCompile(pat)
	firstPageUrl := re.FindString(inputUrl)
	//	fmt.Printf("%+v\n", firstPageUrl)
	BaseUrl := re.ReplaceAllString(inputUrl, "")
	//	fmt.Printf("%+v\n", BaseUrl)
	if firstPageUrl == "" || BaseUrl == "" {
		fmt.Println("输入的小说的首页URL错误,正确示例" + os.Args[0] + " http://www.yunlaige.com/html/18/18327/8493460.html")
		return
	}
	//小说保存的文件名称
	doc, err := goquery.NewDocument(BaseUrl + firstPageUrl)
	if err != nil {
		fmt.Println("爬取首页失败,请确认URL或者网络是否正常")
		return
	}
	bookName := doc.Find("div.bookname h1")
	fmt.Println(ConvertGB2312ToUtf8(bookName.Text(), "gbk", "utf-8"))
	var novelTxtName = "./" + strings.TrimSpace(ConvertGB2312ToUtf8(bookName.Text(), "gbk", "utf-8")) + ".txt"
	if err := httpGet(BaseUrl, firstPageUrl, novelTxtName); err == nil {
		fmt.Println("爬取完成!")
	}
}

func httpGet(base, url, filename string) error {
	//	resp, err := http.Get(url)
	fmt.Println("爬取" + base + url)
	doc, err := goquery.NewDocument(base + url)
	if err != nil {
		return err
	}
	title := doc.Find("p.ctitle")
	fmt.Println(ConvertGB2312ToUtf8(title.Text(), "gbk", "utf-8"))
	//	fmt.Println(ConvertGB2312ToUtf8(title.Text(), "gbk", "utf-8"))
	WirteFile(ConvertGB2312ToUtf8(title.Text(), "gbk", "utf-8"), filename)
	WirteFile("\r\n", filename)

	content := doc.Find("#content")
	if htmlStr, err := content.Html(); err == nil {
		words := strings.Split(htmlStr, "<br/><br/>")
		for _, v := range words {
			if strings.LastIndex(v, "div") == -1 {
				WirteFile(strings.Replace(ConvertGB2312ToUtf8(v, "gbk", "utf-8"), "聽", " ", -1), filename)
				WirteFile("\r\n", filename)
				//fmt.Println(strings.Replace(ConvertGB2312ToUtf8(v, "gbk", "utf-8"), "聽", " ", -1))
			}
		}
	} else {
		return err
	}
	//获取下一页的URL
	var nextPageUrl = ""
	doc.Find("div.bottomlink a").Each(func(index int, item *goquery.Selection) {
		linkTag := item
		if strings.LastIndex(ConvertGB2312ToUtf8(linkTag.Text(), "gbk", "utf-8"), "下一页") >= 0 {
			nextPageUrl, _ = linkTag.Attr("href")
		}
		//		linkTag := item
		//		link, _ := linkTag.Attr("href")
		//		linkText := ConvertGB2312ToUtf8(linkTag.Text(), "gbk", "utf-8")
		//		fmt.Printf("Link #%d: '%s' - '%s'\n", index, linkText, link)
	})
	if len(nextPageUrl) > 0 && strings.TrimSpace(nextPageUrl) != "index.html" {
		httpGet(base, nextPageUrl, filename)
	}
	return nil
}

func ConvertGB2312ToUtf8(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

func WirteFile(wireteString, fileName string) {
	var fileHandler *os.File
	if checkFileIsExist(fileName) { //如果文件存在
		if f, err1 := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR, 0666); err1 == nil {
			fileHandler = f
		} else {
			fmt.Printf("打开文件失败:", err1)
		}
	} else {
		if f, err1 := os.Create(fileName); err1 == nil {
			fileHandler = f
		} else {
			fmt.Printf("新建文件失败:", err1)
		}
	}
	//	wireteString := "[ERROR]" + time.Now().Format("2006-01-02 15:04:05") + "\t" + logstring + "\n"
	if _, err1 := io.WriteString(fileHandler, wireteString); err1 != nil {
		//	if _, err1 := io.WriteString(l.fileHandler, wireteString); err1 != nil {
		fmt.Printf("写文件失败:", err1)
	}
	fileHandler.Close()
}
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
