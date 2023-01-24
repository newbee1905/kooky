package brave

import (
	"github.com/newbee1905/kooky"
	"github.com/newbee1905/kooky/internal/chrome"
	"github.com/newbee1905/kooky/internal/chrome/find"
	"github.com/newbee1905/kooky/internal/cookies"
)

type BraveFinder struct{}

var _ kooky.CookieStoreFinder = (*BraveFinder)(nil)

func init() {
	kooky.RegisterFinder(`Brave`, &BraveFinder{})
}

func (f *BraveFinder) FindCookieStores() ([]kooky.CookieStore, error) {
	files, err := find.FindBraveCookieStoreFiles()
	if err != nil {
		return nil, err
	}

	var ret []kooky.CookieStore
	for _, file := range files {
		ret = append(
			ret,
			&cookies.CookieJar{
				CookieStore: &chrome.CookieStore{
					DefaultCookieStore: cookies.DefaultCookieStore{
						BrowserStr:           file.Browser,
						ProfileStr:           file.Profile,
						OSStr:                file.OS,
						IsDefaultProfileBool: file.IsDefaultProfile,
						FileNameStr:          file.Path,
					},
				},
			},
		)
	}

	return ret, nil
}
