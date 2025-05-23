# Welcome to Mr Poll, Open Sourced.

Mr Poll is coded in the [Go Programming Language](https://go.dev), you can download Go by clicking [here](https://go.dev/dl/). \
After downloading Go you will have to install the required packages for this Project:
```shell
go mod tidy
```

## Disclaimer 
This is an unfinished project and is likely not the actual code that the Mr Poll bot is running at this time. The aim of this project is to replace Mr Poll's code with a system that is stable and reliable while having a low memory footprint, hence why Go was chosen.

## Packages
If you're interested in looking through the packages that this project uses:
- [Disgo - Discord API Wrapper & More](https://github.com/disgoorg/disgo)
- [Gorm - Database ORM](https://gorm.io)
- [Env - Environment Variable Manager](https://github.com/gofor-little/env)

## Todo

#### Data:
- [ ] Review Database schemas
- [ ] Use a Database
- [ ] Use Redis to cache data

#### General Functionality:
- [ ] List commands on help command

#### Poll Functionality:
- [ ] Create & send polls with command
  - [ ] Yes Or No
  - [ ] Multiple Choice
  - [ ] Single Choice
- [ ] End polls with command
- [ ] Send poll with internal API
- [ ] End polls with internal API
- [ ] Handle user voting
- [ ] Anonymous polls voters
- [ ] Required poll roles
- [ ] Poll timers
