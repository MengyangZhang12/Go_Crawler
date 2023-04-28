package parser

import (
	"fmt"
	"go_crawler/engine"
	"go_crawler/model"
	"regexp"
	"strings"
)

const (
	cityRe      = `<a href="(http://album\.zhenai\.com/u/[0-9]+)"[^>]*>([^<]+)</a>`
	ageRe       = `<td width="[0-9]+"><span class="grayL">年龄：</span>([^<]+)</td>`
	genderRe    = `<td width="[0-9]+"><span class="grayL">性别：</span>([^<]+)</td>`
	marriageRe  = `<td width="[0-9]+"><span class="grayL">婚况：</span>([^<]+)</td>`
	locationRe  = `<td><span class="grayL">居住地：</span>([^<]+)</td>`
	educationRe = `<td><span class="grayL">学   历：</span>([^<]+)</td>`
	heightRe    = `<td width="[0-9]+"><span class="grayL">身   高：</span>([^<]+)</td>`
	incomeRe    = `<td><span class="grayL">月   薪：</span>([^<]+)</td>`
	introduceRe = `<div class="introduce">([^<]+)</div>`
	idUrlRe     = `.*album\.zhenai\.com/u/([\d]+)`
	imageRe     = `<img src="(.+)?.+" alt=".*">`
	nextPageRe  = `<a href="(http://www.zhenai.com/zhenghun/[a-z]+/[0-6]+)">下一页</a>`
)

func getMatches(reRule string, contents []byte) []string {
	reAge := regexp.MustCompile(reRule)
	matchesAge := reAge.FindAllSubmatch(contents, -1)
	retList := make([]string, len(matchesAge))
	for i, m := range matchesAge {
		retList[i] = string(m[1])
	}
	return retList
}

func extractString(contents string, position int, re *regexp.Regexp) string {
	match := re.FindStringSubmatch(contents)
	if len(match) >= 1 {
		return string(match[position])
	} else {
		return "null"
	}
	return ""
}

func getUserInfo(contents string) model.Profile {
	profile := model.Profile{}
	profile.Name = extractString(contents, 2, regexp.MustCompile(cityRe))
	userUrl := extractString(contents, 1, regexp.MustCompile(cityRe))
	profile.UserUrl = userUrl
	profile.UserId = extractString(userUrl, 1, regexp.MustCompile(idUrlRe))
	profile.Age = extractString(contents, 1, regexp.MustCompile(ageRe))
	profile.Gender = extractString(contents, 1, regexp.MustCompile(genderRe))
	profile.Marriage = extractString(contents, 1, regexp.MustCompile(marriageRe))
	profile.Location = extractString(contents, 1, regexp.MustCompile(locationRe))
	profile.Height = extractString(contents, 1, regexp.MustCompile(heightRe))
	profile.Education = extractString(contents, 1, regexp.MustCompile(educationRe))
	profile.Income = extractString(contents, 1, regexp.MustCompile(incomeRe))
	imageUrl := extractString(contents, 1, regexp.MustCompile(imageRe))
	profile.ImageUrl = imageUrl
	profile.ImageList = strings.Split(imageUrl, "?")
	profile.Introduce = extractString(contents, 1, regexp.MustCompile(introduceRe)) //userInfo = fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s, %s", userId, name, age, gender, marriage, location, height, education, income, userUrl, imageList[0])

	return profile
}
func ParseCity(contents []byte) engine.ParseResult {
	re := regexp.MustCompile(`<div class="(list-item"><div class="photo"><a .* <div class="item-btn">打招呼</div></div>)`)
	matches := re.FindAllSubmatch(contents, -1)
	var listStr string
	for _, m := range matches {
		listStr = string(m[0])
	}
	users := strings.Split(listStr, "list-item")
	result := engine.ParseResult{}
	for i, user := range users {
		if i == 0 {
		} else {
			userInfo := getUserInfo(user)
			userItem := engine.Item{}
			userItem.Id = userInfo.UserId
			userItem.Url = userInfo.UserUrl
			userItem.Payload = userInfo
			fmt.Println("Got User Item" + userInfo.Name)
			result.Items = append(result.Items, userItem)
		}
	}
	nextPageUrl := extractString(string(contents), 1, regexp.MustCompile(nextPageRe))
	if nextPageUrl == "null" {
	} else {
		result.Requests = append(result.Requests, engine.Request{
			Url:    nextPageUrl,
			Parser: engine.NewFuncParser(ParseCity, "ParseCity"),
		})
	}
	return result
}
