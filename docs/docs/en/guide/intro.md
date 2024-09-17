# <span class="wave">👋</span> Hello Friend

Artalk is a simple yet feature-rich commenting system that you can effortlessly deploy and integrate into any blog, website, or web application.

![](https://user-images.githubusercontent.com/22412567/141147152-df30a0ff-bf41-42ee-9958-4722206a7a90.png)

## Key Features

**Lightweight Design**

The front end uses TypeScript (Vanilla JS), lightweight and free of redundant dependencies, only ~40KB (gzipped).

The back end, rewritten in Golang (Artalk v2), is cross-platform, compact, and fully featured, allowing for rapid deployment.

**“Small but Complete”**

Artalk's features include but are not limited to:

<!-- features -->
* [Sidebar](https://artalk.js.org/guide/frontend/sidebar.html): Quick management, intuitive browsing
* [Social Login](https://artalk.js.org/guide/frontend/auth.html): Fast login via social accounts
* [Email Notification](https://artalk.js.org/guide/backend/email.html): Various sending methods, email templates
* [Diverse Push](https://artalk.js.org/guide/backend/admin_notify.html): Multiple push methods, notification templates
* [Site Notification](https://artalk.js.org/guide/frontend/sidebar.html): Red dot marks, mention list
* [Captcha](https://artalk.js.org/guide/backend/captcha.html): Various verification types, frequency limits
* [Comment Moderation](https://artalk.js.org/guide/backend/moderator.html): Content detection, spam interception
* [Image Upload](https://artalk.js.org/guide/backend/img-upload.html): Custom upload, supports image hosting
* [Markdown](https://artalk.js.org/guide/intro.html): Supports Markdown syntax
* [Emoji Pack](https://artalk.js.org/guide/frontend/emoticons.html): Compatible with OwO, quick integration
* [Multi-Site](https://artalk.js.org/guide/backend/multi-site.html): Site isolation, centralized management
* [Admin](https://artalk.js.org/guide/backend/multi-site.html): Password verification, badge identification
* [Page Management](https://artalk.js.org/guide/frontend/sidebar.html): Quick view, one-click title navigation
* [Page View Statistics](https://artalk.js.org/guide/frontend/pv.html): Easily track page views
* [Hierarchical Structure](https://artalk.js.org/guide/frontend/config.html#nestmax): Nested paginated list, infinite scroll
* [Comment Voting](https://artalk.js.org/guide/frontend/config.html#vote): Upvote or downvote comments
* [Comment Sorting](https://artalk.js.org/guide/frontend/config.html#listsort): Various sorting options, freely selectable
* [Comment Search](https://artalk.js.org/guide/frontend/sidebar.html): Quick comment content search
* [Comment Pinning](https://artalk.js.org/guide/frontend/sidebar.html): Pin important messages
* [View Author Only](https://artalk.js.org/guide/frontend/config.html): Show only the author's comments
* [Comment Jump](https://artalk.js.org/guide/intro.html): Quickly jump to quoted comment
* [Auto Save](https://artalk.js.org/guide/frontend/config.html): Content loss prevention
* [IP Region](https://artalk.js.org/guide/frontend/ip-region.html): Display user's IP region
* [Data Migration](https://artalk.js.org/guide/transfer.html): Free migration, quick backup
* [Image Lightbox](https://artalk.js.org/guide/frontend/lightbox.html): Quick integration of image lightbox
* [Image Lazy Load](https://artalk.js.org/guide/frontend/img-lazy-load.html): Lazy load images, optimize experience
* [Latex](https://artalk.js.org/guide/frontend/latex.html): Integrate Latex formula parsing
* [Night Mode](https://artalk.js.org/guide/frontend/config.html#darkmode): Switch to night mode
* [Extension Plugin](https://artalk.js.org/develop/plugin.html): Create more possibilities
* [Multi-Language](https://artalk.js.org/guide/frontend/i18n.html): Switch between multiple languages
* [Command Line](https://artalk.js.org/guide/backend/config.html): Command line operation management
* [API Documentation](https://artalk.js.org/http-api.html): Provides OpenAPI format documentation
* [Program Upgrade](https://artalk.js.org/guide/backend/update.html): Version check, one-click upgrade
<!-- /features -->

We are not exhaustive; more exciting features await your discovery!

**“Unlimited Blade Works”**

Artalk is continually growing; your creativity drives its evolution, and your contributions add value!

Whether it's a front-end project using Vue, React, Svelte, or a blog system like WordPress, Typecho, Hexo, Artalk can be quickly integrated. With everyone's ingenuity, we believe Artalk can adapt to various business scenarios.

## User Experience

We believe elegant design brings excellent user experience, which helps projects go further.

"Ordinary but not mediocre design" - The design tool Figma, favored by professional UI designers, played a significant role in Artalk's redesign. We pre-designed user-friendly, modern interfaces with Figma, then wrote front-end styles to seamlessly blend them into modern websites. This resulted in a simple and fresh interface. Additionally, we designed user identity badges, comment tile/infinite nesting modes, comment pagination, and ensured compatibility with various screen sizes, offering unlimited content in limited space.

"Collapse in an instant" - For unoptimized commenting systems, users might need to repeatedly enter personal information, and their painstakingly typed insights could be lost due to unexpected situations. Knowing that an adult's collapse can happen in an instant, Artalk uses browser caching to auto-fill user information and auto-save comment data, allowing users to express their thoughts with minimal effort.

"Rich site expressions to rekindle comment enthusiasm" - Monotonous emoticons might dampen visitor enthusiasm for commenting. Hence, Artalk comes with a carefully selected humorous emoticon set. Additionally, Artalk supports custom image emoticons.

"Is what you love your life?" - User experience isn't just about visitors; it's also about site admins. Artalk features user-friendly design for site admins, integrating management tools into the [dashboard](./frontend/sidebar.md#控制中心) in the sidebar. Admin users can easily manage multiple sites, with all data exchanged through standardized APIs and processed asynchronously, reducing data processing blockages and service resource usage. For potential spam, Artalk supports automatic filtering, reducing admin workload and keeping the site clean.

We hope Artalk not only fulfills the basic functions of a commenting system but also becomes a bridge for **knowledge sharers and learners to exchange ideas**, helping knowledge creators realize their value.

## Community Philosophy

“**Simplify Complexity, Remarkable Simplicity**”

Artalk aims to achieve **rich** functionality while being as **simple** as possible.

On October 2, 2018, Artalk's [first line of code](https://github.com/ArtalkJS/Artalk/commit/66128e2c8d9a8ac00a8d1498ff0ec035a7727daf) was pushed to GitHub. It wasn't until October 20, 2021, that version v2 was released. Due to a small team and limited developer time, the project's overall progress has been slow. We greatly need the power of the community, whether it's reporting bugs or providing new feature ideas, we eagerly anticipate your contributions.

The Artalk community is an inclusive and open community. We welcome people of all skill levels to help/participate in project development. If you are a beginner, besides actively learning project-related knowledge, you can also try deploying and using Artalk, finding and confirming its shortcomings during use, and then posting relevant discussions in the project's [Issues](https://github.com/ArtalkJS/Artalk/issues) or [Discussions](https://github.com/ArtalkJS/Artalk/discussions) to help developers better locate issues and make optimizations faster. If you are a skilled developer, you can find all the project source code at [@ArtalkJS](https://github.com/ArtalkJS). Combined with this documentation, we believe it shouldn't be difficult to understand. Whether it's optimizing the front-end and back-end structure, implementing new features, or writing community projects, we look forward to Artalk being invigorated with fresh blood.

“More action, less talk” - The Artalk community does not welcome meaningless disputes. We hope community members coexist harmoniously and contribute ideas for community development. Before raising a question, you should read "[How to Ask Questions the Smart Way](https://lug.ustc.edu.cn/wiki/doc/smart-questions/)", as this may determine whether you get a useful answer. Before expressing an opinion, you should maintain basic etiquette, such as keeping a calm attitude, using appropriate language, and avoiding abusive language, sarcasm, or disrespecting others' beliefs and positions.

As advocates and practitioners of open-source spirit, we believe the free software we create should be freely used, studied, modified, and shared. This project's main program is open-sourced under the [MIT](https://github.com/ArtalkJS/Artalk/blob/master/LICENSE) license, and the documentation is under the [CC](https://creativecommons.org/licenses/by-nc-sa/4.0/deed.zh) license.
::: tip Want to Contribute to the Community?

- Browse developer resources ([Developer Documentation](../develop/index.md) / [CONTRIBUTING.md](https://github.com/ArtalkJS/Artalk/blob/master/CONTRIBUTING.md))
- Maintain Artalk backend (code repository [@ArtalkJS/Artalk:/](https://github.com/ArtalkJS/Artalk))
- Maintain Artalk frontend (code repository [@ArtalkJS/Artalk:/ui](https://github.com/ArtalkJS/Artalk/tree/master/ui))
- Improve Artalk documentation (code repository [@ArtalkJS/Artalk:/docs](https://github.com/ArtalkJS/Artalk/tree/master/docs))
- Translate multi-language i18n (see [Multi-language Guide](./frontend/i18n.html))
- Enhance data migration tools (code repository [@ArtalkJS/Artransfer](https://github.com/ArtalkJS/Artransfer))
- Share your ideas (leave a comment below / [Discussions](https://github.com/ArtalkJS/Artalk/discussions))
- Write related community projects (plugins, deployment tutorials, etc.)

:::

## In Conclusion

By now, you should have a basic understanding of Artalk. Whether you choose to use Artalk or not, we greatly appreciate your attention. If Artalk does not yet meet your needs, we hope you can offer sharp suggestions to help it grow.

Welcome to Artalk,

Take off! 🛫️
