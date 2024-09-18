package dao_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/artalkjs/artalk/v2/internal/dao"
	"github.com/artalkjs/artalk/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestFindCreateSite(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	t.Run("Create New Site", func(t *testing.T) {
		siteName := "TestCreateNewSite"
		siteURL := "https://artalk.example.com"

		result := app.Dao().FindCreateSite(siteName, siteURL)
		assert.False(t, result.IsEmpty(), "直接获取创建后的站点数据有问题")
		assert.Equal(t, siteName, result.Name)
		assert.Equal(t, siteURL, result.Urls)

		findSite := app.Dao().FindSite(siteName)
		assert.False(t, findSite.IsEmpty(), "找不到创建后的站点")
		assert.Equal(t, app.Dao().CookSite(&result), app.Dao().CookSite(&findSite), "创建后的站点数据有问题")
	})

	t.Run("Find Existed Site", func(t *testing.T) {
		result := app.Dao().FindCreateSite("Site A", "https://qwqaq.com")
		assert.False(t, result.IsEmpty())
		assert.Equal(t, "http://localhost:8080/,https://qwqaq.com", result.Urls) // not modified
	})
}

func TestFindCreatePage(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	t.Run("Create New Page", func(t *testing.T) {
		var (
			pageKey   = "/NewPage.html"
			pageTitle = "New Page Title"
			siteName  = "Site A"
		)

		result := app.Dao().FindCreatePage(pageKey, pageTitle, siteName)
		assert.False(t, result.IsEmpty())
		assert.Equal(t, pageKey, result.Key)
		assert.Equal(t, pageTitle, result.Title)
		assert.Equal(t, siteName, result.SiteName)

		findPage := app.Dao().FindPage(pageKey, siteName)
		assert.False(t, findPage.IsEmpty(), "找不到创建后的页面")
		assert.Equal(t, app.Dao().CookPage(&result), app.Dao().CookPage(&findPage), "创建后的页面数据有问题")
	})

	t.Run("Find Existed Page", func(t *testing.T) {
		result := app.Dao().FindCreatePage("/test/1000.html", "", "Site A")
		assert.False(t, result.IsEmpty())
		findPage := app.Dao().FindPage("/test/1000.html", "Site A")
		assert.Equal(t, app.Dao().CookPage(&findPage), app.Dao().CookPage(&result))
	})

	t.Run("Concurrent FindCreatePage", func(t *testing.T) {
		var (
			pageKey   = "/TEST_CONCURRENT_PAGE_KEY.html"
			pageTitle = "New Page Title " + time.Now().String()
			siteName  = "Site A"
		)

		// simulate concurrent requests
		var wg sync.WaitGroup

		ready := make(chan struct{})

		var idMap sync.Map
		n := 10000
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-ready // wait for all goroutines to start at the same time
				result := app.Dao().FindCreatePage(pageKey, pageTitle, siteName)
				idMap.Store(result.ID, true)
			}()
		}

		close(ready)

		wg.Wait()

		// count the number of different pages
		count := 0
		idMap.Range(func(k, v interface{}) bool {
			t.Log("Page ID", k)
			count++
			return true
		})

		assert.Equal(t, 1, count, fmt.Sprintf("Concurrent FindCreatePage should return the same page, but got %d different pages", count))
	})
}

func TestFindCreateUser(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	t.Run("Create New User", func(t *testing.T) {
		var (
			userName  = "NewUser"
			userEmail = "NewUser@gmail.com"
			userLink  = "https://qwqaq.com"
		)

		result, err := app.Dao().FindCreateUser(userName, userEmail, userLink)
		assert.NoError(t, err)
		assert.False(t, result.IsEmpty())
		assert.Equal(t, userName, result.Name)
		assert.Equal(t, userEmail, result.Email)
		assert.Equal(t, userLink, result.Link)

		findUser := app.Dao().FindUser(userName, userEmail)
		assert.False(t, findUser.IsEmpty(), "找不到创建后的用户")
		assert.Equal(t, app.Dao().CookUser(&result), app.Dao().CookUser(&findUser), "创建后的用户数据有问题")
	})

	t.Run("Valid User Values", func(t *testing.T) {
		args := []struct {
			name   string
			email  string
			link   string
			result bool
		}{
			{"", "", "", false},
			{"userA", "", "", false},
			{"", "user_a@example.com", "", false},
			{"userB", "user_b", "", false},
			{"userC", "user_c@example.com", "https://xxxx.com", true},
		}
		for _, arg := range args {
			_, err := app.Dao().FindCreateUser(arg.name, arg.email, arg.link)
			assert.Equal(t, arg.result, err == nil, "FindCreateUser(%s, %s, %s) should return %v", arg.name, arg.email, arg.link, arg.result)
		}

		// Invalid user link
		u, err := app.Dao().FindCreateUser("userD", "user_d@example.com", "xxxx.com")
		assert.NoError(t, err)
		assert.Equal(t, "", u.Link, "The user should be create but link is empty because it's invalid")
	})

	t.Run("Find Existed User", func(t *testing.T) {
		result, err := app.Dao().FindCreateUser("userA", "user_a@qwqaq.com", "")
		assert.NoError(t, err)
		assert.False(t, result.IsEmpty())
		assert.Equal(t, app.Dao().FindUser("userA", "user_a@qwqaq.com"), result)
	})
}

type mockEntity struct {
	ID   int
	Name string
}

func (e mockEntity) IsEmpty() bool {
	return e.ID == 0
}

func TestFindCreateAction(t *testing.T) {
	app, _ := test.NewTestApp()
	defer app.Cleanup()

	var calledTimes int32

	ready := make(chan struct{})

	t.Run("Concurrent FindCreateAction", func(t *testing.T) {
		var wg sync.WaitGroup

		var atomicResult atomic.Value // Assume that Find and Create actions is Thread-safe
		atomicResult.Store(mockEntity{})

		findCreateFunc := func() (mockEntity, error) {
			return dao.FindCreateAction("key", func() (mockEntity, error) {
				// findAction
				time.Sleep(10 * time.Millisecond)     // Simulate a slow find action
				r := atomicResult.Load().(mockEntity) // Thread-safe read result
				return r, nil
			}, func() (mockEntity, error) {
				// createAction
				atomic.AddInt32(&calledTimes, 1)  // Record the number of times createAction is called
				time.Sleep(10 * time.Millisecond) // Simulate a slow create action
				r := mockEntity{ID: 1, Name: "create"}
				atomicResult.Store(r) // Thread-safe update result
				return r, nil
			})
		}

		const n = 1000
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-ready // wait for all goroutines to start at the same time
				r, err := findCreateFunc()
				assert.NoError(t, err)
				assert.False(t, r.IsEmpty(), "FindCreateAction should always return a non-empty entity")
			}()
		}

		// start all goroutines at the same time
		close(ready)

		wg.Wait()

		if got := atomic.LoadInt32(&calledTimes); got != 1 {
			t.Fatalf("FindCreateAction should only call createAction once, but got %d times", got)
		}
	})
}
