//go:build !windows
// +build !windows

package konqueror

import (
	"os"
	"path/filepath"

	"github.com/zellyn/kooky"
	"github.com/zellyn/kooky/internal"
)

type konquerorFinder struct{}

var _ kooky.CookieStoreFinder = (*konquerorFinder)(nil)

func init() {
	kooky.RegisterFinder(`konqueror`, &konquerorFinder{})
}

func (f *konquerorFinder) FindCookieStores() ([]kooky.CookieStore, error) {
	roots, err := konquerorRoots()
	if err != nil {
		return nil, err
	}

	var ret []kooky.CookieStore

	for _, root := range roots {
		var s konquerorCookieStore
		d := internal.DefaultCookieStore{
			BrowserStr:  `konqueror`,
			FileNameStr: filepath.Join(root, `kcookiejar`, `cookies`),
		}
		internal.SetCookieStore(&d, &s)
		s.DefaultCookieStore = d
		ret = append(ret, &s)
	}

	if len(ret) > 0 {
		if cookieStore, ok := ret[len(ret)-1].(*konquerorCookieStore); ok {
			cookieStore.IsDefaultProfileBool = true
		}
	}

	return ret, nil
}

func konquerorRoots() ([]string, error) {
	var ret []string
	// fallback
	if home, err := os.UserHomeDir(); err == nil {
		ret = append(ret, filepath.Join(home, `.local`, `share`))
	}
	if dataDir, ok := os.LookupEnv(`XDG_DATA_HOME`); ok {
		ret = append(ret, dataDir)
	}

	return ret, nil
}
