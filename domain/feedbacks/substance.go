package feedbacks

import (
	"context"
	"fmt"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/jackc/pgtype"
	"strings"
	"time"
)

type Audience struct {
	Id         int64
	Title      string
	Emails     []string
	Disabled   bool
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DisabledAt *time.Time
}

type InvitationType string

const (
	InvitationMultiLink  InvitationType = "multi-link"
	InvitationSingleLink InvitationType = "single-link"
)

type Group struct {
	Id             int64
	Title          string
	InvitationType InvitationType
	Audiences      []int64
	Disabled       bool
	Published      bool
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
	DisabledAt     *time.Time
	PublishedAt    *time.Time
}

type QuestionType string

const (
	QuestionEssay          QuestionType = "essay"
	QuestionMultipleChoice QuestionType = "choice"
)

type Question struct {
	Id         int64
	Sequence   int
	Question   string
	Expect     QuestionType
	Options    []string
	GroupId    int64
	Disabled   bool
	Mandatory  bool
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DisabledAt *time.Time
}

type substanceManager struct {
	dbm database.IManager
}

type ISubstance interface {
	FindAudiences(ctx context.Context, page, limit int) ([]Audience, error)
	FindAudiencesByIds(ctx context.Context, ids ...int64) ([]Audience, error)
	SaveAudiences(ctx context.Context, audiences ...Audience) error
	RemoveAudiencesByIds(ctx context.Context, ids ...int64) error

	FindGroups(ctx context.Context, page, limit int) ([]Group, error)
	FindGroupsByIds(ctx context.Context, ids ...int64) ([]Group, error)
	SaveGroups(ctx context.Context, groups ...Group) error
	RemoveGroupsByIds(ctx context.Context, ids int64) error

	FindQuestions(ctx context.Context, page, limit int) ([]Question, error)
	FindQuestionsByGroupId(ctx context.Context, id int64) ([]Question, error)
	SaveQuestions(ctx context.Context, questions ...Question) error
	UpdateQuestion(ctx context.Context, question Question) error
	RemoveQuestionsByIds(ctx context.Context, ids ...int64) error
}

const (
	tableAudience  string = "audience"
	tableGroups    string = "groups"
	tableQuestions string = "questions"

	substanceDefaultLimit = 20
)

func NewSubstanceDomain(dbm database.IManager) ISubstance {
	return &substanceManager{dbm: dbm}
}

func (s *substanceManager) FindAudiencesByIds(ctx context.Context, ids ...int64) (audiences []Audience, err error) {
	audiences = make([]Audience, 0)
	q := fmt.Sprintf(`
		SELECT 
			id, title, emails, disabled, created_at, updated_at, disabled_at 
		FROM %s
		WHERE id ANY ($1)`, tableAudience)
	aids := &pgtype.Int8Array{}
	if err = aids.Set(ids); err != nil {
		return
	}
	var rows database.Rows
	rows, err = s.dbm.Query(ctx, q, aids)
	if err != nil {
		return
	}
	for rows.Next() {
		var audience Audience
		err = rows.Scan(
			&audience.Id,
			&audience.Title,
			&audience.Emails,
			&audience.Disabled,
			&audience.CreatedAt,
			&audience.UpdatedAt,
			&audience.DisabledAt,
		)
		if err != nil {
			return
		}
		audiences = append(audiences, audience)
	}
	return
}

func (s *substanceManager) FindGroupsByIds(ctx context.Context, ids ...int64) (groups []Group, err error) {
	groups = make([]Group, 0)
	q := fmt.Sprintf(`
		SELECT 
			id, title, invitation_type, audiences, disabled, published, created_at, updated_at, disabled_at, published_at 
		FROM %s
		WHERE id ANY ($1)`, tableGroups)
	aids := &pgtype.Int8Array{}
	if err = aids.Set(ids); err != nil {
		return
	}
	var rows database.Rows
	rows, err = s.dbm.Query(ctx, q, aids)
	if err != nil {
		return
	}
	for rows.Next() {
		var group Group
		err = rows.Scan(
			&group.Id,
			&group.Title,
			&group.InvitationType,
			&group.Audiences,
			&group.Disabled,
			&group.Published,
			&group.CreatedAt,
			&group.UpdatedAt,
			&group.DisabledAt,
			&group.PublishedAt,
		)
		if err != nil {
			return
		}
		groups = append(groups, group)
	}
	return
}

func (s *substanceManager) FindQuestionsByGroupId(ctx context.Context, id int64) (questions []Question, err error) {
	questions = make([]Question, 0)
	q := fmt.Sprintf(`
		SELECT 
			id, sequence, question, expect, options, group_id, mandatory, disabled, created_at, updated_at, disabled_at 
		FROM %s
		WHERE id = $1`, tableQuestions)
	var rows database.Rows
	rows, err = s.dbm.Query(ctx, q, id)
	if err != nil {
		return
	}
	for rows.Next() {
		var question Question
		err = rows.Scan(
			&question.Id,
			&question.Sequence,
			&question.Question,
			&question.Expect,
			&question.Options,
			&question.GroupId,
			&question.Mandatory,
			&question.Disabled,
			&question.CreatedAt,
			&question.UpdatedAt,
			&question.DisabledAt,
		)
		if err != nil {
			return
		}
		questions = append(questions, question)
	}
	return
}

