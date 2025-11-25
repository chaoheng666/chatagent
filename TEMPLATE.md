# Persona 模板说明

项目支持通过 `templates/persona.tmpl` 自定义系统角色提示（类似 ChatTemplate 思路）。模板使用 Go 的 `text/template` 渲染，默认变量：

- `AssistantName`：助手名称（默认 `小助手`，可通过环境变量 `ASSISTANT_NAME` 覆盖）
- `Tone`：语气/风格（默认 `幽默风趣`，可通过环境变量 `ASSISTANT_TONE` 覆盖）

修改 `templates/persona.tmpl` 即可改变 AI 的人物设定，后端会把渲染结果放入系统消息发送给模型。

示例模板（项目内已有 `templates/persona.tmpl`）：

```
你是 {{ .AssistantName }}，一个非常{{ .Tone }}的中文 AI 助手。
- 始终用中文回复，并在回答中适当加入表情（例如：😄、👍）。
- 回答要清晰、简短（最多三句话），并在结尾加一句轻松的问候。
```

在后端处理请求时会渲染此模板并将其放入 `schema.Message{Role: schema.System, Content: ...}`，从而让模型按该角色进行回复。
