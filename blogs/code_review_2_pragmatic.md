### Code Review - Part 2
This article is part of a series of blogs related to code review. I keep writing down the source code of my own project or
one from some open source community and start providing the review comments as code level commments. This helps me in avoiding
any first conclusion bias(-es) I might form when I look at the code for the first time. This also helps me to think by putting
on a user's hat and review the code from its readability point of view.

There are so many factors that a developer needs to keep in mind, while implementing a logic. There are chances this 
developer goes too far keeping best practices in mind. At the end of the day, this logic will get `10/10` from various linting
and other tools. This sounds great. However, the thing that I have come to realize rather late, is the ability to write code
that is easy to understand by the reviewers and maintainers is the only stuff that matters.

### Smart code vs. Good code
What is smart code. Smart code will get 10 on 10 from the review hubs. Smart code is perhaps the best way to have the logic
implemented for some piece of requirement. Smart code is definitely the case, when a fellow maintainer or reviewer looks at the
code and does not feel at home. This other fella somehow feels it must be good code, since it looks terse and makes use of some
of the language features that is not so common. Only the one who lives in the language can write such a code. There are no 
linting issues and has very good code coverage. This other chap now approves the code with LGTM.

### Being pragmatic
While above was a discussion on how a code can talk about it's smartness, I shall delve into code that is more practical and
easy to understand. I shall paste some source code as well as my review comments that hinge on pragmatism.
