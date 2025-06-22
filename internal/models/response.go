package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AnalysisResponse struct {
	Language         string            `json:"language"`
	ComplexityScore  int               `json:"complexity_score"`
	SecurityIssues   []SecurityIssue   `json:"security_issues"`
	StyleSuggestions []StyleSuggestion `json:"style_suggestions"`
	Metrics          CodeMetrics       `json:"metrics"`
	OverallScore     float64           `json:"overall_score"`
}

type SecurityIssue struct {
	Type        string `json:"type"`
	Line        int    `json:"line"`
	Column      int    `json:"column"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Pattern     string `json:"pattern"`
}

type StyleSuggestion struct {
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Message  string `json:"message"`
	Rule     string `json:"rule"`
	Severity string `json:"severity"`
}

type CodeMetrics struct {
	LinesOfCode          int     `json:"lines_of_code"`
	LinesOfComments      int     `json:"lines_of_comments"`
	BlankLines           int     `json:"blank_lines"`
	IdentifierCount      int     `json:"identifiers_count"`
	KeywordCount         int     `json:"keywords_count"`
	OperatorCount        int     `json:"operators_count"`
	Functions            int     `json:"functions"`
	Classes              int     `json:"classes"`
	CommentRatio         float64 `json:"comment_ratio"`
	AverageLineLength    float64 `json:"average_line_length"`
	MaxNestingDepth      int     `json:"max_nesting_depth"`
	CyclomaticComplexity int     `json:"cyclomatic_complexity"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	TimeStamp int64  `json:"time_stamp"`
}

type Analysis struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Code      string             `json:"code" bson:"code"`
	Language  string             `json:"language" bson:"language"`
	Response  AnalysisResponse   `json:"response" bson:"response"`
	Timestamp time.Time          `json:"timestamp" bson:"timestamp"`
	Title     string             `json:"title,omitempty" bson:"title,omitempty"`
	Tags      []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
