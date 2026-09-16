package service

import (
	"regexp"
	"sort"
	"strings"

	"notification/internal/model"
)

// variablePattern 匹配模板内容中的 {{variable}} 占位符。
var variablePattern = regexp.MustCompile(`\{\{\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*\}\}`)

// RenderResult 模板渲染结果。
type RenderResult struct {
	Subject string   `json:"subject"`
	Content string   `json:"content"`
	Missing []string `json:"missing,omitempty"`
}

// ExtractVariables 从文本中提取全部占位变量名（去重）。
func ExtractVariables(text string) []string {
	matches := variablePattern.FindAllStringSubmatch(text, -1)
	seen := make(map[string]bool)
	out := make([]string, 0)
	for _, m := range matches {
		if len(m) < 2 || seen[m[1]] {
			continue
		}
		seen[m[1]] = true
		out = append(out, m[1])
	}
	sort.Strings(out)
	return out
}

// RenderTemplate 用给定变量渲染模板的主题与内容，并返回缺失变量列表。
func (s *Service) RenderTemplate(id string, vars map[string]string) (*RenderResult, error) {
	t, err := s.store.GetTemplate(id)
	if err != nil {
		return nil, err
	}
	subject, missSubject := renderString(t.Subject, vars)
	content, missContent := renderString(t.Content, vars)
	missing := uniqueStrings(append(missSubject, missContent...))
	return &RenderResult{Subject: subject, Content: content, Missing: missing}, nil
}

// renderString 替换文本中的占位变量，返回渲染结果与缺失变量名。
func renderString(input string, vars map[string]string) (string, []string) {
	missing := make([]string, 0)
	out := variablePattern.ReplaceAllStringFunc(input, func(match string) string {
		name := strings.TrimSpace(strings.Trim(match, "{}"))
		if v, ok := vars[name]; ok {
			return v
		}
		missing = append(missing, name)
		return match
	})
	return out, missing
}

// uniqueStrings 去重并排序字符串切片。
func uniqueStrings(in []string) []string {
	seen := make(map[string]bool)
	out := make([]string, 0, len(in))
	for _, v := range in {
		if seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

// ValidateTemplateVariables 校验给定变量是否满足模板要求，返回缺失变量。
func (s *Service) ValidateTemplateVariables(id string, vars map[string]string) ([]string, error) {
	t, err := s.store.GetTemplate(id)
	if err != nil {
		return nil, err
	}
	required := ExtractVariables(t.Content)
	required = append(required, ExtractVariables(t.Subject)...)
	missing := make([]string, 0)
	for _, name := range uniqueStrings(required) {
		if v, ok := vars[name]; !ok || strings.TrimSpace(v) == "" {
			missing = append(missing, name)
		}
	}
	return missing, nil
}

// mergeTemplateVariables 合并模板声明的变量与实际内容提取的变量。
func mergeTemplateVariables(t *model.Template) []string {
	declared := append([]string{}, t.Variables...)
	declared = append(declared, ExtractVariables(t.Subject)...)
	declared = append(declared, ExtractVariables(t.Content)...)
	return uniqueStrings(declared)
}
