# Changelog

## [0.0.25](https://github.com/jjdufu/kt-api/compare/v0.0.24...v0.0.25) (2026-11-12)


### Features

* **payment:** 易支付支持 RSA 签名方式 ([#8291](https://github.com/jjdufu/kt-api/issues/8291)) ([fd414d7](https://github.com/jjdufu/kt-api/commit/fd414d77ab9a51d41ae28b358a4123ab1417e2ea))


### Bug Fixes

* **channels:** 切换渠道时重置模型重定向同步状态 ([#8304](https://github.com/jjdufu/kt-api/issues/8304)) ([bb2420f](https://github.com/jjdufu/kt-api/commit/bb2420fa02a8be2ead38a5921dc97ac155dc81ae))

## [0.0.24](https://github.com/jjdufu/kt-api/compare/v0.0.23...v0.0.24) (2026-11-11)


### Features

* **channels:** 高级自定义渠道增加 DeepSeek 搜索与抓取路径选项 ([#8250](https://github.com/jjdufu/kt-api/issues/8250)) ([bace2ab](https://github.com/jjdufu/kt-api/commit/bace2ab76c135e421a1539f9a2ae1b7932c33828))
* **deepseek:** 补全 Responses 工具历史缺失的 reasoning ([#8162](https://github.com/jjdufu/kt-api/issues/8162)) ([267cb25](https://github.com/jjdufu/kt-api/commit/267cb256cbd755dba5eb9d43aa32279a4ff93d6f))
* **plugins:** 新增 gemini omni 1.2 flash 并把所有输入图片发送给 Vertex ([#8228](https://github.com/jjdufu/kt-api/issues/8228)) ([37c62f3](https://github.com/jjdufu/kt-api/commit/37c62f31c7717c22b6be338e49c82936829a0777))
* 支持用 HOST 环境变量设置监听地址 ([#8269](https://github.com/jjdufu/kt-api/issues/8269)) ([5793286](https://github.com/jjdufu/kt-api/commit/5793286e05d9b91f0b6247f5bd70b83fca0581c8))


### Bug Fixes

* **channel:** 所有 key 失效后重新启用某个 key 时恢复多 key 渠道 ([57a0eb8](https://github.com/jjdufu/kt-api/commit/57a0eb887b91e5c0b1938a2211c6684ac391fc8b))
* **channel:** 请求失败时保留手动禁用的渠道状态 ([e5604c3](https://github.com/jjdufu/kt-api/commit/e5604c3b341633a3734e1db6db2eecb2c674bd4b))
* **deepseek:** Chat Completions 中把 developer 角色映射为 system ([#8287](https://github.com/jjdufu/kt-api/issues/8287)) ([4951179](https://github.com/jjdufu/kt-api/commit/4951179fc142a593e22b077173b41b7b8908b97e))
* **oauth:** 创建 OAuth 用户时保存邀请人 ([#8247](https://github.com/jjdufu/kt-api/issues/8247)) ([cd9a341](https://github.com/jjdufu/kt-api/commit/cd9a3410d1d4034496408f6d23147ef712b61858))
* **sqlite:** SQLITE_PATH 为普通路径时保留默认 SQLite 参数 ([b29d339](https://github.com/jjdufu/kt-api/commit/b29d33931f36299fd635682d15b72877067532a8))
* **theme:** 认证页内联链接改用主题 link 色并避免专名断行 ([a080a77](https://github.com/jjdufu/kt-api/commit/a080a77fa9daa535a69edfaaace4894f5243e94b))
* **web:** 登录与注册页的法律条款文案支持翻译 ([#8254](https://github.com/jjdufu/kt-api/issues/8254)) ([2fe9917](https://github.com/jjdufu/kt-api/commit/2fe991770146f30c396d454ac6bf69d8447631cb))

## [0.0.23](https://github.com/jjdufu/kt-api/compare/v0.0.22...v0.0.23) (2026-11-07)


### Features

* 为 gpt-image-2.6 转发图像复杂度和标记参数 ([3dd6d88](https://github.com/jjdufu/kt-api/commit/3dd6d88bbdc43141564aeb1b8dca36d229a2acf8))


### Bug Fixes

* **auth:** 校验并修复 Redis 鉴权缓存完整性 ([a450257](https://github.com/jjdufu/kt-api/commit/a450257fc1d332b611f8fca4536dcb7908dc7333))
* **cache:** 修复用户和令牌鉴权缓存的 Redis 跨槽错误 ([755b351](https://github.com/jjdufu/kt-api/commit/755b351aaa851ae082855a65feae1153a159c139))
* **plugins:** 将上传插件同时挂载到顶层 /v1 等路径 ([4fbba52](https://github.com/jjdufu/kt-api/commit/4fbba52530b8eff86569e81f1014038b17539a40))
* **pricing:** 保存前逐项核对模型价格变更 ([8dea25a](https://github.com/jjdufu/kt-api/commit/8dea25a48cc9dea21f5e78eace96342d3604f173))
* **pricing:** 引导用户排查模型定价问题 ([4812ed5](https://github.com/jjdufu/kt-api/commit/4812ed5646e394d7b8ee2940728f91c4d9bff5fe))
* **theme:** 统一未保存价格状态的语义配色并遵循主题设置 ([5ed72b7](https://github.com/jjdufu/kt-api/commit/5ed72b73a5ac0215104e4edd2308c88ea4ffd02f))
* **vertex:** 按有效期缓存服务账号访问令牌 ([98cae58](https://github.com/jjdufu/kt-api/commit/98cae587606a3752fe3bd8d0d5cd89bbd2b592a0))
* 同步上游九项改进并适配开天主题与定价确认 ([b492165](https://github.com/jjdufu/kt-api/commit/b49216541b482ef583737ad520a649cdf0b69fa6))
* 明确会话缓存操作语义并支持按规则清理 ([c4868cc](https://github.com/jjdufu/kt-api/commit/c4868cc69737f25bc477d535774ece4e984a1fe4))

## [0.0.22](https://github.com/jjdufu/kt-api/compare/v0.0.21...v0.0.22) (2026-11-04)


### Features

* **channel:** 创建渠道时返回新 ID 及其与输入的对应关系 ([#8164](https://github.com/jjdufu/kt-api/issues/8164)) ([6960f9b](https://github.com/jjdufu/kt-api/commit/6960f9b1685b0c3b1ddc155a31a8d9109fcfbfbb))
* **gemini:** 新增 Argon 字段与推理参数映射 ([#8177](https://github.com/jjdufu/kt-api/issues/8177)) ([254f218](https://github.com/jjdufu/kt-api/commit/254f21866499c871abc7b67826fa2d4db6fac085))
* **plugins:** Gemini Omni 按 Vertex 返回的 token 计费 ([#8153](https://github.com/jjdufu/kt-api/issues/8153)) ([fe112e4](https://github.com/jjdufu/kt-api/commit/fe112e4c62ecdb2193704c394afbaf9749dee933))
* **plugins:** Vertex AI 渠道支持用 Gemini Omni 生成视频 ([#8142](https://github.com/jjdufu/kt-api/issues/8142)) ([eff4e1a](https://github.com/jjdufu/kt-api/commit/eff4e1a2c1de4e3b4b4ebe1c64ffe9887557f10a))
* **plugins:** 任务插件挂在厂商官方路径上并按厂商格式应答 ([#8117](https://github.com/jjdufu/kt-api/issues/8117)) ([ee3fe8a](https://github.com/jjdufu/kt-api/commit/ee3fe8a3702c307430f0a7db7b9884eff5ff16b4))
* **plugins:** 提示插件无法计费的任务模型 ([6ff33c4](https://github.com/jjdufu/kt-api/commit/6ff33c41f744e989368ab271367bb4a0867d6c8f))
* **relay:** 新增可选的 HTTP/2 健康探测 ([#8059](https://github.com/jjdufu/kt-api/issues/8059)) ([424d9d2](https://github.com/jjdufu/kt-api/commit/424d9d2a3bc9bd3557ab8673c19a89b40276fd34))
* **vertex:** 新增按渠道切换 v1beta1 端点的开关 ([#8179](https://github.com/jjdufu/kt-api/issues/8179)) ([0ab877d](https://github.com/jjdufu/kt-api/commit/0ab877d7864795575881d8cf15cadbbf73f00759))
* 模型性能统计纳入立即成功的任务 ([#8067](https://github.com/jjdufu/kt-api/issues/8067)) ([9dc7184](https://github.com/jjdufu/kt-api/commit/9dc7184d2da248efe25ca3085731e5c5ce4cc146))


### Bug Fixes

* **aws:** Bedrock Invoke 请求保留安全防护字段 ([#8071](https://github.com/jjdufu/kt-api/issues/8071)) ([82a10ba](https://github.com/jjdufu/kt-api/commit/82a10bafd1fdc826bc27c2a5619f8191920378c7))
* **billing:** 计费 prompt_tokens_details 中上报的 cache_creation_input_tokens ([#7909](https://github.com/jjdufu/kt-api/issues/7909)) ([e0f1699](https://github.com/jjdufu/kt-api/commit/e0f1699e4d2de108d76cad6e87f2f637c5f493d0))
* **cloudflare:** 保留 embedding 响应数据 ([#8041](https://github.com/jjdufu/kt-api/issues/8041)) ([ad13c5b](https://github.com/jjdufu/kt-api/commit/ad13c5b99324298d4f1c884b89bb07dece38f5af))
* **deps:** 升级 DOMPurify 到 3.4.17 ([#8116](https://github.com/jjdufu/kt-api/issues/8116)) ([eb5f457](https://github.com/jjdufu/kt-api/commit/eb5f4579a244791eb46d79f8d2c46acf393e81c0))
* **gemini:** chat 格式转换时保留顶层 imageConfig ([71aaec6](https://github.com/jjdufu/kt-api/commit/71aaec6ff74e19b4e4e8755b0b8fb44aa9737553))
* **gemini:** embedding 按上游 usageMetadata 计费 ([#8182](https://github.com/jjdufu/kt-api/issues/8182)) ([773e900](https://github.com/jjdufu/kt-api/commit/773e90055222b0cae1c9690d371cb76c21958f82))
* **logger:** 日志轮转后 gin 访问日志与 panic 恢复日志仍写入当前日志文件 ([#8098](https://github.com/jjdufu/kt-api/issues/8098)) ([222d261](https://github.com/jjdufu/kt-api/commit/222d26178c29e71b0ec4d01098b12ba100327bed))
* **plugins:** 加载自行声明 /{key} 地址的插件 ([#8119](https://github.com/jjdufu/kt-api/issues/8119)) ([7d7e9f2](https://github.com/jjdufu/kt-api/commit/7d7e9f2faa18480af25dd6ccd9bfa6a15da9641a))
* **plugins:** 可灵 4.0 时长限制在范围内并映射所有视频尺寸 ([#8118](https://github.com/jjdufu/kt-api/issues/8118)) ([33cb91f](https://github.com/jjdufu/kt-api/commit/33cb91fec1a83eaaf8e3209124f265f419e0c8c8))
* **plugins:** 可灵 4.0 现金计费优先使用可灵返回的 list_price ([#8125](https://github.com/jjdufu/kt-api/issues/8125)) ([e969c7f](https://github.com/jjdufu/kt-api/commit/e969c7f3fba7adee68a66686df34140cb6a47e82))
* **plugins:** 预扣或任务写入失败时仍对已转发的流计费 ([#8121](https://github.com/jjdufu/kt-api/issues/8121)) ([a14d8f9](https://github.com/jjdufu/kt-api/commit/a14d8f9964d748eb2486e8778696bc8db93234a9))
* **pricing:** 正确归类按请求动态计费的模型 ([#8010](https://github.com/jjdufu/kt-api/issues/8010)) ([97dc471](https://github.com/jjdufu/kt-api/commit/97dc4717627e5e6aca0cf05e5a119985888ff0ca))
* **profile:** 个人资料页遵循全局侧边栏显示设置 ([#7429](https://github.com/jjdufu/kt-api/issues/7429)) ([92b05f0](https://github.com/jjdufu/kt-api/commit/92b05f0950eb8bf2fbff094ab294aa37f53b93ea))
* **relay:** gpt-6.1-sol 与 gpt-6.1-terra 也应用 GPT-6 的 chat 参数规则 ([#8064](https://github.com/jjdufu/kt-api/issues/8064)) ([c1450bc](https://github.com/jjdufu/kt-api/commit/c1450bc5110abf2cb8004b334881c120bce7320d))
* **relaykit:** Claude 响应中保留工具参数里的数字 ([#7996](https://github.com/jjdufu/kt-api/issues/7996)) ([bc0a7c6](https://github.com/jjdufu/kt-api/commit/bc0a7c6dcd8c2acd7ca16c45c329c57cb6fa4c6a))
* **relaykit:** 保留与工具调用同时返回的 Gemini 文本 ([#7998](https://github.com/jjdufu/kt-api/issues/7998)) ([d3a307c](https://github.com/jjdufu/kt-api/commit/d3a307c5728c4be67f185d24caf97a3162f7dc92))
* **relaykit:** 修复 stop 既非字符串也非字符串数组时的协议转换问题 ([#7935](https://github.com/jjdufu/kt-api/issues/7935)) ([b5473d1](https://github.com/jjdufu/kt-api/commit/b5473d15eda4a340a9fad69abef5cfdd406df0c2))
* **relay:** 仅在遇到精确的裸 [DONE] 标记时结束流 ([#8074](https://github.com/jjdufu/kt-api/issues/8074)) ([24f3689](https://github.com/jjdufu/kt-api/commit/24f368973af1c16c7f0fefeaeb1efdc7a769fae0))
* **relay:** 允许 Azure 渠道开启 Responses WebSocket ([#8047](https://github.com/jjdufu/kt-api/issues/8047)) ([#8075](https://github.com/jjdufu/kt-api/issues/8075)) ([573c02a](https://github.com/jjdufu/kt-api/commit/573c02aa897ca5698a3a71b766c9c7e06285ccc8))
* **security:** 修复日志脱敏处理 URL 尾随标点的平方级开销 ([6677e01](https://github.com/jjdufu/kt-api/commit/6677e0109adfccc2055a88b19ed55904107715fe))
* **security:** 升级 @visactor/vchart 与 @visactor/react-vchart 到 2.1.7 ([4a98878](https://github.com/jjdufu/kt-api/commit/4a988782c5439101c2518cc22e84757ce26bfda3))
* **security:** 限制 zstd 解压窗口与请求体大小 ([700d94f](https://github.com/jjdufu/kt-api/commit/700d94fe0599e1261e37f172abc5d170533218a2))
* **task:** 任务超时前最多等待 72 小时 ([a4873cb](https://github.com/jjdufu/kt-api/commit/a4873cbfc099c8124f187d1adec4a5ba220b4370))
* **twofa:** 锁定到期后重新开始失败计数 ([#8172](https://github.com/jjdufu/kt-api/issues/8172)) ([cdfe6c5](https://github.com/jjdufu/kt-api/commit/cdfe6c546a7f61d10e1eb9f8fcc30ccb1be442b2))
* **xai:** 流式用量转换时保留 token 明细 ([#6420](https://github.com/jjdufu/kt-api/issues/6420)) ([4670297](https://github.com/jjdufu/kt-api/commit/4670297f75ca9be6019ac6a0bbe7467f712828dd))
* 上下文取消后仍统计已完成的非流式转发 ([#7844](https://github.com/jjdufu/kt-api/issues/7844)) ([60e8d1c](https://github.com/jjdufu/kt-api/commit/60e8d1c8c7a2f5a55f38e1ad79a3729204b7d495))
* 亲和渠道回退从最高优先级开始 ([#8078](https://github.com/jjdufu/kt-api/issues/8078)) ([b6df924](https://github.com/jjdufu/kt-api/commit/b6df9244f7238b2dc0a16fd3119dae23c0c023e1))
* 先保存 Turnstile 密钥再允许启用 ([#7960](https://github.com/jjdufu/kt-api/issues/7960)) ([0190632](https://github.com/jjdufu/kt-api/commit/0190632950fcc77aac1b2617c61c1e97a9025dc9))
* 退出登录时卸载已登录布局，停止无凭证的重复请求 ([#8110](https://github.com/jjdufu/kt-api/issues/8110)) ([f389819](https://github.com/jjdufu/kt-api/commit/f38981943c0291c8ad73448bb1a10ac39225b909))


### Performance Improvements

* **relay:** 参数覆盖中的正则不再每次请求都重新编译 ([#8061](https://github.com/jjdufu/kt-api/issues/8061)) ([14a2c9f](https://github.com/jjdufu/kt-api/commit/14a2c9ffec3c20929781f55fd627467e8cece9b1))
* 用 moego 读取和编辑请求 JSON 与 base64 ([#8148](https://github.com/jjdufu/kt-api/issues/8148)) ([d97a7e3](https://github.com/jjdufu/kt-api/commit/d97a7e37b99f35c22206f9cc74e1a5398596bbb8))

## [0.0.21](https://github.com/jjdufu/kt-api/compare/v0.0.20...v0.0.21) (2026-10-24)


### Bug Fixes

* **logging:** 访问日志中隐藏查询参数里的凭据 ([#7972](https://github.com/jjdufu/kt-api/issues/7972)) ([63bd564](https://github.com/jjdufu/kt-api/commit/63bd564c3a55a7432fb170c887fc015eefa5032b))
* **web:** 切换编辑抽屉记录时不再残留上一条的表单值 ([#7989](https://github.com/jjdufu/kt-api/issues/7989)) ([c54c049](https://github.com/jjdufu/kt-api/commit/c54c049ecdb4c08bbdf3fcf8a79c5ec259c8d221))
* 同步上游 new-api 两个修复（访问日志脱敏、编辑抽屉串值） ([7c4164f](https://github.com/jjdufu/kt-api/commit/7c4164f248403dad3e029ece215956d2864d0706))

## [0.0.20](https://github.com/jjdufu/kt-api/compare/v0.0.19...v0.0.20) (2026-10-23)


### Features

* **theme:** 开天漫画主题下禁用字体选项并提示 ([9190210](https://github.com/jjdufu/kt-api/commit/9190210dedd44f7033a806cb78d2855cbc404920))


### Bug Fixes

* **theme:** 开天漫画主题下圆角/密度设置生效，字体选项置灰 ([063bd52](https://github.com/jjdufu/kt-api/commit/063bd52815b1acfd53f8ff187636471f7f6990ea))
* **theme:** 开天漫画主题下圆角与密度设置生效 ([808d707](https://github.com/jjdufu/kt-api/commit/808d707063f1d318239dbf14eae71da5984c2da6))

## [0.0.19](https://github.com/jjdufu/kt-api/compare/v0.0.18...v0.0.19) (2026-10-23)


### Features

* 文档页外壳复用 marketing 组件 + Toast 关闭按钮 token 化 ([f67631f](https://github.com/jjdufu/kt-api/commit/f67631fba6b4062704a894e933b40a1663f186b9))

## [0.0.18](https://github.com/jjdufu/kt-api/compare/v0.0.17...v0.0.18) (2026-10-22)


### Features

* **theme:** 语义 token 层与 ui-bridge，统一下拉菜单样式 ([a5506a9](https://github.com/jjdufu/kt-api/commit/a5506a9a317387573cabd25d22d3588ace9a2cb3))

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
