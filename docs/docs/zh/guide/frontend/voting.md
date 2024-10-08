# 投票功能

Artalk 支持对评论和页面进行投票，用户可以点击“赞同”或“反对”按钮进行投票。评论列表支持通过投票数进行排序。评论投票功能可以帮助用户更好地了解评论的质量。

## 评论投票

评论投票功能默认启用，用户可以对评论进行投票。

你可以在控制台设置界面找到「界面配置」，修改「投票按钮」选项来启用或禁用评论投票功能。

投票按钮的环境变量：

```
ATK_FRONTEND_VOTE=1
```

投票按钮的配置文件：

```yaml
frontend:
  vote: true
```

### 反对按钮

默认情况下，Artalk 不会显示反对按钮。你可以在控制台设置界面找到「界面配置」，修改「反对按钮」选项来启用反对按钮。

反对按钮的环境变量：

```
ATK_FRONTEND_VOTE_DOWN=1
```

反对按钮的配置文件：

```yaml
frontend:
  voteDown: true
```

## 页面投票

Artalk 支持对页面进行投票，你需要在页面中添加元素来显示页面投票按钮，Artalk 在加载时会自动初始化页面投票按钮：

```html
<div>
  <span class="artalk-page-vote-up"></span>
  <span class="artalk-page-vote-down"></span>
</div>
```

### 已投票状态样式

当用户点击了页面投票按钮后，元素会被添加 `active` 类名，表示用户已经投票，例如：

```html
<span class="artalk-page-vote-up active"></span>
```

你可以通过 CSS 样式来自定义按钮的已投票样式：

```css
.artalk-page-vote-up.active {
  color: #0083ff;
}
```

默认添加的类名为 `active`，可以在客户端通过 `pageVote.activeClass` 来修改：

```js
Artalk.init({
  pageVote: {
    activeClass: 'active',
  },
})
```

### 自定义元素选择器

Artalk 默认查找 `.artalk-page-vote-up` 和 `.artalk-page-vote-down` 作为投票按钮元素。

修改客户端的 `pageVote.upBtnEl` 和 `pageVote.downBtnEl` 配置可以自定义投票按钮选择器：

```js
Artalk.init({
  pageVote: {
    upBtnEl: '.artalk-page-vote-up',
    downBtnEl: '.artalk-page-vote-down',
  },
})
```

### 进一步自定义页面投票按钮

如果投票按钮内没有子元素，Artalk 会输出文字 "赞同 (n)" 到元素中。

如果你想输出投票数量到单独的元素，可以在按钮中添加一个标签，例如：

```html
<div class="artalk-page-vote">
  <span class="artalk-page-vote-up">
    👍 (<span class="artalk-page-vote-up-count"></span>)
  </span>
  <span class="artalk-page-vote-down">
    👎 (<span class="artalk-page-vote-down-count"></span>)
  </span>
</div>
```

更进一步，你可以将文字修改为图标，或者添加其他样式。

投票数选择器默认为 `.artalk-page-vote-up-count` 和 `.artalk-page-vote-down-count`，

可以修改 `pageVote.upCountEl` 和 `pageVote.downCountEl` 自定义投票数输出元素：

```js
Artalk.init({
  pageVote: {
    upCountEl: '.artalk-page-vote-up-count',
    downCountEl: '.artalk-page-vote-down-count',
  },
})
```
