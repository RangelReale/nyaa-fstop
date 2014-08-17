package nyaa

import (
	"github.com/RangelReale/filesharetop/lib"
	"io/ioutil"
	"log"
)

type Fetcher struct {
	logger *log.Logger
}

func NewFetcher() *Fetcher {
	return &Fetcher{
		logger: log.New(ioutil.Discard, "", 0),
	}
}

func (f *Fetcher) ID() string {
	return "NYAA"
}

func (f *Fetcher) SetLogger(l *log.Logger) {
	f.logger = l
}

func (f *Fetcher) Fetch() (map[string]*fstoplib.Item, error) {
	parser := NewNYParser(f.logger)

	cat, err := f.CategoryMap()
	if err != nil {
		return nil, err
	}

	for _, catlist := range *cat {
		for _, c := range catlist {
			// parse 1 pages ordered by seeders
			err := parser.Parse(c, NYSORT_SEEDERS, NYSORTBY_DESCENDING, 1)
			if err != nil {
				return nil, err
			}

			// parse 1 pages ordered by leechers
			err = parser.Parse(c, NYSORT_LEECHERS, NYSORTBY_DESCENDING, 1)
			if err != nil {
				return nil, err
			}
		}
	}

	return parser.List, nil
}

func (f *Fetcher) CategoryMap() (*fstoplib.CategoryMap, error) {
	return &fstoplib.CategoryMap{
		"ANIME":      []string{"1_37"},
		"LIVEACTION": []string{"5_19"},
		"MUSIC":      []string{"3_0"},
		"MANGA":      []string{"2_12"},
	}, nil
}
