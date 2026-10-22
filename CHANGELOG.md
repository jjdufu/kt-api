# Changelog

## [0.0.17](https://github.com/jjdufu/kt-api/compare/v0.0.16...v0.0.17) (2026-10-22)


### Features

* **billing:** OpenAI 12 小时提示缓存写入按独立价格计费 ([#7965](https://github.com/jjdufu/kt-api/issues/7965)) ([9f4eb8b](https://github.com/jjdufu/kt-api/commit/9f4eb8bd945d094a02c3cff1f5427df40f3b5a52))
* **email:** 新增 Resend 邮件服务商 ([#7880](https://github.com/jjdufu/kt-api/issues/7880)) ([f2214e2](https://github.com/jjdufu/kt-api/commit/f2214e2a4dae313376144c9b22560216a2f86275))
* **pricing:** 按站点汇率调整同步来的上游价格 ([e19e340](https://github.com/jjdufu/kt-api/commit/e19e3400d858107b87b8c619d33287d099a896c6))
* **task:** 原生插件路由接受 multipart 请求体 ([c260ae8](https://github.com/jjdufu/kt-api/commit/c260ae89674dbedacfb8c95c6db0c254c39f4d72))
* **task:** 向任务插件驱动传递客户端请求头 ([bc9fb0b](https://github.com/jjdufu/kt-api/commit/bc9fb0b6456ae63e6cfa9e7b9e59b014c55e7924))
* **task:** 普通渠道可在共享插件端点上参与竞争 ([0b77d79](https://github.com/jjdufu/kt-api/commit/0b77d79fa2a32538a484fba692eda2665bca9d70))
* **web:** 可从使用日志详情复制模型名 ([392b00e](https://github.com/jjdufu/kt-api/commit/392b00e510e6a0fbb053da1f5b79b6e9deaee9b3))
* **web:** 同步上游 0017 后的页面主题适配（默认 / 开天漫画） ([3294292](https://github.com/jjdufu/kt-api/commit/32942925b07a43601d9e21df9c5e2b42f0210b1f))


### Bug Fixes

* **billing:** 对未拆分 5m/1h 的 Claude 缓存写入计费 ([d3753e2](https://github.com/jjdufu/kt-api/commit/d3753e24328eae489fc7d3747f6685069c8f4876))
* **channels:** 保留高级自定义渠道的请求体透传设置 ([#7971](https://github.com/jjdufu/kt-api/issues/7971)) ([29c987f](https://github.com/jjdufu/kt-api/commit/29c987f5258136e93f65111348337a93606f4a4c))
* **relay:** Responses 流已送达时正常结束，不再误报断连错误 ([#7956](https://github.com/jjdufu/kt-api/issues/7956)) ([d9db1dd](https://github.com/jjdufu/kt-api/commit/d9db1dd1aba9840d184f8521405207129f422029))
* **relay:** 保持客户端 JSON Schema 的属性顺序 ([22180a3](https://github.com/jjdufu/kt-api/commit/22180a3669290dc1104a6bcd5b40d951736212eb))
* **relay:** 大的流式事件不再截断慢速客户端 ([5c9d51c](https://github.com/jjdufu/kt-api/commit/5c9d51c6d0996d5da92576ea58dad6fc7258a843))
* **task:** 由任务插件自行校验请求 ([7ee9f91](https://github.com/jjdufu/kt-api/commit/7ee9f9108f6ca5e61c4d8e48614978414b3016bc))
* 同步上游 new-api 0017（Resend 完整合入 + 16 个上游提交），适配两套主题页面 ([d271932](https://github.com/jjdufu/kt-api/commit/d271932d5009079a67ed2319598a9aea64da3361))
* 将 Go 升级到 1.26.10 以修复 CVE-2026-111943 ([#7942](https://github.com/jjdufu/kt-api/issues/7942)) ([3949c10](https://github.com/jjdufu/kt-api/commit/3949c107b1831125b6e5c44b10cc20056714a691))

## [0.0.16](https://github.com/jjdufu/kt-api/compare/v0.0.15...v0.0.16) (2026-10-20)


### Features

* **dashboard:** 在总 token 旁显示缓存用量 ([#7773](https://github.com/jjdufu/kt-api/issues/7773)) ([e07b5a0](https://github.com/jjdufu/kt-api/commit/e07b5a00e88f582ef4cb8cac1dddb719354674cc))
* 支持阿里多模态向量化模型 qwen3-vl-embedding ([#5179](https://github.com/jjdufu/kt-api/issues/5179)) ([39f43a6](https://github.com/jjdufu/kt-api/commit/39f43a63a19ac5d1e02118eb652f03fb82f46d94))


### Bug Fixes

* DeepSeek 美元钱包为零时回退到人民币余额 ([#7917](https://github.com/jjdufu/kt-api/issues/7917)) ([881610b](https://github.com/jjdufu/kt-api/commit/881610b4627346f8a65f400f1ca0349d4e819d92))
* **kling:** 插件升级到 1.1.1 ([8010cfd](https://github.com/jjdufu/kt-api/commit/8010cfd956c83d5f95554b8eb63bc4e1d3cb13ef))
* **kling:** 支持单个 API Key ([#7897](https://github.com/jjdufu/kt-api/issues/7897)) ([9bc29fb](https://github.com/jjdufu/kt-api/commit/9bc29fbd265f3d0615fa69189e015f226fbdbbcb))
* **limiter:** Redis 清空后恢复限流脚本 ([#6486](https://github.com/jjdufu/kt-api/issues/6486)) ([bc3f465](https://github.com/jjdufu/kt-api/commit/bc3f4655351b5c185af389c0898f41b23f5616e6))
* **relay:** Claude 流式请求在准备和重试期间持续保活 ([#7857](https://github.com/jjdufu/kt-api/issues/7857)) ([797ce9c](https://github.com/jjdufu/kt-api/commit/797ce9cc5dda7888341a0d96c9d68172b9bd46fc))
* **relaykit:** 为 Claude 合并相邻助手消息时保留工具调用 ([#7921](https://github.com/jjdufu/kt-api/issues/7921)) ([b311278](https://github.com/jjdufu/kt-api/commit/b31127896e6956dae001d268185d1a70b4a4802f))
* **task:** 修复任务提交与轮询的请求头覆盖 ([#7918](https://github.com/jjdufu/kt-api/issues/7918)) ([f31c966](https://github.com/jjdufu/kt-api/commit/f31c96664b0c5bf2dd8a1e76bb765d96677f3325))
* **task:** 查询重定向时仅向允许的主机发送凭据 ([c48c735](https://github.com/jjdufu/kt-api/commit/c48c7351accb90e73cb5bac5ff358dec086cf2d7))
* **task:** 每次重试重新计价，并在提交前预留额度 ([90f6c5c](https://github.com/jjdufu/kt-api/commit/90f6c5cca0431a0095afb9d4ebd851242e6093d5))
* 同步上游 new-api 13 个修复与功能（不含 Resend），适配缓存用量主题样式 ([56d3dbc](https://github.com/jjdufu/kt-api/commit/56d3dbc70ea5c9f91de0e2c1508c266f64df6571))

## [0.0.15](https://github.com/jjdufu/kt-api/compare/v0.0.14...v0.0.15) (2026-10-20)


### Features

* **home:** 首页、关于、文档改为网格底面 + 首屏赛博层 ([37fe000](https://github.com/jjdufu/kt-api/commit/37fe0003797fe8bd93320ef5bf5df4aa11ef4e4d))
* **i18n:** 首页漫画腔文案改为中性专业表述（7 种语言） ([4af5bb6](https://github.com/jjdufu/kt-api/commit/4af5bb6824162084708e9f3b776c315242f53982))
* **theme:** 删除简易主题，只保留默认与开天漫画 ([4b5ac97](https://github.com/jjdufu/kt-api/commit/4b5ac97330b5898e80f538316694793aa6f15327))
* **theme:** 开天漫画全站统一使用控制台字体 ([90123a0](https://github.com/jjdufu/kt-api/commit/90123a0621525abe80f95e2b1559b6f6a874aa57))
* **toast:** 提示移到右下角，成功提示 2 秒自动消失 ([7afdb74](https://github.com/jjdufu/kt-api/commit/7afdb748a2c71be734a0088788e7e5627f3bba63))
* 开天漫画统一字体、删除简易主题、首页背景重构与提示调整 ([a08c4e7](https://github.com/jjdufu/kt-api/commit/a08c4e78679351e1e763025b4ae9950247669b41))

## [0.0.14](https://github.com/jjdufu/kt-api/compare/v0.0.13...v0.0.14) (2026-10-20)


### Features

* **home:** 默认与简易主题的首页、关于页改为企业级 SaaS 布局 ([ce6585c](https://github.com/jjdufu/kt-api/commit/ce6585c8db65b84bcb726be598635e697ad235ca))
* **layout:** 控制台入口仅登录后显示，头像菜单新增控制台，控制台顶栏三栏布局 ([9c343fa](https://github.com/jjdufu/kt-api/commit/9c343fa39e555275884e77460f3e0a4bb14ec04b))
* **theme:** 开天漫画控制台改用 Chakra Petch + 系统字体，漫画字体仅用于首页/关于/登录 ([1318b54](https://github.com/jjdufu/kt-api/commit/1318b545deafead93d313bbe6cc95a0c7167fed4))
* **theme:** 恢复默认主题替代暗夜，超大字体简易改为简易 + 独立大字体开关 ([c77ae00](https://github.com/jjdufu/kt-api/commit/c77ae000c8841212e0c3dac2c4b34206e9213d11))
* 主题重整（默认 / 开天漫画 / 简易 + 大字体开关）与企业级首页 ([dede5a7](https://github.com/jjdufu/kt-api/commit/dede5a78aa0c4fb8ff65779ecadd848085f8e67c))


### Bug Fixes

* **home:** 首页与关于页各区块不再使用 Reveal 入场动画，数据区移除外层卡片 ([2256ea5](https://github.com/jjdufu/kt-api/commit/2256ea58c76b5795157d05cb3f3c3b9cc3c479a0))
* **layout:** 控制台顶栏导航不换行，居中搜索框按断点收窄 ([756d9af](https://github.com/jjdufu/kt-api/commit/756d9af982632968f4ff2ed0a47b6e5cefdf1d96))

## [0.0.13](https://github.com/jjdufu/kt-api/compare/v0.0.12...v0.0.13) (2026-10-20)


### Features

* **auth:** 登录页移除立绘，改为全屏赛博朋克场景 ([a2f3b45](https://github.com/jjdufu/kt-api/commit/a2f3b451bf76b8d294d0e2b35b73eef94e0d3206))
* **home:** 首页与关于页使用一张固定连续的页面背景 ([c42881d](https://github.com/jjdufu/kt-api/commit/c42881dfb109ad3eef5932e5e69081e62bc35133))
* **layout:** 统一 SiteHeader，控制台与公共页共用同一顶栏 ([b72a1a4](https://github.com/jjdufu/kt-api/commit/b72a1a42179c89cae0c426c6da4f596eecb6a469))
* 统一顶栏、连续页面背景、登录页场景与表格空状态居中 ([55bc4b6](https://github.com/jjdufu/kt-api/commit/55bc4b6c546f13b5e1b8a22dc8e62fbb42f9e92c))


### Bug Fixes

* **data-table:** 空状态改为表格外层非滚动容器中的绝对定位覆盖层 ([26e278c](https://github.com/jjdufu/kt-api/commit/26e278ceed87c0dd4bef21d8b79ad9b64fa3cb8d))

## [0.0.12](https://github.com/jjdufu/kt-api/compare/v0.0.11...v0.0.12) (2026-10-20)


### Features

* **home:** enterprise SaaS finish for home and about outside KT Comic ([fad6bf2](https://github.com/jjdufu/kt-api/commit/fad6bf25780e7f082435f28b4eeca2bc958a3ea8))
* **home:** one layout system for home and about ([be3a900](https://github.com/jjdufu/kt-api/commit/be3a900de7f886600e485129db2f91929f26d947))
* **layout:** full-width sticky public header on every public page ([93971eb](https://github.com/jjdufu/kt-api/commit/93971eb3a5ebcf45863c417460ed61ac26ab7681))


### Bug Fixes

* **auth:** one session-expired toast and a quiet sign-out ([d0f59b1](https://github.com/jjdufu/kt-api/commit/d0f59b135a12d4ecc3e40e70e5c0b5f5615835b7))
* **data-table:** centre the empty state in the visible table viewport ([6a11f14](https://github.com/jjdufu/kt-api/commit/6a11f1454a8ce3ec0f788c6699956c1a4b38ceea))

## [0.0.11](https://github.com/jjdufu/kt-api/compare/v0.0.10...v0.0.11) (2026-10-20)


### Features

* **brand:** monochrome currentColor logo outside KT Comic ([b1f15b8](https://github.com/jjdufu/kt-api/commit/b1f15b8d39310651b3d1c47c21678e4c34a636d2))
* **home:** clean theme-token variant of home, about, docs and auth pages ([d0f4195](https://github.com/jjdufu/kt-api/commit/d0f41955170afbd2fd6571123c9597be9a7ba719))
* **layout:** comic-styled public header navigation and icon buttons ([9db7b30](https://github.com/jjdufu/kt-api/commit/9db7b30b1cf181ff8395e68f0b7507f2003b5fe8))
* **theme:** full comic fonts for KT Comic via unicode-range slices ([f0fecb3](https://github.com/jjdufu/kt-api/commit/f0fecb3d89a78ce7e7a79ce99a0d62f8f0ffc003))
* **theme:** keep only underground, simple-large and kt-comic presets ([cebff1f](https://github.com/jjdufu/kt-api/commit/cebff1ff4edf4f2de7185a157f6f0c1588e4d339))
* **theme:** 三主题收敛、漫画全字体与页面跟随主题 ([abb272b](https://github.com/jjdufu/kt-api/commit/abb272b978f3aa79a7255146f8fc8f757edb016c))


### Bug Fixes

* **pricing:** single table outline and readable headers in model details ([53a07e7](https://github.com/jjdufu/kt-api/commit/53a07e7d729b3d040634b1c1ec127312003fa844))

## [0.0.10](https://github.com/jjdufu/kt-api/compare/v0.0.9...v0.0.10) (2026-10-20)


### Features

* **web:** add KT Comic (开天漫画) console theme preset ([9316775](https://github.com/jjdufu/kt-api/commit/93167754f95aacf340b78581c381b9479097d069))
* **web:** comic display fonts for console UI via unicode-range subsets ([36779e0](https://github.com/jjdufu/kt-api/commit/36779e06f80621812e60aab6cadcb5b36147da10))
* **web:** make KT Comic the default theme preset ([159c740](https://github.com/jjdufu/kt-api/commit/159c7401ac6effb7b66dc916117aa04a0e21200b))
* **web:** 控制台「开天漫画」主题预设（默认） ([cb9eafd](https://github.com/jjdufu/kt-api/commit/cb9eafd15a164299da1a6ec411259a3b59eca6c0))

## [0.0.9](https://github.com/jjdufu/kt-api/compare/v0.0.8...v0.0.9) (2026-10-20)


### Features

* **brand:** 开天大模型 / KTAPI logo and default site name ([a8d5fa9](https://github.com/jjdufu/kt-api/commit/a8d5fa91b23c1b7c42d57aad2108ab88e9cc4667))
* **web:** comic restyle for docs, about, auth pages and the user menu ([5188528](https://github.com/jjdufu/kt-api/commit/5188528ee24ed5811432d7185de07d45799f5833))
* 开天大模型品牌 + 全站漫画风格 + 前端性能优化 ([3d83264](https://github.com/jjdufu/kt-api/commit/3d83264c2b392b8d585eb4295d5860b4e714b03c))


### Bug Fixes

* **web:** seamless model marquee and sign-in first step on home page ([d16e248](https://github.com/jjdufu/kt-api/commit/d16e248a5458943d71b1c9babc6b9deb86b66fab))


### Performance Improvements

* **web:** lazy-load locales and slim the lobe icon registry ([de87e51](https://github.com/jjdufu/kt-api/commit/de87e51dc1e570db34bf6d279034d1a333c46e21))
* **web:** preload the home hero image and headline font ([965c8ac](https://github.com/jjdufu/kt-api/commit/965c8acb1aa9b069694a227a8f7b487b5a40fbfa))

## [0.0.8](https://github.com/jjdufu/kt-api/compare/v0.0.7...v0.0.8) (2026-10-20)


### Features

* **web:** comic cyberpunk home page ([44a1fad](https://github.com/jjdufu/kt-api/commit/44a1fad142c3223639b61bf0b1a20560705276fa))
* **web:** comic cyberpunk home page ([d419732](https://github.com/jjdufu/kt-api/commit/d419732af333e163f952b6f4b0aa2443900b259c))

## [0.0.7](https://github.com/jjdufu/kt-api/compare/v0.0.6...v0.0.7) (2026-10-19)


### Features

* require sign-in for the model square and announcements, remove rankings ([0675743](https://github.com/jjdufu/kt-api/commit/0675743ecb33414ebbf591d38893eb4e8989192f))
* **web:** follow the browser language and remove the language switcher ([c1551fd](https://github.com/jjdufu/kt-api/commit/c1551fd2bb2c05d20024a29090f838b2bb75bc0f))


### Bug Fixes

* **asset:** ListAssets 越权修复；模型广场/公告登录可见；删除排行榜；语言跟随浏览器 ([b0dc7f5](https://github.com/jjdufu/kt-api/commit/b0dc7f5549e65bcbd53c3386987cc6e7a53d0df2))
* **asset:** verify every ListAssets group id belongs to the caller ([32e6fbf](https://github.com/jjdufu/kt-api/commit/32e6fbf648d607683c229aa14cd1a91fac88156f))

## [0.0.6](https://github.com/jjdufu/kt-api/compare/v0.0.5...v0.0.6) (2026-10-19)


### Features

* **web:** add API docs pages ([7ffcf60](https://github.com/jjdufu/kt-api/commit/7ffcf6092ebffa4c882eec3cf4001f963a341f27))
* **web:** add API docs pages ([2956fc3](https://github.com/jjdufu/kt-api/commit/2956fc3af0e8c6ca68d6af4bd65b0ebdff560a77))

## [0.0.5](https://github.com/jjdufu/kt-api/compare/v0.0.4...v0.0.5) (2026-10-19)


### Features

* **channel:** read Seedance channel balance from New API compatible suppliers ([39d8ed6](https://github.com/jjdufu/kt-api/commit/39d8ed6e6eb928b014ffccfa9e9b0cb8fa373404))


### Bug Fixes

* **billing:** log task pre-deducts as pending settlement and stop assuming US dollars ([69875bc](https://github.com/jjdufu/kt-api/commit/69875bcaa22422e66358185c53c9439fa1353022))
* **seedance:** render a queued body before the first poll and send a queued callback on submit ([05b417b](https://github.com/jjdufu/kt-api/commit/05b417bd2e617b3b81239e96a91389b6fc87240a))
* **seedance:** 提交即返回 queued 并发送 queued 回调；任务预扣日志改为待结算；豆包视频渠道余额 ([247905a](https://github.com/jjdufu/kt-api/commit/247905a8fd8f6fdc9477a1ee59235cd78abe6052))

## [0.0.4](https://github.com/jjdufu/kt-api/compare/v0.0.3...v0.0.4) (2026-10-19)


### Features

* **web:** remove the upstream update checker and "update available" badge ([ab58cb4](https://github.com/jjdufu/kt-api/commit/ab58cb41344c00da16d3c630a31527038bf1a1de))


### Bug Fixes

* **relay:** stop sending upstream attribution headers to OpenRouter ([a8045a6](https://github.com/jjdufu/kt-api/commit/a8045a633ca7e62912b7242b35f940e84b0f6e80))
* **server:** keep the "panic detected" token in panic logs for alerting ([eff88c4](https://github.com/jjdufu/kt-api/commit/eff88c4e7fb62d152e5ad37c6781c0ebbfc3ebdf))
* **server:** log panics server-side and return a generic 500 ([929c333](https://github.com/jjdufu/kt-api/commit/929c3330c1be8f6f4399aa5b00a406ca726fb6c7))

## [0.0.3](https://github.com/jjdufu/kt-api/compare/v0.0.2...v0.0.3) (2026-10-19)


### Features

* **monitor:** add task callback metrics and alerts (step 16) ([c5b8eb4](https://github.com/jjdufu/kt-api/commit/c5b8eb453a49a7a2e66555fad27a5d13daaedeb6))
* Seedance task callbacks (callback_url) ([8884f22](https://github.com/jjdufu/kt-api/commit/8884f22667faa27168318d7ccc47a53d59a38d3e))
* **seedance:** accept callback_url on video task submit ([6b2dbb7](https://github.com/jjdufu/kt-api/commit/6b2dbb75b842abd0fb7c7e9cfbc6980afa4ed605))
* **task:** add optional HMAC signing for task callbacks ([798a9bd](https://github.com/jjdufu/kt-api/commit/798a9bd1283cdd787c61a329d1b44c360fffc67b))
* **task:** add task callback settings ([64e566e](https://github.com/jjdufu/kt-api/commit/64e566ea098fd75d3eb6651a2016e38fe1da3d6e))
* **task:** add task callback subscription and outbox models ([116dfab](https://github.com/jjdufu/kt-api/commit/116dfab1592b0c124725e63665b371a7b74ef9f5))
* **task:** add upstream callback fast path (step 13) ([9e43618](https://github.com/jjdufu/kt-api/commit/9e436188769a93621b1a35729ac6c9b62949e3c7))
* **task:** deliver task callbacks with retries and leases ([3500bda](https://github.com/jjdufu/kt-api/commit/3500bda43d734a873426b7edafd554203b3179fb))
* **task:** enqueue task callback events on status transitions ([9e6d316](https://github.com/jjdufu/kt-api/commit/9e6d3168c07e0fddc9aa804743ff0272f409c1f5))
* **task:** honor execution_expires_after for task timeout ([b62b551](https://github.com/jjdufu/kt-api/commit/b62b551b8d2a99748a338435ef2ad54fd4e62519))
* **task:** validate customer callback_url with SSRF protection ([5cc02f5](https://github.com/jjdufu/kt-api/commit/5cc02f57e4c1f22649cd655d4ec705c4ccd24759))
* **web:** show task callback deliveries and settings (step 15) ([d62acb5](https://github.com/jjdufu/kt-api/commit/d62acb556b3ff90bddbc54d49f494195ebd748c9))


### Bug Fixes

* **doubao:** treat supplier HTTP 400 failed-task body as failure ([6711efb](https://github.com/jjdufu/kt-api/commit/6711efb414bb0bc570ef4f237c14e7f1fb1789ed))
* **task:** keep task data consistent with locally decided terminal states ([dae9542](https://github.com/jjdufu/kt-api/commit/dae9542e78554a1631e8d2140bf0400038d03f9c))

## [0.0.2](https://github.com/jjdufu/kt-api/compare/v0.0.1...v0.0.2) (2026-10-19)


### Features

* Seedance asset relay and video path alias ([#10](https://github.com/jjdufu/kt-api/issues/10)) ([60f4dc5](https://github.com/jjdufu/kt-api/commit/60f4dc52ea63f1436441513adec9ef520f69916e))

## [0.0.1](https://github.com/jjdufu/kt-api/compare/v0.0.0...v0.0.1) (2026-10-19)


### Bug Fixes

* **ci:** skip duplicate test runs on release PRs and release merges ([b215bd2](https://github.com/jjdufu/kt-api/commit/b215bd201af7ca3dd086710a12a5083072472595))
* **ci:** skip duplicate test runs on release PRs and release merges ([1a3cbbc](https://github.com/jjdufu/kt-api/commit/1a3cbbccbb0821f571e054d66b18154389878d5b))

## 0.0.0 (2026-10-18)


### Features

* release pipeline with release-please and Docker Hub/ACR publish ([8e53c8f](https://github.com/jjdufu/kt-api/commit/8e53c8f5c197ce1ae294c7ebbefda16d9cd1d3a5))
* release pipeline with release-please and Docker Hub/ACR publish ([3481eb7](https://github.com/jjdufu/kt-api/commit/3481eb7b149f8d0aa755bb1a19075ef7e5262877))
