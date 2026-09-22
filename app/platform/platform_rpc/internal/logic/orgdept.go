package logic

import (
	"errors"
	"strings"

	"beaver/app/platform/platform_models/org"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func loadDepartment(db *gorm.DB, deptID string) (*org_models.Department, error) {
	var dept org_models.Department
	err := db.Where("dept_id = ?", deptID).Take(&dept).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "部门不存在")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "查询部门失败")
	}
	return &dept, nil
}

func loadDepartments(db *gorm.DB) (map[string]org_models.Department, error) {
	var list []org_models.Department
	if err := db.Find(&list).Error; err != nil {
		return nil, status.Error(codes.Internal, "查询部门失败")
	}
	out := make(map[string]org_models.Department, len(list))
	for _, item := range list {
		out[item.DeptID] = item
	}
	return out, nil
}

// parentIsUnder 判断 parentID 是不是 deptID 自己或它的子孙。
func parentIsUnder(depts map[string]org_models.Department, deptID, parentID string) bool {
	seen := map[string]bool{}
	cur := parentID
	for cur != "" {
		if cur == deptID {
			return true
		}
		if seen[cur] {
			return true
		}
		seen[cur] = true
		parent, ok := depts[cur]
		if !ok {
			return false
		}
		cur = parent.ParentID
	}
	return false
}

func trimName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", status.Error(codes.InvalidArgument, "部门名不能为空")
	}
	return name, nil
}
