package nyaa

import (
	"errors"
	"fmt"
	gq "github.com/PuerkitoBio/goquery"
	"github.com/RangelReale/filesharetop/lib"
	"log"
	"net/url"
	"strconv"
	//"strings"
	//"time"
)

type NYSort int

const (
	NYSORT_SEEDERS  NYSort = 2
	NYSORT_LEECHERS NYSort = 3
	NYSORT_COMPLETE NYSort = 4
)

type NYSortBy int

const (
	NYSORTBY_ASCENDING  NYSortBy = 2
	NYSORTBY_DESCENDING NYSortBy = 1
)

type NYParser struct {
	List   map[string]*fstoplib.Item
	logger *log.Logger
}

func NewNYParser(l *log.Logger) *NYParser {
	return &NYParser{
		List:   make(map[string]*fstoplib.Item),
		logger: l,
	}
}

func (p *NYParser) Parse(category string, sort NYSort, sortby NYSortBy, pages int) error {

	if pages < 1 {
		return errors.New("Pages must be at least 1")
	}

	posct := int32(0)
	for pg := 1; pg <= pages; pg++ {
		var doc *gq.Document
		var e error

		// parse the page
		if doc, e = gq.NewDocument(fmt.Sprintf("http://www.nyaa.se/?cats=%s&sort=%d&order=%d&offset=%d", category, sort, sortby, pg)); e != nil {
			return e
		}

		/*
			// find the ordered column link using the requested order
			valid := doc.Find(fmt.Sprintf("div.b-content table table.lista tr td > a[href^=\"/index.php?page=torrents&active=0&discount=0&order=%d\"]", sort)).First()
			if valid.Length() == 0 {
				return errors.New("Doc not valid")
			}

			// On the ordered column an up or down arrow is added, check if it is present
			if !strings.ContainsAny(valid.Parent().Text(), "\u2191\u2193") {
				return errors.New("Doc not valid 2")
			}
		*/

		// Iterate on each record
		doc.Find("div.content table.tlist tr.tlistrow").Each(func(i int, s *gq.Selection) {
			var se error

			link := s.Find("td.tlistname a").First()
			if link.Length() == 0 {
				//p.logger.Println("ERROR: Link not found")
				return
			}

			href, hvalid := link.Attr("href")
			if !hvalid || href == "" {
				p.logger.Printf("ERROR: Link not found")
				return
			}

			hu, se := url.Parse(href)
			if se != nil {
				p.logger.Printf("ERROR: %s", se)
				return
			}

			lid := hu.Query().Get("tid")
			if lid == "" {
				p.logger.Printf("ERROR: ID not found")
				return
			}

			category := s.Find("td.tlisticon a").First()
			if category.Length() == 0 {
				p.logger.Printf("ERROR: Category not found")
				return
			}
			cathref, catvalid := category.Attr("href")
			if !catvalid || cathref == "" {
				p.logger.Printf("ERROR: Cat link not found")
				return
			}

			cu, se := url.Parse(cathref)
			if se != nil {
				p.logger.Printf("ERROR: %s", se)
				return
			}
			catid := cu.Query().Get("cats")

			seeder := s.Find("td.tlistsn").First()
			if seeder.Length() == 0 {
				p.logger.Printf("ERROR: Seeder not found")
				return
			}
			leecher := s.Find("td.tlistln").First()
			if leecher.Length() == 0 {
				p.logger.Printf("ERROR: Leecher not found")
				return
			}
			complete := s.Find("td.tlistdn").First()
			if complete.Length() == 0 {
				p.logger.Printf("ERROR: Complete not found")
				return
			}
			comments := s.Find("td.tlistmn").First()
			if comments.Length() == 0 {
				p.logger.Printf("ERROR: Comments not found")
				return
			}

			nseeder, se := strconv.ParseInt(seeder.Text(), 10, 32)
			if se != nil {
				p.logger.Printf("ERROR: %s", se)
				return
			}
			nleecher, se := strconv.ParseInt(leecher.Text(), 10, 32)
			if se != nil {
				p.logger.Printf("ERROR: %s", se)
				return
			}
			ncomplete, se := strconv.ParseInt(complete.Text(), 10, 32)
			if se != nil {
				p.logger.Printf("ERROR: %s", se)
				return
			}
			ncomments, se := strconv.ParseInt(comments.Text(), 10, 32)
			if se != nil {
				p.logger.Printf("ERROR: %s", se)
				return
			}

			//fmt.Printf("%s: %s\n", link.Text(), hu.Query().Get("id"))
			item, ok := p.List[lid]
			if !ok {
				item = fstoplib.NewItem()
				item.Id = lid
				item.Title = link.Text()
				item.Link = hu.String()
				item.Count = 0
				item.Category = catid
				//item.AddDate = nadddate.Format("2006-01-02")
				item.Seeders = int32(nseeder)
				item.Leechers = int32(nleecher)
				item.Complete = int32(ncomplete)
				item.Comments = int32(ncomments)
				p.List[lid] = item
			}
			item.Count++
			posct++
			if sort == NYSORT_SEEDERS {
				item.SeedersPos = posct
			} else if sort == NYSORT_LEECHERS {
				item.LeechersPos = posct
			} else if sort == NYSORT_COMPLETE {
				item.CompletePos = posct
			}
		})
	}

	return nil
}
