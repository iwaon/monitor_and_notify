package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/slack-go/slack"
)

func getPage(link string) (data, content string) {
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	content = string(body)
	// return string(content)

	bodyReader := bytes.NewReader(body)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the content
		title := s.Text()
		if strings.Contains(title, "第４回 中学校説明会") {
			fmt.Printf("%d: '%s'\n", i, title)
			data += title
		}
	})

	return
}

func main() {
	var err error

	data, content := getPage("https://foo.net/usr/bar/event/evtIndex.jsf")

	// err = os.WriteFile("data_scraping.html", []byte(data), 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fileNew := "data_scraping_bar_new.html"
	fileOld := "data_scraping_bar.html"

	fileContentNew := "content_scraping_bar_new.html"
	fileContentOld := "content_scraping_bar.html"

	err = os.WriteFile(fileNew, []byte(data), 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(fileContentNew, []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(fileOld); err != nil {
		c, err := os.Create(fileOld)
		if err != nil {
			log.Fatal(err)
		}
		r, err := os.Open(fileNew)
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(c, r)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("created file:", fileOld)
		postToSlack("foo-alive", fmt.Sprintln("created file:", fileOld, ", time:", time.Now()))
	}
	if _, err := os.Stat(fileContentOld); err != nil {
		c, err := os.Create(fileContentOld)
		if err != nil {
			log.Fatal(err)
		}
		r, err := os.Open(fileContentNew)
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(c, r)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("created file:", fileContentOld)
		postToSlack("foo-alive", fmt.Sprintln("created file:", fileContentOld, ", time:", time.Now()))
	}

	out, err := exec.Command("diff", fileOld, fileNew).Output()
	if err != nil {
		fmt.Print(err.Error()) // when diff exists
		// log.Fatal(err)
	}
	fmt.Println("diff length between old and new", len(out))
	postToSlack("foo-alive", fmt.Sprintln("diff len:", len(out), ", time:", time.Now()))

	if len(out) == 0 {
		return
	}

	fmt.Println("diff", string(out))

	// b, err := json.Marshal("diff has been detected:\n" + string(out))
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("diff", string(b))
	// // postToSlack("foo", "diff has been detected: hoge")
	// postToSlack("foo", string(b))
	postToSlack("foo", "<@U047AHGQJCW> diff has been detected\n"+string(out))

	fmt.Println("backup old file")

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}
	timeString := time.Now().In(jst).Format("20060102_150405")

	dest := timeString + "_bar.html"
	c, err := os.Create(dest)
	if err != nil {
		log.Fatal(err)
	}
	r, err := os.Open(fileOld)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(c, r)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(timeString+"_diff_bar.txt", out, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("move new file to old file")
	if err := os.Rename(fileNew, fileOld); err != nil {
		fmt.Print(err.Error())
	}

	dest = timeString + "_content_bar.html"
	c, err = os.Create(dest)
	if err != nil {
		log.Fatal(err)
	}
	r, err = os.Open(fileContentOld)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(c, r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("move new file to old file")
	if err := os.Rename(fileContentNew, fileContentOld); err != nil {
		fmt.Print(err.Error())
	}
}

func postToSlack(channel string, text string) {
	// generate a client with access token
	tkn := "xoxb-foobar-foobar-foobar"
	c := slack.New(tkn)

	// setting the 2nd argument of MsgOptionText() to true escapses special characters
	_, _, err := c.PostMessage(channel,
		slack.MsgOptionText(text, false),
	)
	if err != nil {
		fmt.Print(err.Error())
	}
}
