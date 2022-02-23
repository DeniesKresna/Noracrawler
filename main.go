package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Data struct {
	Date     time.Time
	Title    string
	HClub    string
	AClub    string
	Pred1    int
	PredX    int
	Pred2    int
	PredTips string
	Odd1     float64
	OddX     float64
	Odd2     float64
	HGoal    int
	AGoal    int
	TGoal    int
}

func main() {
	var dateHtml string
	fmt.Println("Masukkan tanggal (dd-mm-yyyy) : ")
	fmt.Scanln(&dateHtml)

	re := regexp.MustCompile(`^([0-2][0-9]|(3)[0-1])(\-)(((0)[0-9])|((1)[0-2]))(\-)\d{4}$`)
	if !re.MatchString(dateHtml) {
		fmt.Println("Inputan Salah")
		return
	}
	dateHtmls := strings.Split(dateHtml, "-")
	yearSearch, _ := strconv.Atoi(dateHtmls[2])

	c := colly.NewCollector(
		colly.AllowedDomains("norabet.com/",
		"www.norabet.com/","https://norabet.com/",
		"norabet.com",
		"www.norabet.com","https://norabet.com"),
	)
	fmt.Println("set allowed domain norabet ok")

	datas := []Data{}
	c.OnHTML("table.content_table", func(e *colly.HTMLElement){
		//log.Println("Course found", e.Request.URL)
		e.ForEach("table.content_table > tbody > tr", func(i int, tr *colly.HTMLElement){
			if i < 2 || i > 2 {
				return
			}
			var newData = Data{}
			tr.ForEach("td", func(j int, td *colly.HTMLElement){
				fmt.Println("baris ", j)
				switch td.Index {
					case 0:{
						var textnya = td.ChildText("noscript")
						var texts = strings.Split(textnya, ",")
						var textDate = strings.Split(texts[0],"-")
						var textHour = strings.Split(texts[1],":")
						var dt,_ = strconv.Atoi(textDate[0])
						var mt,_ = strconv.Atoi(textDate[1])
						var hr,_ = strconv.Atoi(textHour[0])
						var mn,_ = strconv.Atoi(textHour[1])
						matchTime := time.Date(yearSearch, time.Month(mt), dt, hr, mn, 0, 0, time.Now().Local().UTC().Location())
						newData.Date = matchTime
					}			
					case 1: {
						var textnya = td.ChildAttr("img", "title")
						newData.Title = textnya
						var clubMatch = td.Text
						var clubs = strings.Split(clubMatch, "-")
						newData.HClub = strings.Trim(clubs[0], " ")
						newData.AClub = strings.Trim(clubs[1], " ")
					}		
					case 6: {
						var textnya = td.Text
						newData.Pred1, _ = strconv.Atoi(strings.TrimRight(textnya, "%"))
					}
					case 7: {
						var textnya = td.Text
						newData.PredX, _ = strconv.Atoi(strings.TrimRight(textnya, "%"))
					}
					case 8: {
						var textnya = td.Text
						newData.Pred2, _ = strconv.Atoi(strings.TrimRight(textnya, "%"))
					}
					case 9: {
						var textnya = td.Text
						newData.PredTips = textnya
					}
					case 12: {
						var textnya = td.Text
						newData.Odd1, _ = strconv.ParseFloat(textnya, 32)
					}
					case 13: {
						var textnya = td.Text
						newData.OddX, _ = strconv.ParseFloat(textnya, 32)
					}
					case 14: {
						var textnya = td.Text
						newData.Odd2, _ = strconv.ParseFloat(textnya, 32)
					}
					case 15: {
						var textnya = td.Text
						var goals = strings.Split(textnya, ":")
						if len(goals) == 2 {
							var hGoal, _ = strconv.Atoi(goals[0])
							var aGoal, _ = strconv.Atoi(goals[1])
							newData.HGoal = hGoal
							newData.AGoal = aGoal
							newData.TGoal = aGoal + hGoal
						}
					}
				}
			})
			datas = append(datas, newData)
		})
	})
	fmt.Println("finish set scrape configuration")

	c.OnHTML("title", func(e *colly.HTMLElement){
		fmt.Println(e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	var StatusCode = 400
	c.OnResponse(func(r *colly.Response){
		StatusCode = r.StatusCode
		fmt.Println(StatusCode)
	})

	fmt.Println("Scraping start to https://www.norabet.com/predictions-" + dateHtml + ".html")
	c.Visit("https://www.norabet.com/predictions-" + dateHtml + ".html")
	fmt.Println(datas)
}