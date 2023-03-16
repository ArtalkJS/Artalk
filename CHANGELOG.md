
## [v2.5.0](https://github.com/ArtalkJS/Artalk/compare/v2.4.4...v2.5.0) (2023-03-10)

### Features

* migrate from `echo` to `go-fiber`
* upgrade go to v1.20.1
* display IP region of comment ([#418](https://github.com/ArtalkJS/Artalk/issues/418)) ([#447](https://github.com/ArtalkJS/Artalk/issues/447))
* docker ci add support for building arm64 wheels
* more functions to handle artalk lifecycle ([#426](https://github.com/ArtalkJS/Artalk/issues/426))
* **ui:** add some static methods
* **captcha:** add support for reCAPTCHA and hCaptcha ([#456](https://github.com/ArtalkJS/Artalk/issues/456))
* **captcha:** support turnstile captcha by cloudflare ([#453](https://github.com/ArtalkJS/Artalk/issues/453))
* **i18n:** add i18n support for backend ([#343](https://github.com/ArtalkJS/Artalk/issues/343))
* **i18n:** translations for backend ([#344](https://github.com/ArtalkJS/Artalk/issues/344))
* **i18n:** add i18n support for sidebar ([#353](https://github.com/ArtalkJS/Artalk/issues/353))
* **i18n:** add zh-TW i18n translation for sidebar and app

### Bug Fixes

* **ui:** hash goto function check condition issue
* **lint:** add tsc check before vite compile ([#440](https://github.com/ArtalkJS/Artalk/issues/440))
* **sidebar:** array type config option cannot be shown
* **email:** duplicate sending with multiple admins using same email addrs ([#375](https://github.com/ArtalkJS/Artalk/issues/375))
* **email:** email queue initialization issue ([#374](https://github.com/ArtalkJS/Artalk/issues/374))
* **email:** failback to `email.mail_tpl` if `admin_notify.email.mail_tpl` is empty
* `timeAgo` function does not display the now
* add `.npmignore` to fix NPM publish inclusion issue
* sidebar navigation sorting ([#361](https://github.com/ArtalkJS/Artalk/issues/361))

### Performance Improvements

* improve some css styles
* add graceful shutdown

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