func (s *substanceManager) FindAudiences(ctx context.Context, page, limit int) (audiences []Audience, err error) {
	limit = s.limit(limit)
	page = s.page(page)
	audiences = make([]Audience, 0)
	q := fmt.Sprintf(`
		SELECT 
			id, title, emails, disabled, created_at, updated_at, disabled_at 
		FROM %s
		LIMIT %d OFFSET %d`, tableAudience, limit, (page-1)*limit)
	var rows database.Rows
	rows, err = s.dbm.Query(ctx, q)
	if err != nil {
		return
	}
	for rows.Next() {
		var audience Audience
		err = rows.Scan(
			&audience.Id,
			&audience.Title,
			&audience.Emails,
			&audience.Disabled,
			&audience.CreatedAt,
			&audience.UpdatedAt,
			&audience.DisabledAt,
		)
		if err != nil {
			return
		}
		audiences = append(audiences, audience)
	}
	return
}

func (s *substanceManager) SaveAudiences(ctx context.Context, audiences ...Audience) error {
	q := fmt.Sprintf(`INSERT INTO %s (title, emails, disabled, created_at, disabled_at) VALUES`, tableAudience)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)
	for _, audience := range audiences {
		placeholders = append(placeholders, "(?,?,?,?,?)")
		var disabledAt interface{} = nil
		if audience.Disabled {
			disabledAt = "NOW()"
		}
		values = append(values, audience.Title, audience.Emails, audience.Disabled, "NOW()", disabledAt)
	}
	q = fmt.Sprintf(`
		%s %s
		ON CONFLICT DO UPDATE SET
			title = EXCLUDED.title,
			emails = EXCLUDED.emails,
			disabled = EXCLUDED.disabled
	`, q, strings.Join(placeholders, ","))
	_, err := s.dbm.Prepare(ctx, "save_audiences", q)
	if err != nil {
		return err
	}
	cmd, err2 := s.dbm.Exec(ctx, "save_audiences", values)
	if err2 != nil {
		return err2
	}
	if cmd.RowsAffected() > 0 {
		return nil
	}
	return fmt.Errorf("no rows created")
}

func (s *substanceManager) RemoveAudiencesByIds(ctx context.Context, ids ...int64) (err error) {
	q := fmt.Sprintf(`UPDATE %s SET disabled = 1, disabled_at = NOW() WHERE id = ANY ($1)`, tableAudience)
	aids := &pgtype.Int8Array{}
	if err = aids.Set(ids); err != nil {
		return
	}
	var cmd database.CommandTag
	cmd, err = s.dbm.Exec(ctx, q, aids)
	if err != nil {
		return
	}
	if cmd.RowsAffected() < 1 {
		return fmt.Errorf(`no rows found to delete`)
	}
	return
}

func (s *substanceManager) FindGroups(ctx context.Context, page, limit int) (groups []Group, err error) {
	limit = s.limit(limit)
	page = s.page(page)
	groups = make([]Group, 0)
	q := fmt.Sprintf(`
		SELECT 
			id, title, invitation_type, audiences, disabled, published, created_at, updated_at, disabled_at, published_at 
		FROM %s
		LIMIT %d OFFSET %d`, tableGroups, limit, (page-1)*limit)
	var rows database.Rows
	rows, err = s.dbm.Query(ctx, q)
	if err != nil {
		return
	}
	for rows.Next() {
		var group Group
		err = rows.Scan(
			&group.Id,
			&group.Title,
			&group.InvitationType,
			&group.Audiences,
			&group.Disabled,
			&group.Published,
			&group.CreatedAt,
			&group.UpdatedAt,
			&group.DisabledAt,
			&group.PublishedAt,
		)
		if err != nil {
			return
		}
		groups = append(groups, group)
	}
	return
}

func (s *substanceManager) SaveGroups(ctx context.Context, groups ...Group) error {
	q := fmt.Sprintf(`
		INSERT INTO %s (title, invitation_type, audiences, disabled, published, created_at, disabled_at, published_at) 
		VALUES`, tableGroups)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)
	for _, group := range groups {
		placeholders = append(placeholders, "(?,?,?,?,?,?,?,?)")
		var (
			disabledAt, publishedAt interface{} = nil, nil
		)
		if group.Disabled {
			disabledAt = "NOW()"
		}
		if group.Published {
			publishedAt = "NOW()"
		}
		values = append(values, group.Title, group.InvitationType, group.Audiences, group.Disabled, group.Published,
			"NOW()", disabledAt, publishedAt)
	}
	q = fmt.Sprintf(`
		%s %s
		ON CONFLICT DO UPDATE SET
			title = EXCLUDED.title,
			invitation_type = EXCLUDED.invitation_type,
			audiences = EXCLUDED.audiences,
			disabled = EXCLUDED.disabled,
			published = EXCLUDED.published
	`, q, strings.Join(placeholders, ","))
	_, err := s.dbm.Prepare(ctx, "save_groups", q)
	if err != nil {
		return err
	}
	cmd, err2 := s.dbm.Exec(ctx, "save_groups", values)
	if err2 != nil {
		return err2
	}
	if cmd.RowsAffected() > 0 {
		return nil
	}
	return fmt.Errorf("no rows created")
}

