package models

type AnalyzeOptions struct {
	CheckSecurity   bool `json:"check_security"`
	CheckStyle      bool `json:"check_style"`
	CheckComplexity bool `json:"check_complexity"`
	CheckMetrics    bool `json:"check_metrics"`
}

type AnalyzeRequest struct {
	Code     string         `json:"code" binding:"required"`
	Language string         `json:"language,omitempty"`
	Options  AnalyzeOptions `json:"options"`
}

type UpdateAnalysisRequest struct {
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}
