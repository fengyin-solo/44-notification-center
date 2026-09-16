package service

import (
	"sort"
	"strings"
	"time"

	"notification/internal/model"
	"notification/pkg/idgen"
)

// CreateTemplate 创建消息模板。
func (s *Service) CreateTemplate(input model.Template) (*model.Template, error) {
	t := input
	t.Variables = normalizeVariables(t.Variables)
	if err := t.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	t.ID = idgen.Hex()
	t.CreatedAt = now
	t.UpdatedAt = now
	if err := s.store.CreateTemplate(&t); err != nil {
		return nil, err
	}
	s.log.Infof("模板创建成功: id=%s name=%s", t.ID, t.Name)
	return &t, nil
}

// GetTemplate 获取模板详情。
func (s *Service) GetTemplate(id string) (*model.Template, error) {
	return s.store.GetTemplate(id)
}

// ListTemplates 分页查询模板。
func (s *Service) ListTemplates(filter model.TemplateFilter, page, size int) ([]*model.Template, int, error) {
	all := s.store.ListTemplates()
	matched := make([]*model.Template, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		if matched[i].Version != matched[j].Version {
			return matched[i].Version > matched[j].Version
		}
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	return paginate(matched, page, size), total, nil
}

// UpdateTemplate 更新模板并递增版本号。
func (s *Service) UpdateTemplate(id string, input model.Template) (*model.Template, error) {
	exist, err := s.store.GetTemplate(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Type = input.Type
	exist.Subject = input.Subject
	exist.Content = input.Content
	exist.Variables = normalizeVariables(input.Variables)
	exist.Version++
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateTemplate(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

// DeleteTemplate 删除模板。
func (s *Service) DeleteTemplate(id string) error {
	return s.store.DeleteTemplate(id)
}

// TransitionTemplate 执行模板状态机流转。
func (s *Service) TransitionTemplate(id, to string) (*model.Template, error) {
	t, err := s.store.GetTemplate(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionTemplate(t.Status, to) {
		return nil, model.NewValidationError("status", "模板状态不允许从 "+t.Status+" 流转到 "+to)
	}
	t.Status = to
	t.UpdatedAt = time.Now()
	if err := s.store.UpdateTemplate(t); err != nil {
		return nil, err
	}
	return t, nil
}

// normalizeVariables 去重并去除空白变量名。
func normalizeVariables(vars []string) []string {
	seen := make(map[string]bool)
	out := make([]string, 0, len(vars))
	for _, v := range vars {
		v = strings.TrimSpace(v)
		if v == "" || seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
	}
	return out
}