func (s *substanceManager) RemoveGroupsByIds(ctx context.Context, ids int64) (err error) {
	q := fmt.Sprintf(`UPDATE %s SET disabled = 1, disabled_at = NOW() WHERE id = ANY ($1)`, tableGroups)
	aids := &pgtype.Int8Array{}
	if err = aids.Set(ids); err != nil {
		return
	}
	var cmd database.CommandTag
	cmd, err = s.dbm.Exec(ctx, q, aids)
	if err != nil {
		return
	}
	if cmd.RowsAffected() < 1 {
		return fmt.Errorf(`no rows found to delete`)
	}
	return
}

func (s *substanceManager) FindQuestions(ctx context.Context, page, limit int) (questions []Question, err error) {
	limit = s.limit(limit)
	page = s.page(page)
	questions = make([]Question, 0)
	q := fmt.Sprintf(`
		SELECT 
			id, sequence, question, expect, options, group_id, mandatory, disabled, created_at, updated_at, disabled_at 
		FROM %s
		LIMIT %d OFFSET %d`, tableQuestions, limit, (page-1)*limit)
	var rows database.Rows
	rows, err = s.dbm.Query(ctx, q)
	if err != nil {
		return
	}
	for rows.Next() {
		var question Question
		err = rows.Scan(
			&question.Id,
			&question.Sequence,
			&question.Question,
			&question.Expect,
			&question.Options,
			&question.GroupId,
			&question.Mandatory,
			&question.Disabled,
			&question.CreatedAt,
			&question.UpdatedAt,
			&question.DisabledAt,
		)
		if err != nil {
			return
		}
		questions = append(questions, question)
	}
	return
}

func (s *substanceManager) SaveQuestions(ctx context.Context, questions ...Question) error {
	q := fmt.Sprintf(`
		INSERT INTO %s (sequence, question, expect, options, group_id, mandatory, disabled, created_at, disabledAt) 
		VALUES`, tableQuestions)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)
	for _, question := range questions {
		placeholders = append(placeholders, "(?,?,?,?,?,?,?,?,?)")
		var disabledAt interface{} = nil
		if question.Disabled {
			disabledAt = "NOW()"
		}
		values = append(values, question.Sequence, question.Question, question.Expect, question.Options, question.GroupId,
			question.Mandatory, question.Disabled, "NOW()", disabledAt)
	}
	q = fmt.Sprintf(`%s %s`, q, strings.Join(placeholders, ","))
	_, err := s.dbm.Prepare(ctx, "save_questions", q)
	if err != nil {
		return err
	}
	cmd, err2 := s.dbm.Exec(ctx, "save_questions", values)
	if err2 != nil {
		return err2
	}
	if cmd.RowsAffected() > 0 {
		return nil
	}
	return fmt.Errorf("no rows created")
}

func (s *substanceManager) UpdateQuestion(ctx context.Context, question Question) error {
	q := fmt.Sprintf(`
		UPDATE %s 
		SET 
			sequence = ?,
			question = ?,
			expect = ?,
			options = ?,
			group_id = ?,
			mandatory = ?,
			disabled = ?,
			updated_at = NOW(),
			disabled_at = ?
		WHERE id = ?`, tableQuestions)
	_, err := s.dbm.Prepare(ctx, "update_questions", q)
	if err != nil {
		return err
	}
	var disabledAt interface{} = nil
	if question.Disabled {
		disabledAt = "NOW()"
	}
	cmd, err2 := s.dbm.Exec(
		ctx, "update_questions",
		question.Sequence, question.Question, question.Expect, question.Options, question.GroupId,
		question.Mandatory, question.Disabled, disabledAt)
	if err2 != nil {
		return err2
	}
	if cmd.RowsAffected() > 0 {
		return nil
	}
	return fmt.Errorf("no rows created")
}

func (s *substanceManager) RemoveQuestionsByIds(ctx context.Context, ids ...int64) (err error) {
	q := fmt.Sprintf(`UPDATE %s SET disabled = 1, disabled_at = NOW() WHERE id = ANY ($1)`, tableQuestions)
	aids := &pgtype.Int8Array{}
	if err = aids.Set(ids); err != nil {
		return
	}
	var cmd database.CommandTag
	cmd, err = s.dbm.Exec(ctx, q, aids)
	if err != nil {
		return
	}
	if cmd.RowsAffected() < 1 {
		return fmt.Errorf(`no rows found to delete`)
	}
	return
}

func (s *substanceManager) limit(limit int) int {
	if limit < 1 {
		return substanceDefaultLimit
	}
	return limit
}

func (s *substanceManager) page(page int) int {
	if page < 1 {
		return 1
	}
	return page
}
