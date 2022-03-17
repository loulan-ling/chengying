package main

import (
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"os"
	"regexp"
	"strconv"
	"time"
)

var (
	domainArray []string
	mailArray   []string
	domain      string
	key         string

)

func main() {
	help()
	flag.Parse()
	flag.VisitAll(func (f *flag.Flag) {
		if f.Value.String()=="" {
			flag.PrintDefaults()
			os.Exit(1)
		}
	})
	getGithubDomain(domain,key)
	fmt.Println("---------------------------------------------------------")
	fmt.Println("子域名地址：")
	for _,v := range domainArray{
		fmt.Println(v)
	}
	fmt.Println("---------------------------------------------------------")
	fmt.Println("邮箱地址：")
	for _,v := range mailArray{
		fmt.Println(v)
	}
}

func help() {
	flag.StringVar(&domain, "u", "", "google.com")
	flag.StringVar(&key, "k", "", "GitHub API Key")
}

func getGithubDomain(domain, key string) {
	fmt.Printf ("%s%s%s \r\n","获取",domain,"在GitHub的数据")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: key},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	opts := &github.SearchOptions{
		TextMatch: true,
	}
	q := domain
	for {
		code, resp, err := client.Search.Code(ctx, q, opts)
		if err != nil {
			times := resp.Response.Header["X-Ratelimit-Reset"][0]
			t := time.Now().Unix()
			i, _ := strconv.Atoi(times)
			timSleep := i - int(t)
			time.Sleep(time.Duration(timSleep+3) * time.Second)
		} else {
			for _, y := range code.CodeResults {
				domainName := regexp.MustCompile(`[a-zA-Z0-9\-\\.]{2,64}` + domain + `[:\d*]{0,5}`)
				if domainName == nil {
					fmt.Println("获取域名信息出现错误")
					return
				}
				domainNames := domainName.FindAllStringSubmatch(y.String(), -1)
				for _, y := range domainNames {
					if y[0] != "www."+domain {
						if domainRemoveRepeat(y[0]) == true {
							domainArray = append(domainArray, y[0])
						}
					}
				}
				mail := regexp.MustCompile(`[a-zA-Z0-9\@\_]{2,64}` + domain)
				if mail == nil {
					fmt.Println("获取邮箱信息出现错误")
					return
				}
				mails := mail.FindAllStringSubmatch(y.String(), -1)
				for _, y := range mails {
					if y[0] != "www."+domain {
						if mailRemoveRepeat(y[0]) == true {
							mailArray = append(mailArray, y[0])
						}
					}
				}

			}
			fmt.Printf ("%s%d%s \n\r","正在获取Github第",resp.NextPage-1,"页数据")
			if resp.NextPage == 0 {
				break
			} else {
				time.Sleep(2 * time.Second)
				opts.Page = resp.NextPage
			}
		}
	}
}

func domainRemoveRepeat(domainInfo string) bool {
	for _, x := range domainArray {
		if x == domainInfo {
			return false
		}
	}
	return true
}

func mailRemoveRepeat(mailInfo string) bool {
	for _, x := range mailArray {
		if x == mailInfo {
			return false
		}
	}
	return true
}
