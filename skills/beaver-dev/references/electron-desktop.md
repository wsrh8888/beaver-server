# Electron 客户端（beaver-desktop）

> 来源：`.cursor/rules/rule.mdc`（PC 端部分）+ `.cursor/rules/beaver-desktop-ui.mdc`。

## 技术架构

Vite + Electron + Vue3。

### ts 别名

```js
commonModule:  path.resolve(__dirname, 'src/common'),
mainModule:    path.resolve(__dirname, 'src/main'),
renderModule:  path.resolve(__dirname, 'src/render'),
preloadModule: path.resolve(__dirname, 'src/preload'),
```

### 文件命名

- 所有文件不能以大写字母开头；多个单词用驼峰命名。

### CSS 编写

- 用 Less，采用嵌套结构。

## UI 组件（强制）

写 `src/render` 界面时：

- **优先使用** `src/render/components/ui` 下已有组件（Input、Button、Checkbox、Dialog、Tabs 等）。
- 禁止业务页手写原生 `<input>` / `<button>` / checkbox 再自己铺一套样式，除非 ui 目录确实没有对应能力。
- 密码类字段用 `BeaverInput`：`type="password"` + `show-password`（默认隐藏，右侧小眼睛切换显示）。
- 缺组件时先在 `components/ui` 补通用组件，再在业务页引用；不要在业务页复制一套交互样式。

## 版权头（强制）

新建或改文件时，顶部版权头必须与现有文件一致，完整保留 `beaver-desktop-header-v2` 整段（含中英文说明）。禁止缩短、删改、替换文件头版权与许可声明。版本标记：`beaver-desktop-header-v2`。
