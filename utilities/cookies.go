package utilities

import (
	"net/http"
	"time"
)

func SetLangCookies(w http.ResponseWriter, lang string) {
	cookie := &http.Cookie{
		Name:     "Lang",
		Value:    lang,
		Path:     "/",
		Expires:  time.Now().Add(336 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

func GetCookie(r *http.Request) *http.Cookie {
	cookie, err := r.Cookie("Lang")
	if err != nil {
		return nil
	}
	return cookie
}
