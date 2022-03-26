//go:build windows

package edge

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/zellyn/kooky"
	"github.com/zellyn/kooky/internal"
	"github.com/zellyn/kooky/internal/chrome"
	"github.com/zellyn/kooky/internal/chrome/find"
)

type edgeFinder struct{}

var _ kooky.CookieStoreFinder = (*edgeFinder)(nil)

func init() {
	kooky.RegisterFinder(`edge`, &edgeFinder{})
}

func (f *edgeFinder) FindCookieStores() ([]kooky.CookieStore, error) {
	locApp := os.Getenv(`LocalAppData`)
	if len(locApp) == 0 {
		return nil, errors.New(`%LocalAppData% is empty`)
	}

	var cookiesFiles []kooky.CookieStore

	// old versions
	var se edgeCookieStore
	d := internal.DefaultCookieStore{
		BrowserStr:  `edge`,
		FileNameStr: filepath.Join(locApp, `Microsoft`, `Windows`, `WebCache`, `WebCacheV01.dat`),
	}
	internal.SetCookieStore(&d, &se)
	se.DefaultCookieStore = d
	cookiesFiles = append(cookiesFiles, &se)

	// Blink based
	newRoot := func() ([]string, error) {
		return []string{filepath.Join(locApp, `Microsoft`, `Edge`, `User Data`)}, nil
	}
	blinkCookiesFiles, err := find.FindCookieStoreFiles(newRoot, `edge`)
	if err != nil {
		return nil, err
	}
	for _, cookiesFile := range blinkCookiesFiles {
		var sc chrome.CookieStore
		d := internal.DefaultCookieStore{
			BrowserStr:           cookiesFile.Browser,
			ProfileStr:           cookiesFile.Profile,
			OSStr:                cookiesFile.OS,
			IsDefaultProfileBool: cookiesFile.IsDefaultProfile,
			FileNameStr:          cookiesFile.Path,
		}
		internal.SetCookieStore(&d, &sc)
		sc.DefaultCookieStore = d
		cookiesFiles = append(cookiesFiles, &sc)
	}

	return cookiesFiles, nil
}

/*
https://www.nirsoft.net/utils/edge_cookies_view.html
starting from Fall Creators Update 1709 of Windows 10, the cookies of Microsoft Edge Web browser are stored in the WebCacheV01.dat database
ESE database at %USERPROFILE%\AppData\Local\Microsoft\Windows\WebCache\WebCacheV01.dat (%LocalAppData%\Microsoft\Windows\WebCache\WebCacheV01.dat)
CookieEntryEx_##

https://www.linkedin.com/pulse/windows-10-microsoft-edge-browser-forensics-brent-muir
https://bsmuir.kinja.com/windows-10-microsoft-edge-browser-forensics-1733533818

locations:
%LocalAppData%\Microsoft\Windows\WebCache\WebCacheV01.dat
%LocalAppData%\Microsoft\Edge\User Data\Default

https://www.foxtonforensics.com/browser-history-examiner/microsoft-edge-history-location
v79+:
Edge Cookies are stored in the 'Cookies' SQLite database, within the 'cookies' table.

up to v44:
Edge Cookies are stored in the 'WebCacheV01.dat' ESE database, within the 'CookieEntryEx' containers.

older:
Older versions of Edge stored cookies as separate text files in locations specified within the ESE database.
*/
