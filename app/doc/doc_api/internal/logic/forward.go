package logic

import (
	"context"
	"reflect"
)

// forward 是 doc_api 网关层的统一协议转换：HTTP DTO -> doc_rpc proto 请求 -> 调用 -> 回写 HTTP DTO。
//
// 契约约束：doc.api 的 json 名与 doc_rpc.proto 的 snake_case 字段经 goctl/protoc 生成后
// 都归一成同一套 PascalCase Go 字段名（spaceId / space_id -> SpaceId），因此这里按字段名映射，
// 不再为每个端点手写一遍字段搬运。唯一例外是鉴权主体：HTTP 侧统一叫 UserId（从
// Beaver-User-Id header 注入），而多数 RPC 请求里叫 OperatorUserId，由 alias 处理。
//
// 找不到对应字段时保持目标零值：允许 RPC 侧有 HTTP 不暴露的字段（内部字段、纪元等），
// 也允许 HTTP 侧有网关自用字段（分页默认值等），不产生运行时噪声。
func forward[APIReq any, APIRes any, PBReq any, PBRes any](
	ctx context.Context,
	apiReq *APIReq,
	call func(context.Context, *PBReq) (*PBRes, error),
) (*APIRes, error) {
	pbReq := new(PBReq)
	if err := mapStruct(apiReq, pbReq); err != nil {
		return nil, err
	}
	pbRes, err := call(ctx, pbReq)
	if err != nil {
		return nil, err
	}
	apiRes := new(APIRes)
	if pbRes != nil {
		if err := mapStruct(pbRes, apiRes); err != nil {
			return nil, err
		}
	}
	return apiRes, nil
}

// mapStruct 把 src 结构体的字段按名映射到 dst 结构体。两者都传指针。
func mapStruct(src any, dst any) error {
	sv := reflect.ValueOf(src)
	dv := reflect.ValueOf(dst)
	for sv.Kind() == reflect.Pointer {
		if sv.IsNil() {
			return nil
		}
		sv = sv.Elem()
	}
	for dv.Kind() == reflect.Pointer {
		if dv.IsNil() {
			return nil
		}
		dv = dv.Elem()
	}
	if sv.Kind() != reflect.Struct || dv.Kind() != reflect.Struct {
		return nil
	}
	return mapFields(dv, sv)
}

func mapFields(dst, src reflect.Value) error {
	t := dst.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		field := src.FieldByName(f.Name)
		if !field.IsValid() && f.Name == "OperatorUserId" {
			field = src.FieldByName("UserId")
		}
		if !field.IsValid() {
			continue
		}
		if err := assign(dst.Field(i), field); err != nil {
			return err
		}
	}
	return nil
}

// assign 做单个字段的值转换：支持数值宽度转换、指针/值结构体互转、切片元素递归转换。
func assign(dst, src reflect.Value) error {
	if !src.IsValid() {
		return nil
	}
	for src.Kind() == reflect.Pointer {
		if src.IsNil() {
			return nil
		}
		src = src.Elem()
	}
	for dst.Kind() == reflect.Pointer {
		if dst.IsNil() {
			if !dst.CanSet() {
				return nil
			}
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		dst = dst.Elem()
	}
	if !dst.CanSet() {
		return nil
	}
	if src.Type().AssignableTo(dst.Type()) {
		dst.Set(src)
		return nil
	}
	switch dst.Kind() {
	case reflect.Struct:
		if src.Kind() == reflect.Struct {
			return mapFields(dst, src)
		}
	case reflect.Slice:
		if src.Kind() == reflect.Slice {
			out := reflect.MakeSlice(dst.Type(), src.Len(), src.Len())
			for i := 0; i < src.Len(); i++ {
				if err := assign(out.Index(i), src.Index(i)); err != nil {
					return err
				}
			}
			dst.Set(out)
		}
	case reflect.String:
		if src.Kind() == reflect.String {
			dst.SetString(src.String())
		}
	case reflect.Bool:
		if src.Kind() == reflect.Bool {
			dst.SetBool(src.Bool())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		switch src.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			if src.Type().ConvertibleTo(dst.Type()) {
				dst.Set(src.Convert(dst.Type()))
			}
		}
	}
	return nil
}
