
## [v2.5.0](https://github.com/ArtalkJS/Artalk/compare/v2.4.4...v2.5.0) (2023-03-10)

### Bug Fixes

* admin_page_get default Sort by create_at DESC.
* Fix wrong SQL syntax to support mysql database.
* add `X-Forwarded-For` proxy header read
* `timeAgo` appears negative number minutes ago
* `timeAgo` function does not display the now
* github actions image cache
* Typecho Importer some bugs.
* Prevent unsafe svg image uploads
* naming conflict of node circleci
* crash when frontend config is not set
* Truncate Chinese character issue.
* Wrong logic of captcha 'always' config.
* sidebar navigation sorting ([#361](https://github.com/ArtalkJS/Artalk/issues/361))
* docker image build script
* add `.npmignore` to fix NPM publish inclusion issue
* missing trusted_domains field
* config alias and Dockerfile
* frontend repository directory issue
* Add lock in ChildCommentCacheSave
* destroy method remove root element
* Add lock prevent query db repeatedly while cache miss.
* Dockerfile install bash
* Cache add mutex lock
* replace gravatar default mirror
* PageEdit not work when Page is not created.
* add whitelist for origin checker
* Geetest js resource url
* Rename field total_parents to total_roots.
* admin_site_edit URL Validator bug.
* Reduce the use of native raw SQL to process boolean judgments and transfer them to ORM for processing.
* **UA:** Cooperate with client-side to correct Win11 UA.
* **admin_edit:** Remove cache before admin_edit.
* **admin_site_edit:** Split urls and TrimSpace to allow Space around sep char.
* **api:** Modify to only api subpath enable SiteOrigin checker
* **api_main_get:** Msg center disable comment-pin.
* **api_version:** Trim version prefix v char
* **artrans:** Frontend import Payload modify to JSON object.
* **artrans:** Modify Typecho boot param keys.
* **artransfer:** trim json string before checking array type.
* **artransfer:** typecho date recovery.
* **artransfer:** incorrect IsPinned reference to IsPending ([#49](https://github.com/ArtalkJS/Artalk/issues/49))
* **artransfer:** improve logger & twikoo importer bug.
* **artransfer:** CreatedAt field recovery & json unmarshal number to float64 issue.
* **cache:** ChildCommentCacheSave repeatedly cache while exist
* **cache:** Use SingleFlight to avoid Cache breakdown.
* **cache:** Type conversion issue in Captcha part when using redis.
* **cache:** Duplicate display reply comments when pid cache is deleted
* **captcha:** Always-mode fail on the first req & modify font size
* **captcha:** Fix lots of bugs.
* **comment_add:** ip region search in comment add api
* **comment_add:** PendingDefault ignore Admin user
* **comment_edit:** Resend email when comment is_pending status is modified
* **conf_admins:** Preserve admin user which not in config file.
* **config:** replace viper with koanf to make map case sensitivity ([#47](https://github.com/ArtalkJS/Artalk/issues/47))
* **cors:** Import cors domains from db site urls.
* **db:** Time type field use pointer to solve 0000-00-00 issue.
* **db:** Fix query syntax error & importer vote recover.
* **db:** Support postgres database.
* **docker:** Docker-compose file add build option
* **docker:** Entrypoint script gen command
* **docker:** Run artalk-go anywhere with a right config
* **dockerfile:** Expose host for external device
* **email:** email queue initialization issue ([#374](https://github.com/ArtalkJS/Artalk/issues/374))
* **email:** duplicate sending with multiple admins using same email addrs ([#375](https://github.com/ArtalkJS/Artalk/issues/375))
* **email:** failback to `email.mail_tpl` if `admin_notify.email.mail_tpl` is empty
* **frontend-conf:** Use pointer with omitempty to ignore frontend config output
* **geetest:** Static page prevent cache
* **img-upload:** Set img_upload.path defualt value
* **importer:** UrlResolver disabled by default
* **importer:** UrlResolver avoid end of slash being removed
* **importer:** Fix too many sql variables Err while too much items.
* **importer:** json decode by universal type.
* **importer:** UrlResolver disabled by default when a TargetSiteUrl is given
* **limit-middleware:** Precise match & fix disable captcha config.
* **lint:** add tsc check before vite compile ([#440](https://github.com/ArtalkJS/Artalk/issues/440))
* **notify:** Remove atk-emoticon img tags
* **sidebar:** array type config option cannot be shown
* **site_origin:** issue of same-origin request under the proxy ([#454](https://github.com/ArtalkJS/Artalk/issues/454))
* **time:** Embed IANA timezone database in Windows build
* **transfer:** Importer boot param parse.
* **trusted_domains:** Extract from full URL with slash suffix & improve referer interceptor
* **trusted_domains:** Disable CORS restrictions
* **typo:** CaptchaCheck typo
* **ui:** hash goto function check condition issue
* **upgrade:** Upgrade cmd load config
* **utils:** func SplitAndTrimSpace remove blank items by default
* **validator:** ValidateURL modify to IsRequestURL
* **vote:** Vote up and down at the same time.

### Code Refactoring

* renamed from artalk-go to artalk
* http origin checker
* Improve GetLinkToReply logic
* Bind API actions under a struct.
* abstract email service
* project package structure
* remove version two-way check ([#452](https://github.com/ArtalkJS/Artalk/issues/452))
* batch removing artalk `-go` postfix
* Replace global variable lib.DB with injected db.
* build scripts and CI tests
* move from tmp dir to root
* merge ci workflows
* merge backend to monorepo branch
* migrate from echo to go-fiber
* replace pkger with go:embed
* **CI:** improve build and release workflows ([#358](https://github.com/ArtalkJS/Artalk/issues/358))
* **api:** Modify admin_comment_edit user merge logic
* **api:** Remove unnecessary parameter of ParamsDecode method
* **api_get:** Comments get api.
* **captcha:** abstract captcha service ([#455](https://github.com/ArtalkJS/Artalk/issues/455))
* **comment:** seperate comment ui renders from single file ([#427](https://github.com/ArtalkJS/Artalk/issues/427))
* **conf:** Rename notify to admin_notify & admin_notify.email option
* **docker:** Use a folder to place config instead of mounting a single file
* **notify_launcher:** notify_launcher refactor and fix some bugs.
* **scripts:** double quote to prevent globbing and word splitting
* **style:** convert to use Sass as a style interpreter ([#439](https://github.com/ArtalkJS/Artalk/issues/439))
* **ui:** automatic dependency injection ([#429](https://github.com/ArtalkJS/Artalk/issues/429))

### Documentation

* add `CODE_OF_CONDUCT.md`
* Update README.md
* add simplified README for artalk npm package
* add translation section to `CONTRIBUTING.md`
* Update README.md
* fix content and update index
* fix wrong config value ([#371](https://github.com/ArtalkJS/Artalk/issues/371))
* fix broken links ([#364](https://github.com/ArtalkJS/Artalk/issues/364))
* add open api ([#360](https://github.com/ArtalkJS/Artalk/issues/360))
* add `Project Structure` section to `CONTRIBUTING.md`
* Update README.md
* refine and add frontend api docs
* example in zh-CN
* update setup-example-site.sh script usage
* init artalk with new frontend api
* migrate ArtalkJS/Docs to monorepo docs
* Update README.md
* Update README.md
* Update README.md
* Update README.md
* Update README.md
* update CONTRIBUTING.md `make debug-build`
* **deploy:** add `restart=always` for docker to auto restart ([#425](https://github.com/ArtalkJS/Artalk/issues/425))
* **extras:** add deploy guide for vuepress ([#436](https://github.com/ArtalkJS/Artalk/issues/436))

### Features

* Typecho Impoter PageKey builder.
* add graceful shutdown
* password support to be encrypted by md5.
* Global site & origin checker & support cookie
* Update config example file.
* display IP region of comment ([#418](https://github.com/ArtalkJS/Artalk/issues/418)) ([#447](https://github.com/ArtalkJS/Artalk/issues/447))
* upgrade go to v1.20.1
* more functions to handle artalk lifecycle ([#426](https://github.com/ArtalkJS/Artalk/issues/426))
* Add script to build artalk-sidebar
* Use bluemonday sanitizer
* Global trusted_domain config
* launch with vscode debugger
* docker ci add support for building arm64 wheels
* one-key site creating with artalk integrated
* trusted_domains disable referer check
* ApiVersion & FeMinVersion.
* Page PV counter
* Pagination of admin_page_get.
* FeMinVersion update to 2.0.5
* root path redict to sidebar login page
* FindCommentChildren query db while cache miss
* Add goreleaser with golang-cross.
* Remove example config allow_origins field
* Specific comments get flat mode.
* Add robots.txt.
* Typecho Impoter.
* admin_users password support to be encrypted by using bcrypt.
* Page title fetch support meta tag redirect link
* Admin update all title of pages.
* Typecho Impoter add some hint.
* **anti-spam:** Support Aliyun for Spam detection.
* **anti-spam:** Support KeyWords for Spam detection.
* **api:** Add Accessible URL fields
* **api:** Modify some APIs method to post
* **api:** Get login status API
* **api:** New statistic api
* **api:** Site cannot find err mark
* **api:** settings get & save
* **api-import-upload:** Frontend importer upload before.
* **api_add:** Allow patch user-agent parameter
* **api_admin_cache:** Add api handle cache
* **api_admin_page_fetch:** Add GetStatus parameter.
* **api_login_status:** Add is_admin field
* **artrans:** Import from Artalk v1 (PHP).
* **artrans:** Export and Import Artrans type data.
* **artrans:** Run on frontend client.
* **cache:** Almost full cache coverage.
* **cache:** Support redis & memcache
* **cache:** Replace built-in json lib with vmihailenco/msgpack to speed up
* **captcha:** add support for reCAPTCHA and hCaptcha ([#456](https://github.com/ArtalkJS/Artalk/issues/456))
* **captcha:** support turnstile captcha by cloudflare ([#453](https://github.com/ArtalkJS/Artalk/issues/453))
* **cmd:** Add workdir parameter & gen cmd.
* **cmd:** add admin command to create new user
* **cmd:** Add config cmd
* **comment:** Add field content_marked to render the markdown content
* **conf:** Support control Frontend Conf by Backend.
* **conf:** add conf option for specific email tpl for admins
* **conf:** auto-gen config file when initializing the app
* **conf:** Support for admin_notify.email.enabled option
* **conf:** Support for admin_notify.noise_mode option
* **conf:** Config database options not just dsn.
* **db:** Remove formal foreign key constraints
* **db:** Allow set table_prefix
* **email-tpl:** add ip variable to email template
* **geetest:** Provide geetest check API for frontend.
* **geetest:** Access geetest API
* **gen:** Gen cmd support overwrite file
* **go:** upgrade to Go 1.19 and update some deps
* **go:** Update golang to v1.18.1
* **i18n:** add i18n support for backend ([#343](https://github.com/ArtalkJS/Artalk/issues/343))
* **i18n:** translations for backend ([#344](https://github.com/ArtalkJS/Artalk/issues/344))
* **i18n:** add i18n support for sidebar ([#353](https://github.com/ArtalkJS/Artalk/issues/353))
* **i18n:** add zh-TW i18n translation for sidebar and app
* **img-upload:** Add img-upload option to config file.
* **img-upload:** Image Upload Api
* **img-upload:** Upload rate limitation.
* **img-upload:** Upload using upgit
* **img-upload:** Admin always enable img_upload function.
* **importer:** import from Twikoo.
* **importer:** import from valine.
* **importer:** Speed up importer & fix some bugs.
* **model:** Add IsPinned field into commnent model.
* **multi-site:** Isolate admin users among sites
* **notify:** Support WebHook notify
* **notify:** Various ways to send notify message for admin.
* **notify:** Support template & fix some bugs
* **pin:** Update admin_comment_edit api to edit IsPinned
* **pin:** Pin comments keep most top & Child-comments recursive query use Cooked.
* **pv:** Make route rule & create pvs table
* **self-upgrade:** Add cmd Self-upgrade function
* **tencent_tms:** Support tencent tms for Spam detection.
* **ui:** add some static methods
* **upgit:** Config del_local to remove local file after uploaded.
* **upgrade:** Upgrade command allow ignore version check.
* **user:** Admin email sending isolation between sites.
* **version:** Update Frontend min version to 2.2.1

### Performance Improvements

* **conf:** only provide conf for frontend when fetching the first page
* **docker-build:** cache by separating the downloading of deps
* **trusted_domains:** Improve origin checker

