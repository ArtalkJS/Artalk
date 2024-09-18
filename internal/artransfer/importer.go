package artransfer

import (
	"cmp"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/artalkjs/artalk/v2/internal/utils"
	"github.com/cheggaaa/pb/v3"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type ImportParams struct {
	TargetSiteName string `json:"target_site_name" form:"target_site_name" validate:"optional"` // The target site name
	TargetSiteURL  string `json:"target_site_url" form:"target_site_url" validate:"optional"`   // The target site url
	URLResolver    bool   `json:"url_resolver" form:"url_resolver" validate:"optional"`         // Enable URL resolver
	URLKeepDomain  bool   `json:"url_keep_domain" form:"url_keep_domain" validate:"optional"`   // Keep domain
	JsonFile       string `json:"json_file,omitempty" form:"json_file" validate:"optional"`     // The JSON file path
	JsonData       string `json:"json_data,omitempty" form:"json_data" validate:"optional"`     // The JSON data
	Assumeyes      bool   `json:"assumeyes" form:"assumeyes" validate:"optional"`               // Automatically answer yes for all questions

	console *Console `json:"-"`
}

func (p *ImportParams) SetConsole(c *Console) {
	p.console = c
}

func (p *ImportParams) Console() *Console {
	if p.console == nil {
		p.console = NewConsole()
	}
	return p.console
}

func importArtrans(tx *gorm.DB, params *ImportParams, comments []*entity.Artran) error {
	console := params.Console()
	isTrue := func(val string) bool { return val == "true" || val == "1" }

	if len(comments) == 0 {
		return fmt.Errorf(i18n.T("No comment"))
	}

	if params.TargetSiteURL != "" && !utils.ValidateURL(params.TargetSiteURL) {
		return fmt.Errorf(i18n.T("Invalid {{name}}", map[string]interface{}{"name": i18n.T("Target Site") + " " + "URL"}))
	}

	console.Println()
	console.Print("# " + i18n.T("Please review") + ":\n\n")

	// Print the first comment
	console.PrintEncodeData(i18n.T("First comment"), comments[0])

	// Print parameters
	console.PrintTable([][]interface{}{
		{i18n.T("Target Site") + " " + i18n.T("Name"), cmp.Or(params.TargetSiteName, i18n.T("Unspecified"))},
		{i18n.T("Target Site") + " URL", cmp.Or(params.TargetSiteURL, i18n.T("Unspecified"))},
		{i18n.T("Comment count"), fmt.Sprintf("%d", len(comments))},
		{i18n.T("URL Resolver"), lo.If(params.URLResolver, "on").Else("off")},
	})

	console.Println()

	// Confirm to continue
	if !params.Assumeyes && !console.Confirm(i18n.T("Confirm to continue?")) {
		os.Exit(0)
	}

	console.Println()

	// ---------------------
	//  Start importing
	// ---------------------
	importComments := []*entity.Comment{}
	rawId2GenId := buildGenIdMap(comments)                           // Original ID => GenId (GenId is comment index +1)
	rawRid2RootGenId := buildRid2RootGenIdMap(comments, rawId2GenId) // Original Rid => RootGenId
	createdDates := map[int]time.Time{}
	updatedDates := map[int]time.Time{}

	for i, c := range comments {
		// ---------------------
		//  Prepare site
		// ---------------------
		siteName := strings.TrimSpace(cmp.Or(params.TargetSiteName, c.SiteName))
		siteURLs := strings.TrimSpace(cmp.Or(params.TargetSiteURL, c.SiteURLs))
		if siteName == "" {
			console.Warn(fmt.Sprintf("skip comment id %s since `importParams.target_site_name` and `comment.site_name` are both empty", c.ID))
			continue
		}

		site, sErr := prepareSite(tx, siteName, siteURLs)
		if sErr != nil {
			return fmt.Errorf("failed to prepare site, %w", sErr)
		}

		// ---------------------
		//  Prepare page
		// ---------------------
		pageKey := strings.TrimSpace(c.PageKey)
		if pageKey == "" {
			console.Warn(fmt.Sprintf("skip comment id %s since `comment.page_key` is empty", c.ID))
			continue
		}

		if params.URLResolver { // Enable URL resolver
			splitURLs := utils.SplitAndTrimSpace(params.TargetSiteURL, ",")
			if len(splitURLs) == 0 {
				return fmt.Errorf("\"target_site_url\" cannot be empty if URL resolver is enabled")
			}
			// Use the first URL (form the TargetSiteUrl of import params) as the PageKey (domain part)
			pageKey = getResolvedPageKey(splitURLs[0], c.PageKey)
		}

		if !params.URLResolver && !params.URLKeepDomain { // Strip domain from PageKey
			pageKey = stripDomainFromURL(pageKey)
		}

		page, err := findCreatePage(tx, pageKey, c.PageTitle, site.Name)
		if err != nil {
			return fmt.Errorf("failed to prepare page, %w", err)
		}

		adminOnlyVal := isTrue(c.PageAdminOnly)
		if page.AdminOnly != adminOnlyVal {
			page.AdminOnly = adminOnlyVal
			if pErr := dbSave(tx, &page); pErr != nil {
				return fmt.Errorf("failed to update page, %w", pErr)
			}
		}

		// ---------------------
		//  Prepare user
		// ---------------------
		correctUserBasicInfo(c, console)
		user, uErr := findCreateUser(tx, c.Nick, c.Email, c.Link)
		if uErr != nil {
			return fmt.Errorf("failed to prepare user, %w", uErr)
		}

		if !user.IsAdmin {
			userModified := false
			if c.BadgeName != "" && c.BadgeName != user.BadgeName {
				user.BadgeName = c.BadgeName
				userModified = true
			}
			if c.BadgeColor != "" && c.BadgeColor != user.BadgeColor {
				user.BadgeColor = c.BadgeColor
				userModified = true
			}
			if userModified {
				if uErr := dbSave(tx, &user); uErr != nil {
					return fmt.Errorf("failed to update user, %w", uErr)
				}
			}
		}

		// ---------------------
		//  Prepare vote
		// ---------------------
		voteUp, _ := strconv.Atoi(c.VoteUp)
		voteDown, _ := strconv.Atoi(c.VoteDown)

		// ---------------------
		//  Create new comment
		// ---------------------
		nComment := entity.Comment{
			Rid:    rawId2GenId[c.Rid], // [M_Step.1] Rid => GenId
			RootID: rawRid2RootGenId[c.Rid],

			Content: c.Content,

			UA: c.UA,
			IP: c.IP,

			IsCollapsed: isTrue(c.IsCollapsed),
			IsPending:   isTrue(c.IsPending),
			IsPinned:    isTrue(c.IsPinned),

			VoteUp:   voteUp,
			VoteDown: voteDown,

			UserID:   user.ID,
			PageKey:  page.Key,
			SiteName: site.Name,
		}

		// Prepare slices for restoring CreatedAt and UpdatedAt
		createdDates[i] = parseDate(c.CreatedAt)
		if c.UpdatedAt != "" {
			updatedDates[i] = parseDate(c.UpdatedAt)
		} else {
			updatedDates[i] = parseDate(c.CreatedAt)
		}

		// Append to importComments for batch insert
		importComments = append(importComments, &nComment)
	}

	console.Println(i18n.T("Importing") + "...")
	console.Println()

	// ---------------------
	//  Batch insert
	// ---------------------
	// @link https://gorm.io/docs/create.html#Batch-Insert
	if err := tx.CreateInBatches(&importComments, 100).Error; err != nil {
		return fmt.Errorf("failed to batch insert comments, %w", err)
	}

	// GenId => DBRealId Mapping
	genId2DBRealIdMap := map[uint]uint{}
	for i, savedComment := range importComments {
		genId2DBRealIdMap[uint(i+1)] = savedComment.ID // [M_Step.2] Create GenId => DBRealId Map
	}

	// Progress bar
	var bar *pb.ProgressBar
	if !console.IsOutputFuncSet() {
		bar = pb.StartNew(len(comments))
	}

	total := len(comments)
	for i, savedComment := range importComments {
		// Restore CreatedAt and UpdatedAt

		// Invalid operation for GORM to save `CreatedAt` or `UpdatedAt` field values by using `Create()` or `Save()`,
		// only `Updates()` as a alternative way to update these fields.
		// @see https://gorm.io/zh_CN/docs/conventions.html#CreatedAt
		// @see https://github.com/go-gorm/gorm/issues/4827#issuecomment-960480148
		// savedComment.CreatedAt = createdDates[i]
		// savedComment.UpdatedAt = updatedDates[i]

		updateData := map[string]interface{}{
			"CreatedAt": createdDates[i],
			"UpdatedAt": updatedDates[i],
		}

		// Rebuild Rid
		if savedComment.Rid != 0 {
			updateData["Rid"] = genId2DBRealIdMap[savedComment.Rid] // [M_Step.3] GenId => DBRealId
			updateData["RootID"] = genId2DBRealIdMap[savedComment.RootID]
		}

		// Perform update
		err := tx.Model(&savedComment).Updates(updateData)
		if err.Error != nil {
			return fmt.Errorf("failed to update comment, %w", err.Error)
		}

		// Rebuild Vote (mock vote)
		createMockVote := func(voteType entity.VoteType, count int) error {
			for i := 0; i < count; i++ {
				if vErr := dbSave(tx, &entity.Vote{
					TargetID: savedComment.ID,
					Type:     voteType,
				}); vErr != nil {
					return fmt.Errorf("failed to create vote, %w", vErr)
				}
			}
			return nil
		}

		if savedComment.VoteUp > 0 {
			if err := createMockVote(entity.VoteTypeCommentUp, savedComment.VoteUp); err != nil {
				return err
			}
		}
		if savedComment.VoteDown > 0 {
			if err := createMockVote(entity.VoteTypeCommentDown, savedComment.VoteDown); err != nil {
				return err
			}
		}

		// Output progress
		if bar != nil {
			bar.Increment()
		}
		if console.IsOutputFuncSet() && i%50 == 0 {
			console.Print(fmt.Sprintf("%.0f%%... ", float64(i)/float64(total)*100))
		}
	}

	// Finish progress
	if bar != nil {
		bar.Finish()
	}
	if console.IsOutputFuncSet() {
		console.Println()
	}

	// Done
	console.Println()
	console.Info(i18n.T("{{count}} items imported", map[string]interface{}{"count": len(comments)}))

	return nil
}

// Prepare site
func prepareSite(tx *gorm.DB, targetSiteName string, targetSiteURLs string) (*entity.Site, error) {
	if targetSiteName == "" {
		return nil, fmt.Errorf("target_site_name is required for prepareSite()`")
	}

	site := findSite(tx, targetSiteName)

	// Create site
	if site.IsEmpty() {
		site = entity.Site{}
		site.Name = targetSiteName
		site.Urls = targetSiteURLs

		if err := dbSave(tx, &site); err != nil {
			return nil, fmt.Errorf("failed to create site")
		}

		return &site, nil
	}

	// Edit existing site
	originalURLs := strings.Split(site.Urls, ",")
	newURLs := []string{}
	{
		targetURLsSpit := utils.SplitAndTrimSpace(targetSiteURLs, ",")
		for _, u := range targetURLsSpit {
			isInCurrentList := func(u string) bool {
				for _, uu := range originalURLs {
					if uu == u {
						return true
					}
				}
				return false
			}

			// if target url item not exist in current site url list, then append to newURLs
			if !isInCurrentList(u) {
				newURLs = append(newURLs, u)
			}
		}
	}

	// Update site urls
	if len(newURLs) > 0 {
		newURLs = append(newURLs, originalURLs...) // Prepend new urls
		site.Urls = strings.Join(newURLs, ",")

		if err := dbSave(tx, &site); err != nil {
			return nil, fmt.Errorf("failed to update site urls")
		}
	}

	return &site, nil
}

func buildGenIdMap(comments []*entity.Artran) map[string]uint {
	genIdMap := map[string]uint{}
	for i, c := range comments {
		genIdMap[c.ID] = uint(i + 1) // [Step.0] GenId is the comment index + 1
	}
	return genIdMap
}

func getCommentByGenId(comments []*entity.Artran, genId uint) *entity.Artran {
	return comments[int(genId)-1]
}

func buildRid2RootGenIdMap(comments []*entity.Artran, genIdMap map[string]uint) map[string]uint {
	getRooID := func(rid string) uint {
		// loop to find the root comment
		visited := map[uint]bool{}
		rootId := genIdMap[rid]
		for rootId != 0 && !visited[rootId] {
			visited[rootId] = true // avoid infinite loop (rid = id)

			comment := getCommentByGenId(comments, rootId)
			if comment == nil || comment.Rid == "0" || comment.Rid == "" {
				return rootId
			}

			rootId = genIdMap[comment.Rid]
		}
		return rootId
	}

	rootIdMap := map[string]uint{}
	for _, c := range comments {
		rootIdMap[c.Rid] = getRooID(c.Rid)
	}

	return rootIdMap
}

func correctUserBasicInfo(c *entity.Artran, console *Console) {
	if strings.TrimSpace(c.Nick) == "" {
		console.Warn("Detected empty user nick, set to \"Anonymous\" for comment ID: " + c.ID)
		c.Nick = "Anonymous"
	}
	if strings.TrimSpace(c.Email) == "" || !utils.ValidateEmail(c.Email) {
		console.Warn("Detected invalid user email " + strconv.Quote(c.Email) + ", set to \"anonymous@example.org\" for comment ID: " + c.ID)
		c.Email = "anonymous@example.org"
	}
	if strings.TrimSpace(c.Link) != "" && !utils.ValidateURL(c.Link) {
		console.Warn("Detected invalid user link " + strconv.Quote(c.Link) + ", append \"https://\" for comment ID: " + c.ID)
		c.Link = "https://" + c.Link
	}
}
