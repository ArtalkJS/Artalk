# Voting Feature

Artalk supports voting on comments and pages, allowing users to click the "Up" or "Down" buttons to cast their votes. The comment list can be sorted based on the number of votes, which helps users better assess the quality of the comments.

## Comment Voting

The comment voting feature is enabled by default, allowing users to vote on comments.

You can find the "UI Settings" in the Dashboard to modify the "Vote Button" option to enable or disable comment voting.

Environment variable for the vote button:

```
ATK_FRONTEND_VOTE=1
```

Configuration file for the vote button:

```yaml
frontend:
  vote: true
```

### Down Button

By default, Artalk does not display the Down button. You can find the "UI Settings" in the Dashboard to modify the "Vote Down Button" option to enable it.

Environment variable for the Down button:

```
ATK_FRONTEND_VOTE_DOWN=1
```

Configuration file for the Down button:

```yaml
frontend:
  voteDown: true
```

## Page Voting

Artalk supports voting on pages. To enable page voting, you need to add elements in the page to display the voting buttons, which Artalk will automatically initialize on load:

```html
<div>
  <span class="artalk-page-vote-up"></span>
  <span class="artalk-page-vote-down"></span>
</div>
```

### Voted Status Styling

When users click the page voting button, the element will be given an `active` class to indicate that the user has voted. For example:

```html
<span class="artalk-page-vote-up active"></span>
```

You can customize the voted button style using CSS:

```css
.artalk-page-vote-up.active {
  color: #0083ff;
}
```

The default added class name is `active`, but it can be changed using `pageVote.activeClass` in the client configuration:

```js
Artalk.init({
  pageVote: {
    activeClass: 'active',
  },
})
```

### Custom Element Selectors

By default, Artalk searches for `.artalk-page-vote-up` and `.artalk-page-vote-down` as the voting button elements.

You can customize the voting button selectors by modifying the `pageVote.upBtnEl` and `pageVote.downBtnEl` configuration in the client:

```js
Artalk.init({
  pageVote: {
    upBtnEl: '.artalk-page-vote-up',
    downBtnEl: '.artalk-page-vote-down',
  },
})
```

### Further Customizing Page Voting Buttons

If the voting buttons do not contain any child elements, Artalk will output the text "Up (n)" into the element.

If you want to output the vote count into a separate element, you can add a tag inside the button, for example:

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

To further customize, you can replace the text with icons or add other styles.

The default selectors for vote counts are `.artalk-page-vote-up-count` and `.artalk-page-vote-down-count`.

You can modify `pageVote.upCountEl` and `pageVote.downCountEl` to customize the vote count output elements:

```js
Artalk.init({
  pageVote: {
    upCountEl: '.artalk-page-vote-up-count',
    downCountEl: '.artalk-page-vote-down-count',
  },
})
```
