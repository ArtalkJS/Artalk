
## [v2.5.2](https://github.com/ArtalkJS/Artalk/compare/v2.5.1...v2.5.2) (2023-03-19)

### Bug Fixes

* **ui:** destroy instance when calling init func twice

### Performance Improvements

* **notify:** shorten the length of notify key ([#466](https://github.com/ArtalkJS/Artalk/issues/466))
* **style:** add `!important` for white-space css in code block ([#467](https://github.com/ArtalkJS/Artalk/issues/467))


## [v2.5.1](https://github.com/ArtalkJS/Artalk/compare/v2.5.0...v2.5.1) (2023-03-16)

### Features

* **go:** upgrade golang to v1.20.2

### Bug Fixes

* **ui/count-widget:** context api undefined issue ([#464](https://github.com/ArtalkJS/Artalk/issues/464))
* **ui/i18n:** duplicate packaging built-in language in the external script
* **ui/paginator:** showErr func call issue
* **ui/sort-dropdown:** dropdown menu disappears after call reload func ([#461](https://github.com/ArtalkJS/Artalk/issues/461))

### Code Refactoring

* **anti_spam/aliyun:** accessing aliyun green text api without sdk ([#459](https://github.com/ArtalkJS/Artalk/issues/459))
* **ui/user:** modify user to standalone module ([#463](https://github.com/ArtalkJS/Artalk/issues/463))


## [v2.5.0](https://github.com/ArtalkJS/Artalk/compare/v2.4.4...v2.5.0) (2023-03-10)

### Features

* migrate from `echo` to `go-fiber`
* upgrade go to v1.20.1
* display IP region of comment ([#418](https://github.com/ArtalkJS/Artalk/issues/418)) ([#447](https://github.com/ArtalkJS/Artalk/issues/447))
* docker ci add support for building arm64 wheels
* more functions to handle artalk lifecycle ([#426](https://github.com/ArtalkJS/Artalk/issues/426))
* **ui:** add some static methods
* **ui/height_limit:** support scrollable height limit area ([#451](https://github.com/ArtalkJS/Artalk/issues/451))
* **ui/sidebar:** add dark mode support ([#450](https://github.com/ArtalkJS/Artalk/issues/450))
* **captcha:** add support for reCAPTCHA and hCaptcha ([#456](https://github.com/ArtalkJS/Artalk/issues/456))
* **captcha:** support turnstile captcha by cloudflare ([#453](https://github.com/ArtalkJS/Artalk/issues/453))
* **i18n:** add i18n support for backend ([#343](https://github.com/ArtalkJS/Artalk/issues/343))
* **i18n:** translations for backend ([#344](https://github.com/ArtalkJS/Artalk/issues/344))
* **i18n:** add i18n support for sidebar ([#353](https://github.com/ArtalkJS/Artalk/issues/353))
* **i18n:** add zh-TW i18n translation for sidebar and app

### Bug Fixes

* **ui:** hash goto function check condition issue
* **ui/conf:** avoid some conf overrides frontend from the backend ([#449](https://github.com/ArtalkJS/Artalk/issues/449))
* **ui/editor:** disable img upload cannot hide its btn
* **ui/i18n:** subscribe event priority issue
* **ui/sidebar:** array type of preference initial data issue
* **ui/sidebar:** array type config option cannot be shown
* **ui/sidebar:** boolean type setting option save issue ([#431](https://github.com/ArtalkJS/Artalk/issues/431)) ([#444](https://github.com/ArtalkJS/Artalk/issues/444))
* **ui/sidebar:** setting item save follow type of template
* **lint:** add tsc check before vite compile ([#440](https://github.com/ArtalkJS/Artalk/issues/440))
* **email:** duplicate sending with multiple admins using same email addrs ([#375](https://github.com/ArtalkJS/Artalk/issues/375))
* **email:** email queue initialization issue ([#374](https://github.com/ArtalkJS/Artalk/issues/374))
* **email:** failback to `email.mail_tpl` if `admin_notify.email.mail_tpl` is empty
* `timeAgo` function does not display the now
* add `.npmignore` to fix NPM publish inclusion issue
* sidebar navigation sorting ([#361](https://github.com/ArtalkJS/Artalk/issues/361))

### Performance Improvements

* improve some css styles
* add graceful shutdown
* **conf/i18n:** detect and change locale when config file contains chinese
* **ui/list:** remove useless function call
* **ui/sidebar:** improve sidebar i18n

### Code Refactoring

* bump to monorepo
* renamed from artalk-go to artalk
* http origin checker
* abstract email service
* project package structure
* remove version two-way check ([#452](https://github.com/ArtalkJS/Artalk/issues/452))
* build scripts and CI tests
* replace pkger with go:embed
* launch with vscode debugger
* **CI:** one-key site creating with artalk integrated
* **CI:** improve build and release workflows ([#358](https://github.com/ArtalkJS/Artalk/issues/358))
* **captcha:** abstract captcha service ([#455](https://github.com/ArtalkJS/Artalk/issues/455))
* **comment:** separate comment ui renders from single file ([#427](https://github.com/ArtalkJS/Artalk/issues/427))
* **style:** convert to use Sass as a style interpreter ([#439](https://github.com/ArtalkJS/Artalk/issues/439))
* **ui:** automatic dependency injection ([#429](https://github.com/ArtalkJS/Artalk/issues/429))
* **ui/checker:** simplify checker lifecycle function param table ([#428](https://github.com/ArtalkJS/Artalk/issues/428))
* **ui/dark-mode:** separate dark mode logic into its own module ([#430](https://github.com/ArtalkJS/Artalk/issues/430))
* **ui/editor:** modify editor ui to standalone module ([#441](https://github.com/ArtalkJS/Artalk/issues/441))
* **ui/editor:** change functions of editor to standalone modules ([#443](https://github.com/ArtalkJS/Artalk/issues/443))
* **ui/height-limit:** modify height limit function to standalone module ([#435](https://github.com/ArtalkJS/Artalk/issues/435))
* **ui/i18n:** improve i18n function to standalone module ([#434](https://github.com/ArtalkJS/Artalk/issues/434))
* **ui/list:** modify list pagination to standalone module ([#437](https://github.com/ArtalkJS/Artalk/issues/437))
* **page/fetch:** remove goquery dependency when extracting page data ([#442](https://github.com/ArtalkJS/Artalk/issues/442))
* **anti_spam/qcloud-tms:** implement qcloud tms api without sdk ([#438](https://github.com/ArtalkJS/Artalk/issues/438))

### Documentation

* migrate ArtalkJS/Docs to monorepo docs
* add `CODE_OF_CONDUCT.md`
* add simplified README for artalk npm package
* add translation section to `CONTRIBUTING.md`
* fix wrong config value ([#371](https://github.com/ArtalkJS/Artalk/issues/371))
* fix broken links ([#364](https://github.com/ArtalkJS/Artalk/issues/364))
* add open api ([#360](https://github.com/ArtalkJS/Artalk/issues/360))
* add `Project Structure` section to `CONTRIBUTING.md`
* refine and add frontend api docs
* update setup-example-site.sh script usage
* init artalk with new frontend api
* **deploy:** add `restart=always` for docker to auto restart ([#425](https://github.com/ArtalkJS/Artalk/issues/425))
* **extras:** add deploy guide for vuepress ([#436](https://github.com/ArtalkJS/Artalk/issues/436))
