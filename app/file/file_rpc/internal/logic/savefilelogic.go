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

package logic

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"beaver/app/file/file_models"
	"beaver/app/file/file_rpc/internal/svc"
	"beaver/app/file/file_rpc/types/file_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type SaveFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewSaveFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveFileLogic {
	return &SaveFileLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("save_file", ctx)}
}

func (l *SaveFileLogic) SaveFile(in *file_rpc.SaveFileReq) (*file_rpc.SaveFileRes, error) {
	var existing file_models.FileModel
	if err := l.svcCtx.DB.Take(&existing, "md5 = ?", in.Md5).Error; err == nil {
		return &file_rpc.SaveFileRes{FileKey: existing.FileKey}, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		l.logger.Error(model.LogMsg{
			Text: "按md5查重文件失败",
			Data: map[string]any{"md5": in.Md5, "err": err.Error()},
		})
		return nil, err
	}

	if in.FileInfoJson == "" {
		return nil, status.Error(codes.InvalidArgument, "fileInfo不能为空")
	}

	suffix := "jpg"
	if strings.Contains(in.OriginalName, ".") {
		parts := strings.Split(in.OriginalName, ".")
		if len(parts) > 1 {
			suffix = strings.ToLower(parts[len(parts)-1])
		}
	}
	fileKey := in.Md5 + "." + suffix

	source := file_models.QiniuSource
	if in.Source == "local" {
		source = file_models.LocalSource
	}

	fileInfo := &file_models.FileInfo{}
	if err := json.Unmarshal([]byte(in.FileInfoJson), fileInfo); err != nil {
		return nil, status.Error(codes.InvalidArgument, "fileInfo格式不正确")
	}

	newFile := &file_models.FileModel{
		FileKey:      fileKey,
		OriginalName: in.OriginalName,
		Size:         in.Size,
		Path:         in.Path,
		Md5:          in.Md5,
		Type:         in.Type,
		Source:       source,
		FileInfo:     fileInfo,
	}
	if err := l.svcCtx.DB.Create(newFile).Error; err != nil {
		l.logger.Error(model.LogMsg{
			Text: "保存文件失败",
			Data: map[string]any{"fileKey": fileKey, "md5": in.Md5, "err": err.Error()},
		})
		return nil, status.Error(codes.Internal, "保存文件失败")
	}

	l.logger.Info(model.LogMsg{
		Text: "保存文件成功",
		Data: map[string]interface{}{
			"fileKey": fileKey,
			"md5":     in.Md5,
			"size":    in.Size,
			"source":  source,
		},
	})

	return &file_rpc.SaveFileRes{FileKey: fileKey}, nil
}
