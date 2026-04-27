# NoSocial 酒吧小程序 - 顾客端

## 项目概述

面向酒吧场景的微信小程序，使用 **TDesign 组件库** + **微信小程序原生开发**，采用暗色奢华风格设计，主色调为金色 `#D4AF37`。

## 功能模块

### 1. 首页 (`pages/index`)
- 轮播图（自动播放、点击跳转）
- 2x2 功能入口网格（共享股东、生日有礼、酒水塔罗、离店点评）
- 今日推荐酒水卡片
- 店铺信息（地址导航、电话拨打、营业时间、WiFi）

### 2. 展示页 (`pages/showcase`)
- 瀑布流画廊布局
- 图片懒加载 + 淡入动画
- 图片预览功能

### 3. 酒水单 (`pages/menu`)
- 横向滚动分类标签
- 酒水列表卡片（图片、价格、酒精度、配料标签）
- 搜索功能
- 上拉加载更多

### 4. 我的 (`pages/profile`)
- 微信登录 / 登出
- 用户信息展示
- 股东数据统计（累计收益、可提现余额、团队人数）
- 功能入口（优惠券、占卜历史、生日设置）

### 5. 共享股东 (`pages/shareholder`)
- **申请页**: 权益展示、费用说明、微信支付
- **收益页**: 收益概览、提现、收益明细列表
- **团队页**: 团队统计、成员列表

### 6. 酒水塔罗 (`pages/tarot`)
- **占卜页**: 问题输入、洗牌动画、抽牌翻转动画
- **结果页**: 牌面展示、含义解读、推荐酒水
- **历史页**: 过往占卜记录列表

### 7. 离店点评 (`pages/review`)
- **引导页**: 5步操作指引
- **提交页**: 截图上传、评价内容、提交审核

### 8. 生日有礼 (`pages/birthday`)
- 生日判断、礼品展示
- 礼品领取
- 历史记录

### 9. 优惠券 (`pages/coupon`)
- 优惠券列表（未使用/已使用/已过期）
- 券面信息展示

## 技术栈

| 技术 | 说明 |
|------|------|
| 微信小程序原生 | WXML + WXSS + JS |
| TDesign 组件库 | `tdesign-miniprogram` |
| MobX | `mobx-miniprogram` 状态管理 |
| Go + Gin | 后端 API（文档定义） |

## 项目结构

```
nosocial-wxapp/
├── app.js                    # 应用入口
├── app.json                  # 全局配置
├── app.wxss                  # 全局样式
├── pages/
│   ├── index/               # 首页
│   ├── showcase/            # 展示页
│   ├── menu/                # 酒水单
│   ├── profile/             # 我的
│   ├── drink/               # 酒水详情
│   ├── tarot/               # 塔罗占卜/结果/历史
│   ├── shareholder/         # 股东申请/收益/团队
│   ├── review/              # 点评引导/提交
│   ├── birthday/            # 生日礼品
│   ├── coupon/              # 优惠券
│   └── webview/             # 外部网页
├── components/              # 公共组件
├── utils/
│   ├── request.js           # 网络请求封装
│   └── constants.js         # 常量定义
├── stores/
│   └── app.js               # MobX 全局状态
└── assets/images/           # 图片资源
```

## 安装 & 使用

### 1. 安装依赖

```bash
npm install
```

### 2. 构建 TDesign

在微信开发者工具中点击：`工具 -> 构建 npm`

### 3. 配置项目

修改 `project.config.json` 中的 `appid` 为你自己的小程序 AppID。

### 4. 配置后端 API

修改 `utils/request.js` 中的 `BASE_URL` 为你的后端服务地址。

## API 接口说明

所有接口统一返回格式：`{ code: 0, data: {}, msg: "" }`

### 公共接口（无需登录）
- `GET /public/banners` - 轮播图列表
- `GET /public/drinks` - 酒水列表
- `GET /public/drinks/:id` - 酒水详情
- `GET /public/categories` - 分类列表

### 微信小程序接口（需登录）
- `POST /wx/login` - 微信登录
- `GET /wx/user` - 用户信息
- `PUT /wx/user/birthday` - 设置生日
- `POST /wx/shareholder/pay` - 股东注册支付
- `GET /shareholder/profile` - 股东信息
- `GET /shareholder/earnings` - 收益明细
- `GET /shareholder/team` - 团队列表
- `POST /wx/tarot/divine` - 塔罗占卜
- `GET /wx/tarot/history` - 占卜历史
- `GET /wx/birthday/gift` - 生日礼品
- `POST /wx/birthday/claim` - 领取礼品
- `POST /wx/review/submit` - 提交点评
- `GET /wx/coupons` - 优惠券列表

## 主题定制

全局 CSS 变量定义在 `app.wxss`：

```css
--nosocial-gold: #D4AF37;        /* 主色调 - 金色 */
--nosocial-black: #1a1a1a;       /* 背景色 - 深黑 */
--nosocial-black-light: #2a2a2a; /* 卡片背景 */
--nosocial-gray: #8a8a8a;        /* 次要文字 */
```

## 开发规范

1. **页面文件**: 每个页面包含 `.js` `.wxml` `.wxss` `.json` 四个文件
2. **样式命名**: 使用 BEM 命名规范，组件前缀为页面名
3. **状态管理**: 使用 MobX 管理全局状态（用户信息、登录状态）
4. **网络请求**: 统一通过 `utils/request.js` 封装
5. **图片加载**: 使用 `lazy-load` 属性实现懒加载

## 待完善项

1. [ ] TabBar 图标资源（`assets/images/tabbar/`）
2. [ ] 微信支付真实接口对接
3. [ ] 后端 API 联调
4. [ ] 性能优化（图片压缩、分包加载）
5. [ ] 单元测试
