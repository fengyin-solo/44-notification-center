package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// cronBounds 五个字段的取值范围：分、时、日、月、周。
var cronBounds = [5][2]int{{0, 59}, {0, 23}, {1, 31}, {1, 12}, {0, 6}}

// cronNames 五个字段的名称，用于错误提示。
var cronNames = [5]string{"分钟", "小时", "日", "月", "星期"}

// ParseCron 解析并校验 5 字段 cron 表达式（分 时 日 月 周）。
func ParseCron(expr string) error {
	fields := strings.Fields(strings.TrimSpace(expr))
	if len(fields) != 5 {
		return NewValidationError("cron_expr", "Cron 表达式必须包含 5 个字段，用空格分隔")
	}
	for i, field := range fields {
		if err := validateCronField(field, cronBounds[i][0], cronBounds[i][1]); err != nil {
			return NewValidationError("cron_expr", cronNames[i]+"字段不合法: "+err.Error())
		}
	}
	return nil
}

// validateCronField 校验单个 cron 字段。
func validateCronField(field string, min, max int) error {
	if field == "" {
		return fmt.Errorf("字段不能为空")
	}
	if field == "*" {
		return nil
	}
	for _, part := range strings.Split(field, ",") {
		if part == "" {
			return fmt.Errorf("逗号分隔的子项不能为空")
		}
		if err := validateCronPart(part, min, max); err != nil {
			return err
		}
	}
	return nil
}

// validateCronPart 校验单个子项（支持 *、a-b、a-b/n、*/n、单值）。
func validateCronPart(part string, min, max int) error {
	step := 1
	base := part
	if idx := strings.IndexByte(part, '/'); idx >= 0 {
		base = part[:idx]
		stepStr := part[idx+1:]
		if stepStr == "" {
			return fmt.Errorf("步长不能为空")
		}
		n, err := strconv.Atoi(stepStr)
		if err != nil || n <= 0 {
			return fmt.Errorf("步长必须是正整数")
		}
		step = n
	}
	lo, hi := min, max
	if base != "*" {
		if idx := strings.IndexByte(base, '-'); idx >= 0 {
			a, err := strconv.Atoi(base[:idx])
			if err != nil {
				return fmt.Errorf("范围起点必须是整数")
			}
			b, err := strconv.Atoi(base[idx+1:])
			if err != nil {
				return fmt.Errorf("范围终点必须是整数")
			}
			if a > b {
				return fmt.Errorf("范围起点不能大于终点")
			}
			lo, hi = a, b
		} else {
			v, err := strconv.Atoi(base)
			if err != nil {
				return fmt.Errorf("必须是整数或 *")
			}
			lo, hi = v, v
		}
	}
	if lo < min || hi > max {
		return fmt.Errorf("数值必须在 %d-%d 之间", min, max)
	}
	if step > (hi - lo + 1) {
		return fmt.Errorf("步长不能大于取值区间")
	}
	return nil
}

// cronFieldSet 解析后的单个字段取值集合。
type cronFieldSet struct {
	min, max int
	values   map[int]bool
}

// parseCronField 解析单个 cron 字段为取值集合。
func parseCronField(field string, min, max int) (*cronFieldSet, error) {
	set := &cronFieldSet{min: min, max: max, values: make(map[int]bool)}
	if field == "*" {
		for v := min; v <= max; v++ {
			set.values[v] = true
		}
		return set, nil
	}
	for _, part := range strings.Split(field, ",") {
		step := 1
		base := part
		if idx := strings.IndexByte(part, '/'); idx >= 0 {
			base = part[:idx]
			step, _ = strconv.Atoi(part[idx+1:])
		}
		lo, hi := min, max
		if base != "*" {
			if idx := strings.IndexByte(base, '-'); idx >= 0 {
				lo, _ = strconv.Atoi(base[:idx])
				hi, _ = strconv.Atoi(base[idx+1:])
			} else {
				lo, _ = strconv.Atoi(base)
				hi = lo
			}
		}
		for v := lo; v <= hi; v += step {
			set.values[v] = true
		}
	}
	return set, nil
}

// match 判断字段值是否在允许集合内。
func (s *cronFieldSet) match(v int) bool {
	return s.values[v]
}

// NextRunTime 计算给定时间之后的下一次触发时间。
func NextRunTime(expr string, from time.Time) (time.Time, error) {
	if err := ParseCron(expr); err != nil {
		return time.Time{}, err
	}
	fields := strings.Fields(strings.TrimSpace(expr))
	sets := make([]*cronFieldSet, 5)
	for i, field := range fields {
		set, err := parseCronField(field, cronBounds[i][0], cronBounds[i][1])
		if err != nil {
			return time.Time{}, err
		}
		sets[i] = set
	}
	loc := from.Location()
	cur := time.Date(from.Year(), from.Month(), from.Day(), from.Hour(), from.Minute(), 0, 0, loc)
	cur = cur.Add(time.Minute)
	limit := cur.AddDate(1, 0, 0)
	for cur.Before(limit) {
		month := int(cur.Month())
		weekday := int(cur.Weekday())
		if sets[0].match(cur.Minute()) && sets[1].match(cur.Hour()) &&
			sets[2].match(cur.Day()) && sets[3].match(month) && sets[4].match(weekday) {
			return cur, nil
		}
		cur = cur.Add(time.Minute)
	}
	return time.Time{}, fmt.Errorf("一年内不存在匹配的触发时间")
}
