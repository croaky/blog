# TODO Code Comments

> "It is okay to live with technical debt,
> as long as the amount of debt is known."

- [Practical Go](https://dave.cheney.net/practical-go/presentations/qcon-china.html#_dont_comment_bad_code_rewrite_it):

I once felt `TODO` comments should be linted out of a codebase:
either create a card in the project management system or fix it now.

After some experience on teams who used TODOs liberally,
I've adjusted my loose personal guidelines.

## Working alone

When working alone, the code and TODO comments are the only
project management system. Trello / Notion boards are not needed.

## Working on a team

When working on a team, link to a card in the project management system in the
TODO comment. A teammate asking for it during code review is an inevitability.

Annotate a TODO comment with your username:

```
// TODO(dfc) this is O(N^2), find a faster way to do this.
```

The username is not a commitment by the person to fix the issue
but they may be the best person to ask when it is time to fix it.
