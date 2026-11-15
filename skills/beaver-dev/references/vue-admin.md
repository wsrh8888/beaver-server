# Vue 后台系统（beaver-manager）

> 依据仓库现状整理（`package.json` / `src/` 结构）。此端暂无独立 `.cursor/rules`，遵循下面既有约定。

## 技术栈

- Vue `^3.5.17` + TypeScript `~5.8.3` + Vite `^7`。
- 状态：Pinia `3.0.3`；路由：vue-router `4.5.1`。
- UI：Element Plus `2.10.2` + `@element-plus/icons-vue`。
- 样式：Less `4.3.0`；HTTP：axios `1.10.0`。
- 自动导入：`unplugin-auto-import` + `unplugin-vue-components`；代码规范：`@antfu/eslint-config`。

## 目录结构（`src/`）

```
api/       接口封装
assets/    静态资源
config/    配置
layouts/   布局
pages/     页面
pinia/     Pinia store
router/    路由表
types/     类型定义
utils/     工具
App.vue / main.ts
```

## 脚本

```bash
npm run dev          # vite serve --mode=dev
npm run serve_prod   # vite serve --mode=prod
npm run serve_test   # vite serve --mode=test
npm run build_prod   # vite build --mode=prod
npm run build_test   # vite build --mode=test
npm run lint         # eslint . --fix
```

## 约定

- 新接口统一放 `api/` 封装后再在页面调用，不在页面里散写 axios。
- 全局状态走 `pinia/`，类型放 `types/`。
- 组件用 Element Plus 原生能力优先，少自造表单控件样式。
- 遵循全局约定：文件命名不以大写开头、多词驼峰；版权头与现有文件一致、完整保留。
