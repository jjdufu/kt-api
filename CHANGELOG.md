# Changelog

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
