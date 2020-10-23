package db

import "time"

type Group struct {
	Id          int
	Title       string
	Audience    []string
	Disabled    bool
	CreatedDate time.Time
	UpdatedDate time.Time
}

type Question struct {
	Id          int64
	Question    string
	Expect      string
	Options     []map[string]string
	GroupId     int
	Disabled    bool
	CreatedDate time.Time
	UpdatedDate time.Time
}

type Link struct {
	Id              int64
	Hash            string
	PIN             string
	QuestionGroupId int
	Disabled        bool
	CreatedDate     time.Time
	UpdatedDate     time.Time
}

type LinkItem struct {
	Link      Link
	Questions []Question
	Group     Group
}

type QuestionItem struct {
	Question Question
	Group    Group
}

type User struct {
	Id int64
	Username string
	Password string
}