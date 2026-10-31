package entity

import "github.com/google/uuid"

type IntegrityCheckIssueLevel int

const (
	IntegrityCheckIssueLevelWarning IntegrityCheckIssueLevel = 0
	IntegrityCheckIssueLevelError   IntegrityCheckIssueLevel = 1
)

type MovieIntegrityCheckIssue struct {
	ID         uuid.UUID                `db:"id"`
	MovieID    uuid.UUID                `db:"movie_id"`
	IssueLevel IntegrityCheckIssueLevel `db:"issue_level"`
	Issue      string                   `db:"issue"`
}
