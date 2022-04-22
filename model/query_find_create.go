package model

func FindCreateSite(siteName string) Site {
	site := FindSite(siteName)
	if site.IsEmpty() {
		site = NewSite(siteName, "")
	}
	return site
}

func FindCreatePage(pageKey string, pageTitle string, siteName string) Page {
	page := FindPage(pageKey, siteName)
	if page.IsEmpty() {
		page = NewPage(pageKey, pageTitle, siteName)
	}
	return page
}

func FindCreateUser(name string, email string, link string) User {
	user := FindUser(name, email)
	if user.IsEmpty() {
		user = NewUser(name, email, link) // save a new user
	}
	return user
}

func FindCreateNotify(userID uint, lookCommentID uint) Notify {
	notify := FindNotify(userID, lookCommentID)
	if notify.IsEmpty() {
		notify = NewNotify(userID, lookCommentID)
	}
	return notify
}
