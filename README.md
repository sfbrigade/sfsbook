# stop-violence

Waffle Backlog [![Stories in Ready](https://badge.waffle.io/sfbrigade/stop-violence.svg?label=ready&title=Ready)](http://waffle.io/sfbrigade/stop-violence)

## About the project
We are building an interactive resource guide based on the [PDF referral guide](http://sfwar.org/pdf/RefGuide.pdf) compiled by SFWAR. The guide will allow volunteers to help maintain the guide to make sure that it is a publicly available source of up-to-date information about the resources available to survivors in the SF Bay Area.

This is the respository for the project that came out of the Code Across 2016. Currently we're in the exploratory phase of the project.

## How to contribute
We are currently in the process of defining the scope of the project. Some pieces are ready to be designed and built, other still need to be figured out.

[Check out our Waffle Backlog to see what issues are ready to be tackled, and what issues are still being defined](http://waffle.io/sfbrigade/stop-violence)


### Research and Notes
We're using the wiki to share our research findings, notes on the project, and design artifacts.

If you want to get started on the research initives we're currently working on, first check out our [Research Guidelines](https://github.com/sfbrigade/stop-violence/wiki/Research-Guidelines). After reading through the guidelines, if you want to help us conduct research you can find interview and observation guides that will get you out into the field and collecting data in the [Research Tools section](https://github.com/sfbrigade/stop-violence/wiki/Research-Tools)

Want to help us design some research? Reach out to Karl (@karlton) on the Brigade Slack.

## Development
The sfsbook application backend is written in [Go](http://golang.org). So install Go (use version 1.7) and then `git clone` this repository. Build the server like so:

```
go build
./sfsbook -init_passwords
```

and you should end up with a server called `sfsbook` ready to run at the top level. The server will manufacture some getting started state in `./state`, dump a varying amount of logging and be ready to browse to at [https://localhost:10443/index.html](https://localhost:10443/index.html).

If you're interested in helping with development, reach out ot Rob (@rjkroege) on the Brigade Slack.
