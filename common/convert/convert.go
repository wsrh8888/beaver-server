/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * 中文：
 * 本文件为海狸 IM（Beaver IM）开源项目源代码。
 * 版权所有 © 2024-2026 Beaver IM Team，基于 MIT 协议授权。
 * 禁止删除、篡改或替换本文件头部版权与许可声明。
 * 使用与商业授权说明：https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * English:
 * This file is part of the Beaver IM open-source project.
 * Copyright (c) 2024-2026 Beaver IM Team. Licensed under the MIT License.
 * Do not remove, alter, or replace this copyright and license header.
 * Usage & commercial licensing: https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * beaver-server-header-v1
 */

// Package convert 在「同形不同型」的两套结构之间做字段映射转换。
//
// 用在哪：业务 API 的 HTTP 出入参（.api 生成的 types.XxxReq/XxxRes）与同域 RPC 的
// proto 出入参（pb.XxxReq/XxxRes）——同一份业务语义，两套 Go 类型。
//
// 映射契约（重要）：按 **Go 字段名** 对齐，不依赖 json tag。
//   - HTTP 层遵循 go-zero 规范：lowerCamel tag + 缩写全大写（UserID）。
//   - RPC 层遵循 protobuf 规范：snake_case tag + protoc 命名（UserId）。
//   - 两层 wire 格式本就应各自独立，内部转换不该假设两端 tag 一致。
//     字段名规范化（去下划线 + 转小写）后 UserID ≡ UserId ≡ user_id，
//     以稳定的 Go 标识符为契约，重命名可被编译器/IDE 追踪，杜绝 tag 漂移。
//
// 为什么不用 JSON 往返：那会让转换正确性依赖两端 json tag 碰巧相同，
// 任何一侧调整 wire 命名风格都会静默丢字段（编译期无感、线上才暴露）。
package convert

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// To 把 in 按 Go 字段名映射到 *T。
//
// 转换规则：
//   - 字段名按 normKey 对齐（忽略大小写/下划线）；目标无同名 field 则跳过。
//   - 基础类型可赋值则直接赋值，可转换（int→int32 等）则 Convert。
//   - struct / slice / array / map / pointer 递归映射。
//   - 仅 header tag 的字段（网关注入身份，非 payload）不参与映射。
//   - 类型不可兼容时返回 error，调用方应作为业务错误上抛，不要静默吞掉。
func To[T any](in any) (*T, error) {
	out := new(T)
	if in == nil {
		return out, nil
	}
	if err := assign(reflect.ValueOf(out).Elem(), reflect.ValueOf(in)); err != nil {
		return nil, err
	}
	return out, nil
}

// normKey 规范化字段名：去掉下划线/连字符并转小写。
// UserID / UserId / user_id / user-id → "userid"。
func normKey(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == '_' || r == '-' {
			continue
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

// deref 剥掉源值的指针层；nil 返回 invalid。
func deref(v reflect.Value) reflect.Value {
	for v.IsValid() && v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}
	return v
}

// assign 把 src 写入 dst（dst 必须是可设置的）。
func assign(dst, src reflect.Value) error {
	src = deref(src)
	if !src.IsValid() {
		return nil // 源为 nil，目标保持零值
	}
	if !dst.CanSet() {
		return nil
	}

	// 目标是指针：分配后递归写元素。
	if dst.Kind() == reflect.Pointer {
		elem := reflect.New(dst.Type().Elem())
		if err := assign(elem.Elem(), src); err != nil {
			return err
		}
		dst.Set(elem)
		return nil
	}

	// 数值/字符串/bool 等标量与 struct 分开处理：
	// 允许 int↔int32 这类宽度转换，但不要把标量硬塞进 struct 或反之。
	srcScalar := isScalar(src.Kind())
	dstScalar := isScalar(dst.Kind())
	if srcScalar || dstScalar {
		if !srcScalar || !dstScalar {
			return fmt.Errorf("convert: 类型不兼容 %s → %s", src.Type(), dst.Type())
		}
		switch {
		case src.Type().AssignableTo(dst.Type()):
			dst.Set(src)
		case src.Type().ConvertibleTo(dst.Type()):
			dst.Set(src.Convert(dst.Type()))
		default:
			return fmt.Errorf("convert: 无法转换 %s → %s", src.Type(), dst.Type())
		}
		return nil
	}

	switch dst.Kind() {
	case reflect.Struct:
		return assignStruct(dst, src)
	case reflect.Slice, reflect.Array:
		return assignSlice(dst, src)
	case reflect.Map:
		return assignMap(dst, src)
	default:
		if src.Type().AssignableTo(dst.Type()) {
			dst.Set(src)
			return nil
		}
		return fmt.Errorf("convert: 未支持的 kind %s (%s → %s)", dst.Kind(), src.Type(), dst.Type())
	}
}

func assignStruct(dst, src reflect.Value) error {
	src = deref(src)
	if !src.IsValid() || src.Kind() != reflect.Struct {
		return fmt.Errorf("convert: 期望 struct，实际 %s", src.Kind())
	}
	st := src.Type()
	for i := 0; i < st.NumField(); i++ {
		sf := st.Field(i)
		if sf.PkgPath != "" { // 非导出字段
			continue
		}
		// 仅 header 注解、无 json tag 的字段是网关注入的身份头，不是 payload。
		if sf.Tag.Get("header") != "" && sf.Tag.Get("json") == "" {
			continue
		}
		key := normKey(sf.Name)
		df := fieldByKey(dst, key)
		if !df.IsValid() {
			continue // 目标无此字段（两端本就可不同形）
		}
		if err := assign(df, src.Field(i)); err != nil {
			return fmt.Errorf("convert: 字段 %s: %w", sf.Name, err)
		}
	}
	return nil
}

// fieldByKey 按规范化字段名在 dst 中找可设置字段。
func fieldByKey(dst reflect.Value, key string) reflect.Value {
	dt := dst.Type()
	for i := 0; i < dt.NumField(); i++ {
		f := dt.Field(i)
		if f.PkgPath != "" {
			continue
		}
		if normKey(f.Name) == key {
			return dst.Field(i)
		}
	}
	return reflect.Value{}
}

func assignSlice(dst, src reflect.Value) error {
	src = deref(src)
	if !src.IsValid() || (src.Kind() != reflect.Slice && src.Kind() != reflect.Array) {
		return fmt.Errorf("convert: 期望 slice/array，实际 %s", src.Kind())
	}
	out := reflect.MakeSlice(reflect.SliceOf(dst.Type().Elem()), src.Len(), src.Len())
	for i := 0; i < src.Len(); i++ {
		if err := assign(out.Index(i), src.Index(i)); err != nil {
			return fmt.Errorf("convert: [%d]: %w", i, err)
		}
	}
	dst.Set(out)
	return nil
}

func assignMap(dst, src reflect.Value) error {
	src = deref(src)
	if !src.IsValid() || src.Kind() != reflect.Map {
		return fmt.Errorf("convert: 期望 map，实际 %s", src.Kind())
	}
	out := reflect.MakeMapWithSize(dst.Type(), src.Len())
	iter := src.MapRange()
	for iter.Next() {
		ev := reflect.New(dst.Type().Elem()).Elem()
		if err := assign(ev, iter.Value()); err != nil {
			return fmt.Errorf("convert: map[%v]: %w", iter.Key(), err)
		}
		out.SetMapIndex(iter.Key(), ev)
	}
	dst.Set(out)
	return nil
}

func isScalar(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	}
	return false
}
