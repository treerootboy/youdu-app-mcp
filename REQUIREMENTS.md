# YouDu MCP 项目需求文档

## 项目概述

**项目名称**: YouDu IM MCP Server
**项目目标**: 为有度即时通讯(YouDu IM)提供Model Context Protocol (MCP)集成,使AI助手(如Claude)能够通过MCP协议调用有度IM的API
**技术栈**: Go 1.23+, MCP SDK, YouDu SDK

## 核心需求

### 1. 基础架构

#### 1.1 三层架构设计

```
用户层 (CLI/MCP/API)
    ↓
适配器层 (Adapter)
    ↓
SDK层 (YouDu Client)
```

#### 1.2架构设计思路
adapter 层用于将前端（CLI, MCP, API）的请求转换为对 Youdu-SDK 的 API 
调用。其设计要点如下：
 * 方法名简单化
   * 对前端暴露简洁、易于理解的方法名称。
 * 参数映射
   * 负责将前端传入的参数映射（mapping）为 Youdu-SDK 所需的参数格式。
 * 封装与适配
   * 由于 Youdu-SDK 的接口设计偏向底层，其参数设计不一定完全适合前端直接复用。
   * 因此，adapter 层需要进行适配，设计出更符合 CLI、MCP、API 调用需求的参数化接口，或对底层 SDK 进行封装。 

**要求**:
- 适配器层作为中间层,统一管理API定义
- 所有API只需在适配器层定义一次
- 自动生成CLI命令和MCP工具

#### 1.3 技术选型

**必须使用**:
- ✅ 官方MCP SDK: `github.com/modelcontextprotocol/go-sdk`
- ✅ 有度SDK: `github.com/addcnos/youdu/v2`
- ✅ CLI框架: `github.com/spf13/cobra`
- ✅ 配置管理: `github.com/spf13/viper`

**禁止**:
- ❌ 手动实现MCP协议 (必须使用官方SDK)
- ❌ 为每个API手写CLI命令 (必须自动生成)
- ❌ 为每个API手写MCP工具 (必须自动生成)


### 2. 需要实现的有度API

#### 2.1 部门管理 (Department)

| API名称 | SDK方法 | 描述 |
|---------|---------|------|
| get_dept_list | GetDeptList | 获取部门列表 |
| get_dept_user_list | GetDeptUserList | 获取部门用户列表 |
| get_dept_user_simple_list | GetDeptUserSimpleList | 获取部门用户简单列表 |
| create_dept | CreateDept | 创建部门 |
| update_dept | UpdateDept | 更新部门 |
| delete_dept | DeleteDept | 删除部门 |
| get_dept_alias_list | GetDeptAliasList | 获取部门别名列表 |
| get_dept_id_by_alias | GetDeptIDByAlias | 通过别名获取部门ID |

#### 2.2 用户管理 (User)

| API名称 | SDK方法 | 描述 |
|---------|---------|------|
| get_user | GetUser | 获取用户信息 |
| create_user | CreateUser | 创建用户 |
| update_user | UpdateUser | 更新用户 |
| delete_user | DeleteUser | 删除用户 |
| batch_delete_user | BatchDeleteUser | 批量删除用户 |
| get_user_enable_state | GetUserEnableState | 获取用户启用状态 |
| update_user_enable_state | UpdateUserEnableState | 更新用户启用状态 |
| update_user_position | UpdateUserPosition | 更新用户职位 |

#### 2.3 消息发送 (Message)

| API名称 | SDK方法 | 描述 |
|---------|---------|------|
| send_text_message | SendTextMessage | 发送文本消息 |
| send_image_message | SendImageMessage | 发送图片消息 |
| send_file_message | SendFileMessage | 发送文件消息 |
| send_link_message | SendLinkMessage | 发送链接消息 |
| send_mp_news_message | SendMpNewsMessage | 发送图文消息 |
| send_ex_link_message | SendExLinkMessage | 发送外链消息 |
| send_sys_message | SendSysMessage | 发送系统消息 |
| send_pop_window_message | SendPopWindowMessage | 发送弹窗消息 |
| send_message | SendMessage | 通用消息发送 |

#### 2.4 会话消息 (Session Message)

| API名称 | SDK方法 | 描述 |
|---------|---------|------|
| send_text_session_message | SendTextSessionMessage | 发送会话文本消息 |
| send_image_session_message | SendImageSessionMessage | 发送会话图片消息 |
| send_file_session_message | SendFileSessionMessage | 发送会话文件消息 |
| send_voice_session_message | SendVoiceSessionMessage | 发送会话语音消息 |
| send_video_session_message | SendVideoSessionMessage | 发送会话视频消息 |
| send_session_message | SendSessionMessage | 通用会话消息发送 |

#### 2.5 群组管理 (Group)

| API名称 | SDK方法 | 描述 |
|---------|---------|------|
| get_group_list | GetGroupList | 获取群组列表 |
| get_group_info | GetGroupInfo | 获取群组信息 |
| create_group | CreateGroup | 创建群组 |
| update_group | UpdateGroup | 更新群组 |
| delete_group | DeleteGroup | 删除群组 |
| add_group_member | AddGroupMember | 添加群组成员 |
| del_group_member | DelGroupMember | 删除群组成员 |
| is_group_member | IsGroupMember | 检查是否为群组成员 |

#### 2.6 会话管理 (Session)

| API名称 | SDK方法 | 描述 |
|---------|---------|------|
| create_session | CreateSession | 创建会话 |
| get_session | GetSession | 获取会话 |
| update_session | UpdateSession | 更新会话 |

#### 2.7 媒体管理 (Media)

| API名称 | SDK方法 | 描述 |
|---------|---------|------|
| upload_media | UploadMedia | 上传媒体 |
| get_media | GetMedia | 获取媒体 |
| search_media | SearchMedia | 搜索媒体 |

#### 2.8 认证 (Auth)

| API名称 | SDK方法 | 描述 |
|---------|---------|------|
| get_token | GetToken | 获取令牌 |
| identify | Identify | 身份验证 |