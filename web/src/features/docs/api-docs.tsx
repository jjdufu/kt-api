/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { Link } from '@tanstack/react-router'

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

import {
  BaseUrlBlock,
  BulletList,
  Block,
  CodeSample,
  DocsPageLayout,
  DocsSection,
  EndpointCard,
  EndpointOverviewTable,
  ErrorTable,
  FieldTable,
  KeyValueTable,
  MinorHeading,
  Note,
  Paragraph,
  ParamDetail,
  RecommendationList,
  ResponseFieldTable,
  RichText,
  StatusTable,
  SubHeading,
  type DocsSectionItem,
  type ErrorRow,
  type FieldRow,
  type ResponseFieldRow,
  type StatusRow,
} from './components/docs-primitives'
import { SEEDANCE_MODELS } from './constants'
import { useApiBaseUrl } from './hooks/use-api-base-url'

const SECTIONS: DocsSectionItem[] = [
  { id: 'base-url', title: 'BASE_URL' },
  { id: 'api-overview', title: '接口列表' },
  { id: 'create-video-task', title: '1. 创建视频生成任务' },
  { id: 'get-video-task', title: '2. 获取视频任务状态' },
  { id: 'create-asset-group', title: '3. 创建素材组' },
  { id: 'create-asset', title: '4. 创建素材资产' },
  { id: 'list-asset-groups', title: '5. 查询素材资产组合列表' },
  { id: 'get-asset-group', title: '6. 查询单个素材资产组合信息' },
  { id: 'list-assets', title: '7. 查询素材资产列表' },
  { id: 'get-asset', title: '8. 查询素材资产信息' },
  { id: 'update-asset-group', title: '9. 更新素材资产组合信息' },
  { id: 'update-asset', title: '10. 更新素材资产信息' },
  { id: 'delete-asset', title: '11. 删除单个素材资产' },
  { id: 'delete-asset-group', title: '12. 删除指定素材组' },
  { id: 'create-visual-validate-session', title: '13. 拉起真人认证 H5' },
  { id: 'get-visual-validate-result', title: '14. 获取真人认证结果' },
  { id: 'user-isolation', title: '用户隔离与 ID 规则' },
  { id: 'common-errors', title: '常见错误' },
  { id: 'usage-recommendations', title: '调用建议' },
]

const ASSET_ACTION_PATH = '/api/assets/?Action='
const ASSET_VERSION = '&Version=2024-01-01'

const ENDPOINTS = [
  {
    feature: '创建视频生成任务',
    method: 'POST',
    path: '/api/v3/contents/generations/tasks',
  },
  {
    feature: '获取视频任务状态',
    method: 'GET',
    path: '/api/v3/contents/generations/tasks/{task_id}',
  },
  ...[
    ['创建素材组', 'CreateAssetGroup'],
    ['创建素材资产', 'CreateAsset'],
    ['查询素材资产组合列表', 'ListAssetGroups'],
    ['查询单个素材资产组合信息', 'GetAssetGroup'],
    ['查询素材资产列表', 'ListAssets'],
    ['查询素材资产信息', 'GetAsset'],
    ['更新素材资产组合信息', 'UpdateAssetGroup'],
    ['更新素材资产信息', 'UpdateAsset'],
    ['删除单个素材资产', 'DeleteAsset'],
    ['删除指定素材组', 'DeleteAssetGroup'],
    ['拉起真人认证 H5', 'CreateVisualValidateSession'],
    ['获取真人认证结果', 'GetVisualValidateResult'],
  ].map(([feature, action]) => ({
    feature,
    method: 'POST',
    path: `${ASSET_ACTION_PATH}${action}${ASSET_VERSION}`,
  })),
]

function TimestampNote() {
  return (
    <Paragraph text='`CreateTime` / `UpdateTime` 由上游返回并原样透传，统一使用 UTC RFC3339，有精度时保留小数秒；`UpdateTime` 包含素材组成功修改的时间，状态和错误未变化的轮询不会推进 `UpdateTime`；未知时间会省略字段，不返回空字符串或伪造时间。' />
  )
}

function AssetDetailNote() {
  return (
    <>
      <Paragraph text='所有状态下 `AssetType` 统一返回 `Image`、`Video` 或 `Audio`。素材详情与列表在没有审核策略时省略 `Moderation`，错误码 `Code` 和消息 `Message` 均为空时省略 `Error`。' />
      <TimestampNote />
    </>
  )
}

const createTaskFields: FieldRow[] = [
  {
    name: 'model',
    type: 'string',
    required: '是',
    description: `非空；当前可用：${SEEDANCE_MODELS.map((model) => `\`${model}\``).join('、')}。具体可用模型以「模型广场」为准。`,
  },
  {
    name: 'content',
    type: 'array',
    required: '是',
    description: '非空数组；内容项的语义和模型兼容性由模型能力决定。',
  },
  {
    name: 'content[].type',
    type: 'string',
    required: '是',
    description:
      '可使用 `text`、`image_url`、`audio_url`、`video_url` 或 `draft_task`。',
  },
  {
    name: 'content[].text',
    type: 'string',
    required: '条件必填',
    description: '`type=text` 时传入描述文本。',
  },
  {
    name: 'content[].image_url.url',
    type: 'string',
    required: '条件必填',
    description: '`type=image_url` 时传入图片 URL 或 `asset://asset_...`。',
  },
  {
    name: 'content[].audio_url.url',
    type: 'string',
    required: '条件必填',
    description: '`type=audio_url` 时传入音频 URL 或 `asset://asset_...`。',
  },
  {
    name: 'content[].video_url.url',
    type: 'string',
    required: '条件必填',
    description: '`type=video_url` 时传入视频 URL 或 `asset://asset_...`。',
  },
  {
    name: 'content[].draft_task.id',
    type: 'string',
    required: '条件必填',
    description:
      '`type=draft_task` 时传入当前账号已成功完成的 Draft 视频任务 ID（即创建 Draft 任务时本平台返回的 `id`）；仅用于基于 Draft 视频生成正式视频。',
  },
  {
    name: 'content[].role',
    type: 'string',
    required: '否',
    description: '内容项角色。',
  },
  {
    name: 'callback_url',
    type: 'string',
    required: '否',
    description:
      '任务状态变更后由本平台推送回调的地址；不合法时提交直接返回 HTTP 400，不扣费、不创建任务。',
  },
  {
    name: 'return_last_frame',
    type: 'boolean',
    required: '否',
    description: '是否请求返回最后一帧。',
  },
  {
    name: 'service_tier',
    type: 'string',
    required: '否',
    description: '服务等级。',
  },
  {
    name: 'execution_expires_after',
    type: 'integer',
    required: '否',
    description: '任务执行过期时间（秒），范围 3600 至 259200，默认 172800。',
  },
  {
    name: 'ratio',
    type: 'string',
    required: '否',
    description: '省略时由所选模型使用默认行为；支持范围由所选模型决定。',
  },
  {
    name: 'duration',
    type: 'integer',
    required: '否',
    description: '省略时由所选模型使用默认行为；具体时长范围由所选模型决定。',
  },
  {
    name: 'resolution',
    type: 'string',
    required: '否',
    description: '默认720p；支持范围由所选模型决定。',
  },
  {
    name: 'generate_audio',
    type: 'boolean',
    required: '否',
    description: '省略时使用模型默认行为；显式 `false` 会被保留。',
  },
  {
    name: 'omni_reference_task_type',
    type: 'string',
    required: '否',
    description:
      '仅 Seedance 2.5 的全模态参考任务可用：`auto`、`reference`、`edit`、`extend`。`edit` 和 `extend` 会触发上游同步约束校验。',
  },
  {
    name: 'output_format',
    type: 'string',
    required: '否',
    description:
      '输出封装格式；仅 Seedance 2.5 支持，可选 `mp4`、`mov`，默认 `mp4`。',
  },
  {
    name: 'draft',
    type: 'boolean',
    required: '否',
    description:
      'Seedance 2.5 Draft 模式开关；创建样片时必须为 `true`，基于样片创建正式视频时可省略或设为 `false`。',
  },
  {
    name: 'tools',
    type: 'array',
    required: '否',
    description: '工具配置数组。',
  },
  {
    name: 'tools[].type',
    type: 'string',
    required: '否',
    description: '工具类型。',
  },
  {
    name: 'safety_identifier',
    type: 'string',
    required: '否',
    description: '安全标识。',
  },
  {
    name: 'priority',
    type: 'integer',
    required: '否',
    description: '任务优先级。',
  },
  {
    name: 'frames',
    type: 'integer',
    required: '否',
    description: '生成帧数。',
  },
  {
    name: 'seed',
    type: 'integer',
    required: '否',
    description: '随机种子。',
  },
  {
    name: 'camera_fixed',
    type: 'boolean',
    required: '否',
    description: '是否固定镜头。',
  },
  {
    name: 'watermark',
    type: 'boolean',
    required: '否',
    description: '省略时使用模型默认行为；显式 `false` 会被保留。',
  },
]

const createTaskSuccessJson = `{
  "id": "task_BGJr5gwbfZssKxUFWXqLgOt2Dek4lmO5"
}`

const createTaskErrorJson = `{
  "error": {
    "code": "invalid_request",
    "message": "content is required"
  }
}`

function CreateVideoTask() {
  const baseUrl = useApiBaseUrl()
  return (
    <EndpointCard
      anchor='create-video-task-endpoint'
      number='1'
      title='创建视频生成任务'
      method='POST'
      path='/api/v3/contents/generations/tasks'
      contentType='application/json'
    >
      <Block>
        <SubHeading text='请求示例' />
        <CodeSample
          language='bash'
          code={`curl --request POST \\
  "${baseUrl}/api/v3/contents/generations/tasks" \\
  --header "Authorization: Bearer $KT_API_KEY" \\
  --header "Content-Type: application/json" \\
  --data '{
    "model": "doubao-seedance-2-5-260628",
    "content": [
      {
        "type": "text",
        "text": "一只猫在阳光下的花园中奔跑"
      }
    ],
    "generate_audio": true,
    "omni_reference_task_type": "auto",
    "ratio": "16:9",
    "duration": 8,
    "resolution": "720p",
    "output_format": "mov",
    "watermark": false
  }'`}
        />
      </Block>
      <Block>
        <SubHeading text='请求字段' />
        <FieldTable rows={createTaskFields} />
      </Block>
      <ParamDetail name='model'>
        <Paragraph text='必填字符串。填写本平台已开通的模型 ID（见上表或「模型广场」）；不支持直接填写火山方舟 Endpoint ID。模型无可用渠道时返回 HTTP 400。' />
      </ParamDetail>
      <ParamDetail name='content'>
        <Paragraph text='必填对象数组。可组合文本、图片、视频和音频；也可传入样片任务。常见组合如下：' />
        <ul className='text-muted-foreground list-disc space-y-1 pl-6 text-sm leading-relaxed'>
          <li>仅文本。</li>
          <li>文本（可选）与图片。</li>
          <li>文本（可选）与视频。</li>
          <li>文本（可选）与音频；仅 Seedance 2.5 允许没有图片或视频。</li>
          <li>文本（可选）与图片、视频、音频的任意受支持组合。</li>
          <li>
            <RichText text='单个 `draft_task`，用于基于已完成的 Draft 样片创建正式视频；不能与文本或参考素材混用。' />
          </li>
        </ul>
        <Note text='首帧/首尾帧、参考图片以及多模态参考输入是互斥任务类型，不能在同一个任务中混用。Seedance 2.5 和 2.0 系列对含真人面孔的输入素材存在额外限制，应使用平台支持或已授权的素材方案。' />
        <ParamDetail name='omni_reference_task_type'>
          <Paragraph text='仅 Seedance 2.5 的全模态参考生视频任务可用。可省略，或设置为 `auto`、`reference`、`edit`、`extend`：' />
          <ul className='text-muted-foreground list-disc space-y-1 pl-6 text-sm leading-relaxed'>
            <li>
              <RichText text='`auto`：让模型根据提示词和参考素材自动识别参考生视频、视频编辑或视频延长任务。' />
            </li>
            <li>
              <RichText text='`reference`：参考生视频任务，基于参考素材生成新视频；`ratio` 和 `duration` 无特殊限制。' />
            </li>
            <li>
              <RichText text='`edit`：必须包含 `reference_video`，且 `ratio` 为 `adaptive`、`duration` 为 `-1`。' />
            </li>
            <li>
              <RichText text='`extend`：必须包含 `reference_video`，且 `ratio` 为 `adaptive`。' />
            </li>
          </ul>
          <Paragraph text='该字段由上游校验并返回原始错误结果。显式指定 `edit` 或 `extend` 时，上游会在提交阶段校验上述约束；提示词意图与指定类型不一致仍可能导致异步任务失败。' />
        </ParamDetail>
        <Block>
          <MinorHeading>文本元素</MinorHeading>
          <FieldTable
            rows={[
              {
                name: 'content[].type',
                type: 'string',
                required: '是',
                description: '固定 `text`；元素类型。',
              },
              {
                name: 'content[].text',
                type: 'string',
                required: '是',
                description:
                  '中文建议不超过 500 字；英文建议不超过 1000 词；描述要生成的视频。',
              },
            ]}
          />
          <Paragraph text='所有模型支持中文和英文提示词。Seedance 2.5 还支持西班牙语、印度尼西亚语、葡萄牙语、日语、马来语、泰语、阿拉伯语、越南语和韩语；Seedance 2.0 系列还支持西班牙语、印度尼西亚语、葡萄牙语和日语。' />
        </Block>
        <Block>
          <MinorHeading>图片元素</MinorHeading>
          <FieldTable
            rows={[
              {
                name: 'content[].type',
                type: 'string',
                required: '是',
                description: '固定 `image_url`；元素类型。',
              },
              {
                name: 'content[].image_url',
                type: 'object',
                required: '是',
                description: '图片对象。',
              },
              {
                name: 'content[].image_url.url',
                type: 'string',
                required: '是',
                description: '公网 URL、或 `asset://<ASSET_ID>`。',
              },
              {
                name: 'content[].role',
                type: 'string',
                required: '视任务而定',
                description:
                  '`first_frame`、`last_frame`、`reference_image`；图片用途。',
              },
            ]}
          />
          <MinorHeading>图片限制</MinorHeading>
          <KeyValueTable
            kLabel='项目'
            vLabel='限制'
            rows={[
              {
                k: '格式',
                v: 'jpeg、png、webp、bmp、tiff、gif；Seedance 1.5 Pro 及更新模型还支持 heic、heif',
              },
              {
                k: '宽高比',
                v: '0.4 至 2.5',
              },
              {
                k: '宽和高',
                v: '各为 300 至 6000 px',
              },
              {
                k: '单张大小',
                v: '小于 30 MB',
              },
              {
                k: '请求体大小',
                v: '不超过 64 MB',
              },
              {
                k: '首帧任务',
                v: '1 张；role 可省略或设为 first_frame',
              },
              {
                k: '首尾帧任务',
                v: '2 张；分别设为 first_frame 和 last_frame',
              },
              {
                k: 'Seedance 2.5 多模态参考',
                v: '1 至 30 张；role 为 reference_image',
              },
              {
                k: 'Seedance 2.0 系列多模态参考',
                v: '1 至 9 张；role 为 reference_image',
              },
            ]}
          />
          <Paragraph text='首尾帧图片允许相同。如果两张图片的宽高比不同，系统以首帧为基准处理尾帧。' />
        </Block>
        <Block>
          <MinorHeading>视频元素</MinorHeading>
          <FieldTable
            rows={[
              {
                name: 'content[].type',
                type: 'string',
                required: '是',
                description: '固定 `video_url`；元素类型。',
              },
              {
                name: 'content[].video_url',
                type: 'object',
                required: '是',
                description: '视频对象。',
              },
              {
                name: 'content[].video_url.url',
                type: 'string',
                required: '是',
                description: '公网 URL 或 `asset://<ASSET_ID>`。',
              },
              {
                name: 'content[].role',
                type: 'string',
                required: '是',
                description: '固定 `reference_video`；视频用途。',
              },
            ]}
          />
          <MinorHeading>视频限制</MinorHeading>
          <KeyValueTable
            kLabel='项目'
            vLabel='限制'
            rows={[
              {
                k: '容器',
                v: 'MP4（.mp4、video/mp4）或 QuickTime（.mov、video/quicktime）',
              },
              {
                k: '视频编码',
                v: 'H.264/AVC 或 H.265/HEVC',
              },
              {
                k: '音频编码',
                v: 'AAC 或 MP3',
              },
              {
                k: '分辨率',
                v: '480p、720p、1080p、4K',
              },
              {
                k: '宽高比',
                v: '0.4 至 2.5',
              },
              {
                k: '宽和高',
                v: '各为 300 至 6000 px',
              },
              {
                k: '总像素',
                v: '409600 至 8295044',
              },
              {
                k: '单个文件大小',
                v: '不超过 200 MB',
              },
              {
                k: '帧率',
                v: '24 至 60 FPS',
              },
              {
                k: 'Seedance 2.5',
                v: '单个 2 至 30 秒；最多 10 个；总时长不超过 30 秒',
              },
              {
                k: 'Seedance 2.0 系列',
                v: '单个 2 至 15 秒；最多 3 个；总时长不超过 15 秒',
              },
            ]}
          />
        </Block>
        <Block>
          <MinorHeading>音频元素</MinorHeading>
          <FieldTable
            rows={[
              {
                name: 'content[].type',
                type: 'string',
                required: '是',
                description: '固定 `audio_url`；元素类型。',
              },
              {
                name: 'content[].audio_url',
                type: 'object',
                required: '是',
                description: '音频对象。',
              },
              {
                name: 'content[].audio_url.url',
                type: 'string',
                required: '是',
                description:
                  '公网 URL、Base64 Data URL 或 `asset://<ASSET_ID>`。',
              },
              {
                name: 'content[].role',
                type: 'string',
                required: '是',
                description: '固定 `reference_audio`；音频用途。',
              },
            ]}
          />
          <Paragraph text='Base64 格式为 `data:audio/<format>;base64,<data>`，格式名使用小写。' />
          <MinorHeading>音频限制</MinorHeading>
          <KeyValueTable
            kLabel='项目'
            vLabel='限制'
            rows={[
              {
                k: '格式',
                v: 'wav、mp3',
              },
              {
                k: '单个文件大小',
                v: '不超过 15 MB',
              },
              {
                k: '请求体大小',
                v: '不超过 64 MB',
              },
              {
                k: 'Seedance 2.5',
                v: '单段 2 至 30 秒；最多 10 段；总时长不超过 30 秒；允许只传音频',
              },
              {
                k: 'Seedance 2.0 系列',
                v: '单段 2 至 15 秒；最多 3 段；总时长不超过 15 秒；至少同时包含一张图片或一个视频',
              },
            ]}
          />
        </Block>
        <Block>
          <MinorHeading>Draft 任务元素</MinorHeading>
          <Paragraph text='仅在基于 Seedance 2.5 Draft 样片生成正式视频时使用。正式视频请求的 `content` 必须且只能包含该元素，不能再传入文本、图片、视频或音频。' />
          <FieldTable
            rows={[
              {
                name: 'content[].type',
                type: 'string',
                required: '是',
                description: '固定 `draft_task`；元素类型。',
              },
              {
                name: 'content[].draft_task',
                type: 'object',
                required: '是',
                description: 'Draft 任务对象。',
              },
              {
                name: 'content[].draft_task.id',
                type: 'string',
                required: '是',
                description:
                  '当前账号已成功完成的 Draft 视频任务 ID（本平台返回的 `task_` 开头 ID，由本平台换算为上游任务 ID）；自 Draft 任务 `created_at` 起有效 7 天。',
              },
            ]}
          />
        </Block>
      </ParamDetail>
      <ParamDetail name='callback_url'>
        <Paragraph text='可选字符串。任务状态改变时，由本平台（而非上游）向该地址发送 `POST` 请求，`Content-Type: application/json`；回调 body 与同一时刻查询任务接口返回的响应体逐字节一致。`callback_url` 不会转发给上游，也不会出现在查询结果中。会推送的 `status`：' />
        <StatusTable
          rows={[
            { name: 'queued', description: '排队中（提交成功后立即推送）' },
            { name: 'running', description: '运行中' },
            { name: 'succeeded', description: '已成功' },
            {
              name: 'failed',
              description: '已失败，body 中带 `error.code` 与 `error.message`',
            },
            {
              name: 'expired',
              description: '超过任务过期阈值（含本平台判定的超时）',
            },
          ]}
        />
        <BulletList
          items={[
            '`cancelled` 不推送回调，最终状态以查询接口为准。',
            '收到任意 2xx 响应即视为送达，响应体会被忽略；单次请求超时 5 秒；3xx 不跟随，按失败处理。',
            '`succeeded` / `failed` 最多尝试 4 次（首次 + 3 次重试），间隔约 5 秒；`queued` / `running` / `expired` 只发送 1 次。',
            '同一任务的事件按顺序投递；若更新的状态已产生，尚未送达的 `queued` / `running` 事件会被跳过，终态一定会投递。',
            '回调请求带 `X-KtApi-Event-Id`（事件唯一 ID，重试时不变，请据此去重）、`X-KtApi-Delivery-Attempt`（第几次尝试）和 `X-KtApi-Task-Id`（任务 ID）请求头；投递语义为至少一次，接收端需要幂等。',
            '管理员开启回调签名后，请求还会带 `X-KtApi-Timestamp` 与 `X-KtApi-Signature: v1=<hex>`，其中 `<hex> = HMAC-SHA256(用户密钥, timestamp + "." + event_id + "." + 原始 body)`。',
            '`callback_url` 需为可公网访问的 `http`/`https` 地址（默认端口限 80、443、8080、8443，长度不超过 2048），不能带用户名密码，不能解析到私网、回环或链路本地地址；不合法时提交返回 HTTP 400，不扣费、不创建任务。',
          ]}
        />
      </ParamDetail>
      <ParamDetail name='camera_fixed'>
        <Paragraph text='可选布尔值，默认 `false`。设为 `true` 时，请求固定摄像头效果，但结果不作保证；参考图任务不支持该字段。适用于 Seedance 1.5 Pro、1.0 Pro 和 1.0 Pro Fast。' />
      </ParamDetail>
      <ParamDetail name='draft'>
        <Paragraph text='Seedance 2.5 Draft 模式先生成低分辨率样片，确认效果后再基于该样片生成正式视频。Draft 任务与正式视频任务均通过 `POST /api/v3/contents/generations/tasks` 创建，并使用同一查询接口获取状态和 `content.video_url`。' />
        <Block>
          <MinorHeading>创建 Draft 样片</MinorHeading>
          <CodeSample
            language='json'
            code={`{
  "model": "doubao-seedance-2-5-260628",
  "content": [
    {"type": "text", "text": "女孩抱着狐狸，女孩睁开眼，温柔地看向镜头，镜头缓缓拉远"},
    {"type": "image_url", "image_url": {"url": "https://example.com/input.png"}}
  ],
  "duration": 5,
  "ratio": "16:9",
  "resolution": "480p",
  "output_format": "mov",
  "draft": true
}`}
          />
          <ul className='text-muted-foreground list-disc space-y-1 pl-6 text-sm leading-relaxed'>
            <li>
              <RichText text='Draft 视频仅支持 `480p`；传入其他分辨率会报错。' />
            </li>
            <li>Draft 样片按普通 480p 视频的 token 用量和单价计费。</li>
            <li>
              <RichText text='仅状态为 `succeeded` 的 Draft 任务可用于创建正式视频。' />
            </li>
          </ul>
        </Block>
        <Block>
          <MinorHeading>创建正式视频</MinorHeading>
          <CodeSample
            language='json'
            code={`{
  "model": "doubao-seedance-2-5-260628",
  "content": [
    {
      "type": "draft_task",
      "draft_task": {"id": "<DRAFT_TASK_ID>"}
    }
  ],
  "resolution": "1080p",
  "output_format": "mov"
}`}
          />
          <ul className='text-muted-foreground list-disc space-y-1 pl-6 text-sm leading-relaxed'>
            <li>
              <RichText text='正式视频默认且仅支持 `1080p`。' />
            </li>
            <li>
              <RichText text='`model` 必须与创建 Draft 任务时使用的模型完全一致。' />
            </li>
            <li>
              <RichText text='`draft` 可省略或设为 `false`；Draft 任务 ID 自 `created_at` 起有效 7 天。' />
            </li>
          </ul>
          <Paragraph text='以下“自动复用”参数禁止传入；即使值相同，上游也会报错。' />
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>类型</TableHead>
                <TableHead>参数</TableHead>
                <TableHead>正式视频任务规则</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow>
                <TableCell>固定取值</TableCell>
                <TableCell>
                  <code>model</code>
                </TableCell>
                <TableCell>必须与 Draft 任务一致。</TableCell>
              </TableRow>
              <TableRow>
                <TableCell>固定取值</TableCell>
                <TableCell>
                  <code>draft</code>
                </TableCell>
                <TableCell>
                  {'可省略或设为 '}
                  <code>false</code>。
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell>固定取值</TableCell>
                <TableCell>
                  <code>resolution</code>
                </TableCell>
                <TableCell>
                  {'默认且仅支持 '}
                  <code>1080p</code>。
                </TableCell>
              </TableRow>
              <TableRow>
                <TableCell>自动复用（禁止传入）</TableCell>
                <TableCell>
                  <code>content.text</code>、<code>content.image_url</code>、
                  <code>content.video_url</code>、<code>content.audio_url</code>
                </TableCell>
                <TableCell>自动复用 Draft 任务提示词和参考素材。</TableCell>
              </TableRow>
              <TableRow>
                <TableCell>自动复用（禁止传入）</TableCell>
                <TableCell>
                  <code>duration</code>、<code>ratio</code>、<code>seed</code>、
                  <code>generate_audio</code>、
                  <code>omni_reference_task_type</code>
                </TableCell>
                <TableCell>自动复用 Draft 任务对应值。</TableCell>
              </TableRow>
              <TableRow>
                <TableCell>可重新指定</TableCell>
                <TableCell>
                  <code>return_last_frame</code>、<code>output_format</code>、
                  <code>watermark</code>、<code>service_tier</code>、
                  <code>execution_expires_after</code>、<code>priority</code>、
                  <code>callback_url</code>、<code>safety_identifier</code>
                </TableCell>
                <TableCell>
                  未传入时使用 Seedance 2.5 默认值，不沿用 Draft 任务实际取值。
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
          <Paragraph text='Draft 样片与正式视频独立计费：样片按 480p 计费；正式视频按 1080p 目标分辨率计费。正式视频的 token 单价取决于 Draft 第一步是否包含输入视频；计算输入视频时长时，以第一步输入视频为准，Draft 样片本身不计入。' />
        </Block>
      </ParamDetail>
      <ParamDetail name='duration'>
        <Paragraph text='可选整数，单位为秒。`duration` 与 `frames` 二选一；同时传入时以 `frames` 为准。`-1` 表示由模型选择合适时长。' />
        <KeyValueTable
          kLabel='模型'
          vLabel='取值'
          rows={[
            {
              k: 'Seedance 2.5',
              v: '默认 -1；4 至 30，或 -1',
            },
            {
              k: 'Seedance 2.0 系列',
              v: '4 至 15，或 -1',
            },
          ]}
        />
        <Paragraph text='Seedance 2.5 的视频编辑任务只能使用 `-1`，输入视频时长需要位于 4 至 30 秒。' />
        <Paragraph text='查询接口返回的 `duration` 是按总帧数除以 24 后向下取整的整数秒数，可能与媒体文件的精确时长不同。例如最终生成 133 帧时，实际时长约为 5.54 秒，返回值为 5。' />
      </ParamDetail>
      <ParamDetail name='execution_expires_after'>
        <Paragraph text='可选整数，默认 `172800` 秒（48 小时），范围为 3600 至 259200。时间从任务的 `created_at` 开始计算；超过阈值后，排队或运行中的任务会转为 `expired`。若上游未及时报告，本平台会在阈值后再等待 30 分钟宽限期，仍未结束时自行将任务置为 `expired`（`error.code` 为 `task_expired`）并退还预扣费用。' />
      </ParamDetail>
      <ParamDetail name='frames'>
        <Paragraph text='可选整数，用于精细控制视频时长，支持 Seedance 1.0 Pro 和 1.0 Pro Fast。范围为 29 至 289，并且必须符合 `25 + 4n`（`n` 为正整数）。近似时长按 `frames / 24` 秒计算。' />
      </ParamDetail>
      <ParamDetail name='generate_audio'>
        <Paragraph text='可选布尔值，默认 `true`。`true` 表示生成与画面同步的单声道音频（与传入音频声道数无关），`false` 表示生成无声视频。适用于 Seedance 2.5 和 2.0 系列。' />
      </ParamDetail>
      <ParamDetail name='output_format'>
        <Paragraph text='可选字符串，默认 `mp4`，仅 Seedance 2.5 支持配置。' />
        <KeyValueTable
          kLabel='值'
          vLabel='说明'
          rows={[
            {
              k: 'mp4',
              v: '通用输出格式',
            },
            {
              k: 'mov',
              v: '面向调色、抠像、合成等后期流程；采用 H.264、yuv444p 和 PCM，播放器兼容性可能较低',
            },
          ]}
        />
      </ParamDetail>
      <ParamDetail name='priority'>
        <Paragraph text='可选整数，默认 `0`，范围为 0 至 9。数值越大，同一 Endpoint 内的排队位置越靠前；相同优先级仍按 FIFO 排序。它不会中断已经运行的任务，并且 `service_tier=flex` 时不可用。适用于 Seedance 2.5 和 2.0 系列。' />
      </ParamDetail>
      <ParamDetail name='ratio'>
        <Paragraph text='可选字符串：`16:9`、`4:3`、`1:1`、`3:4`、`9:16`、`21:9`、`adaptive`。' />
        <ul className='text-muted-foreground list-disc space-y-1 pl-6 text-sm leading-relaxed'>
          <li>
            <RichText text='Seedance 2.5 和 2.0 系列默认使用 `adaptive`。' />
          </li>
          <li>
            <RichText text='Seedance 2.5 的首帧、首尾帧和视频编辑任务只支持 `adaptive`。' />
          </li>
          <li>图生视频的指定宽高比与输入图不一致时，平台会居中裁剪输入图。</li>
        </ul>
      </ParamDetail>
      <ParamDetail name='resolution'>
        <Paragraph text='可选字符串：`480p`、`720p`、`1080p`、`4k`。' />
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>模型</TableHead>
              <TableHead>默认值</TableHead>
              <TableHead>可选值</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  Seedance 2.5
                </code>
              </TableCell>
              <TableCell className='text-sm'>720p</TableCell>
              <TableCell className='text-sm'>480p、720p、1080p</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  Seedance 2.0
                </code>
              </TableCell>
              <TableCell className='text-sm'>720p</TableCell>
              <TableCell className='text-sm'>480p、720p、1080p、4k</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  Seedance 2.0 Fast
                </code>
              </TableCell>
              <TableCell className='text-sm'>720p</TableCell>
              <TableCell className='text-sm'>480p、720p</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  Seedance 2.0 Mini
                </code>
              </TableCell>
              <TableCell className='text-sm'>720p</TableCell>
              <TableCell className='text-sm'>480p、720p</TableCell>
            </TableRow>
          </TableBody>
        </Table>
        <MinorHeading>输出像素尺寸</MinorHeading>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>分辨率</TableHead>
              <TableHead>比例</TableHead>
              <TableHead>Seedance 2.5</TableHead>
              <TableHead>Seedance 2.0 系列</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow>
              <TableCell>480p</TableCell>
              <TableCell>16:9</TableCell>
              <TableCell className='text-sm'>854×480</TableCell>
              <TableCell className='text-sm'>864×496</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>480p</TableCell>
              <TableCell>4:3</TableCell>
              <TableCell className='text-sm'>752×560</TableCell>
              <TableCell className='text-sm'>752×560</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>480p</TableCell>
              <TableCell>1:1</TableCell>
              <TableCell className='text-sm'>640×640</TableCell>
              <TableCell className='text-sm'>640×640</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>480p</TableCell>
              <TableCell>3:4</TableCell>
              <TableCell className='text-sm'>560×752</TableCell>
              <TableCell className='text-sm'>560×752</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>480p</TableCell>
              <TableCell>9:16</TableCell>
              <TableCell className='text-sm'>480×854</TableCell>
              <TableCell className='text-sm'>496×864</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>480p</TableCell>
              <TableCell>21:9</TableCell>
              <TableCell className='text-sm'>992×432</TableCell>
              <TableCell className='text-sm'>992×432</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>720p</TableCell>
              <TableCell>16:9</TableCell>
              <TableCell className='text-sm'>1280×720</TableCell>
              <TableCell className='text-sm'>1280×720</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>720p</TableCell>
              <TableCell>4:3</TableCell>
              <TableCell className='text-sm'>1112×834</TableCell>
              <TableCell className='text-sm'>1112×834</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>720p</TableCell>
              <TableCell>1:1</TableCell>
              <TableCell className='text-sm'>960×960</TableCell>
              <TableCell className='text-sm'>960×960</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>720p</TableCell>
              <TableCell>3:4</TableCell>
              <TableCell className='text-sm'>834×1112</TableCell>
              <TableCell className='text-sm'>834×1112</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>720p</TableCell>
              <TableCell>9:16</TableCell>
              <TableCell className='text-sm'>720×1280</TableCell>
              <TableCell className='text-sm'>720×1280</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>720p</TableCell>
              <TableCell>21:9</TableCell>
              <TableCell className='text-sm'>1470×630</TableCell>
              <TableCell className='text-sm'>1470×630</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>1080p</TableCell>
              <TableCell>16:9</TableCell>
              <TableCell className='text-sm'>1920×1080</TableCell>
              <TableCell className='text-sm'>1920×1080</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>1080p</TableCell>
              <TableCell>4:3</TableCell>
              <TableCell className='text-sm'>1664×1248</TableCell>
              <TableCell className='text-sm'>1664×1248</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>1080p</TableCell>
              <TableCell>1:1</TableCell>
              <TableCell className='text-sm'>1440×1440</TableCell>
              <TableCell className='text-sm'>1440×1440</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>1080p</TableCell>
              <TableCell>3:4</TableCell>
              <TableCell className='text-sm'>1248×1664</TableCell>
              <TableCell className='text-sm'>1248×1664</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>1080p</TableCell>
              <TableCell>9:16</TableCell>
              <TableCell className='text-sm'>1080×1920</TableCell>
              <TableCell className='text-sm'>1080×1920</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>1080p</TableCell>
              <TableCell>21:9</TableCell>
              <TableCell className='text-sm'>2206×946</TableCell>
              <TableCell className='text-sm'>2206×946</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>4K</TableCell>
              <TableCell>16:9</TableCell>
              <TableCell className='text-sm'>-</TableCell>
              <TableCell className='text-sm'>3840×2160</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>4K</TableCell>
              <TableCell>4:3</TableCell>
              <TableCell className='text-sm'>-</TableCell>
              <TableCell className='text-sm'>3326×2494</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>4K</TableCell>
              <TableCell>1:1</TableCell>
              <TableCell className='text-sm'>-</TableCell>
              <TableCell className='text-sm'>2880×2880</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>4K</TableCell>
              <TableCell>3:4</TableCell>
              <TableCell className='text-sm'>-</TableCell>
              <TableCell className='text-sm'>2494×3326</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>4K</TableCell>
              <TableCell>9:16</TableCell>
              <TableCell className='text-sm'>-</TableCell>
              <TableCell className='text-sm'>2160×3840</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>4K</TableCell>
              <TableCell>21:9</TableCell>
              <TableCell className='text-sm'>-</TableCell>
              <TableCell className='text-sm'>4398×1886</TableCell>
            </TableRow>
          </TableBody>
        </Table>
        <Paragraph text='1080p 不支持 Seedance 2.0 Fast/Mini；4K 仅 Seedance 2.0 支持。' />
        <Paragraph text='Seedance 2.5 的 1080p 输出和 Seedance 2.0 的 4K 输出使用 10-bit H.265/HEVC 编码，播放环境需要支持相应格式；如遇兼容问题，可使用 VLC、MPV、QuickTime Player 等播放器。' />
        <Paragraph text='Draft 模式的分辨率规则优先于本表：创建 Draft 样片时仅允许 `480p`；基于 Draft 样片创建正式视频时仅允许 `1080p`。' />
      </ParamDetail>
      <ParamDetail name='return_last_frame'>
        <Paragraph text='可选布尔值，默认 `false`。设为 `true` 后，可以通过查询任务接口获取生成视频的 PNG 尾帧；图片尺寸与视频一致且不带水印。可将该尾帧作为下一个任务的首帧，以生成连续片段。' />
      </ParamDetail>
      <ParamDetail name='safety_identifier'>
        <Paragraph text='可选字符串，用于标识终端用户。对同一用户应保持固定且唯一，最长 64 个英文字符。建议使用用户名、用户 ID 或邮箱的哈希值，避免提交原始个人信息。' />
      </ParamDetail>
      <ParamDetail name='seed'>
        <Paragraph text='可选整数，默认 `-1`，`-1` 表示随机选择。仅 Seedance 1.5 Pro、1.0 Pro 和 1.0 Pro Fast 支持；Seedance 2.5 和 2.0 系列不支持该参数。' />
      </ParamDetail>
      <ParamDetail name='service_tier'>
        <Paragraph text='默认 `default`。仅 Seedance 1.5 Pro、1.0 Pro 和 1.0 Pro Fast 支持。' />
        <KeyValueTable
          kLabel='值'
          vLabel='说明'
          rows={[
            {
              k: 'default',
              v: '在线推理，适合更关注响应时效的任务',
            },
            {
              k: 'flex',
              v: '离线推理，适合对时延不敏感的任务；价格按在线推理的 50% 计费',
            },
          ]}
        />
        <Paragraph text='任务提交后不能修改服务等级。Seedance 2.5 和 2.0 系列暂不支持 `flex`。' />
      </ParamDetail>
      <ParamDetail name='tools'>
        <Paragraph text='可选对象数组，适用于 Seedance 2.5 和 2.0 系列。' />
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>字段</TableHead>
              <TableHead>类型</TableHead>
              <TableHead>是否必填</TableHead>
              <TableHead>可选值</TableHead>
              <TableHead>说明</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  tools[].type
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                是
              </TableCell>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  web_search
                </code>
              </TableCell>
              <TableCell className='text-sm leading-relaxed'>
                允许模型按提示词判断是否联网搜索
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
        <Paragraph>
          联网搜索可能改善时效性内容的生成效果，但会增加处理时间。查询任务接口的
          <code className='bg-muted mx-1 rounded px-1 py-0.5 text-xs'>
            usage.tool_usage.web_search
          </code>
          可用于查看实际搜索次数。
        </Paragraph>
      </ParamDetail>
      <ParamDetail name='watermark'>
        <Paragraph text='可选布尔值，默认 `false`。设为 `true` 时，视频右下角带有 AI 生成水印。' />
      </ParamDetail>
      <Block>
        <SubHeading text='成功响应' />
        <CodeSample language='json' code={createTaskSuccessJson} />
        <ResponseFieldTable
          rows={[
            {
              name: 'id',
              type: 'string',
              description:
                '视频任务 ID（`task_` 开头，由本平台生成，不同于上游任务 ID）。后续查询必须使用此 ID。',
            },
          ]}
        />
      </Block>
      <Block>
        <SubHeading text='失败响应' />
        <Paragraph text='请求 JSON 无效、参数校验失败或上游拒绝请求时返回 HTTP `400`，错误信封如下（上游拒绝时 `message` 为上游原始错误响应文本）；API Key 无效时返回 HTTP `401`，使用 `{"error": {"message", "type", "code"}}` 信封。' />
        <CodeSample language='json' code={createTaskErrorJson} />
        <ResponseFieldTable
          rows={[
            {
              name: 'error.code',
              type: 'string',
              description:
                '错误码；本平台校验失败或上游拒绝时为 `invalid_request`。',
            },
            {
              name: 'error.message',
              type: 'string',
              description: '错误说明。',
            },
          ]}
        />
        <Note text='失败响应不包含任务 ID，不能用于后续查询。' />
      </Block>
    </EndpointCard>
  )
}

const runningTaskJson = `{
  "id": "task_BGJr5gwbfZssKxUFWXqLgOt2Dek4lmO5",
  "model": "doubao-seedance-2-5-260628",
  "status": "running",
  "content": {
    "video_url": "",
    "last_frame_url": "",
    "file_url": ""
  },
  "usage": {
    "prompt_tokens": 0,
    "completion_tokens": 0,
    "total_tokens": 0,
    "prompt_tokens_details": {
      "cached_tokens": 0
    },
    "completion_tokens_details": {
      "reasoning_tokens": 0
    }
  },
  "frames": null,
  "framespersecond": null,
  "output_format": "mov",
  "created_at": 1786015617,
  "updated_at": 1786015617
}`

const succeededTaskJson = `{
	"id": "task_GlvqIS2axgKamyWpKzSP4yuV0hdf0q78",
	"model": "doubao-seedance-2-0-mini-260615",
	"status": "succeeded",
	"content": {
		"video_url": "https://..."
	},
	"usage": {
		"completion_tokens": 120946,
		"total_tokens": 120946
	},
	"created_at": 1787624623,
	"updated_at": 1787624776,
	"seed": 0,
	"resolution": "480p",
	"ratio": "9:16",
	"duration": 12,
	"framespersecond": 24,
	"service_tier": "default",
	"execution_expires_after": 172800,
	"generate_audio": true,
	"draft": false,
	"priority": 0,
	"output_format": "mp4"
}`

const failedTaskJson = `{
  "created_at": 1792410250,
  "error": {
    "code": "InputTextSensitiveContentDetected",
    "message": "The request failed because the input text may contain sensitive information."
  },
  "id": "task_CBkXCGoJ5l8VezyE5oHwlQwFEVJf5BxV",
  "model": "doubao-seedance-2-0-mini-260615",
  "status": "failed",
  "updated_at": 1792410257
}`

const taskStatuses: StatusRow[] = [
  {
    name: 'queued',
    description: '任务排队中。',
  },
  {
    name: 'running',
    description: '任务运行中。',
  },
  {
    name: 'cancelled',
    description: '任务已在上游被取消（不推送回调）。',
  },
  {
    name: 'succeeded',
    description: '任务成功，可读取 `content.video_url`。',
  },
  {
    name: 'failed',
    description: '视频生成失败，原因见 `error`。',
  },
  {
    name: 'expired',
    description:
      '任务超时（含本平台判定的超时，此时 `error.code` 为 `task_expired`）。',
  },
]

const taskResponseFields: ResponseFieldRow[] = [
  {
    name: 'id',
    type: 'string',
    description:
      '视频任务 ID，与创建任务时返回的 `id` 相同；不会返回上游任务 ID。',
  },
  {
    name: 'model',
    type: 'string',
    description: '查询结果中的模型名。',
  },
  {
    name: 'status',
    type: 'string',
    description: '本次查询得到的任务状态。',
  },
  {
    name: 'content.video_url',
    type: 'string',
    description:
      '任务完成时返回视频 URL（上游签名地址，原样透传，会过期，请及时下载）；未完成时为空或不返回。',
  },
  {
    name: 'seed',
    type: 'integer',
    description: '查询结果中的随机种子。',
  },
  {
    name: 'resolution',
    type: 'string',
    description: '查询结果中的分辨率。',
  },
  {
    name: 'duration',
    type: 'integer',
    description: '查询结果中的视频时长。',
  },
  {
    name: 'ratio',
    type: 'string',
    description: '查询结果中的视频比例。',
  },
  {
    name: 'framespersecond',
    type: 'integer',
    description: '查询结果中的帧率。',
  },
  {
    name: 'output_format',
    type: 'string',
    description: '输出视频封装格式；Seedance 2.5 返回 `mp4` 或 `mov`。',
  },
  {
    name: 'service_tier',
    type: 'string',
    description: '查询结果中的服务等级。',
  },
  {
    name: 'tools',
    type: 'array',
    description: '查询结果中的工具配置。',
  },
  {
    name: 'tools[].type',
    type: 'string',
    description: '工具类型。',
  },
  {
    name: 'usage.completion_tokens',
    type: 'integer',
    description: '完成消耗。',
  },
  {
    name: 'usage.total_tokens',
    type: 'integer',
    description: '总消耗。',
  },
  {
    name: 'usage.tool_usage.web_search',
    type: 'integer',
    description: 'Web 搜索工具消耗。',
  },
  {
    name: 'error',
    type: 'object',
    description: '任务失败或超时时返回，包含 `code` 与 `message`。',
  },
  {
    name: 'created_at',
    type: 'integer',
    description: '任务创建时间（Unix 秒）。',
  },
  {
    name: 'updated_at',
    type: 'integer',
    description: '任务最近更新时间（Unix 秒）。',
  },
]

function GetVideoTask() {
  const baseUrl = useApiBaseUrl()
  return (
    <EndpointCard
      anchor='get-video-task-endpoint'
      number='2'
      title='获取视频任务状态'
      method='GET'
      path='/api/v3/contents/generations/tasks/{task_id}'
      remark='返回本平台后台从上游同步到的最新任务状态（提交后、首次同步前返回 `queued`），状态变化可能有数秒延迟。仅能查询当前认证用户创建的任务。'
    >
      <Block>
        <SubHeading text='请求示例' />
        <CodeSample
          language='bash'
          code={`curl --request GET \\
  "${baseUrl}/api/v3/contents/generations/tasks/task_BGJr5gwbfZssKxUFWXqLgOt2Dek4lmO5" \\
  --header "Authorization: Bearer $KT_API_KEY"`}
        />
      </Block>
      <Block>
        <SubHeading text='路径参数' />
        <FieldTable
          rows={[
            {
              name: 'task_id',
              type: 'string',
              required: '是',
              description: '创建任务时响应中的 `id`。',
            },
          ]}
        />
      </Block>
      <Block>
        <SubHeading text='处理中响应' />
        <CodeSample language='json' code={runningTaskJson} />
      </Block>
      <Block>
        <SubHeading text='完成响应' />
        <CodeSample language='json' code={succeededTaskJson} />
      </Block>
      <Block>
        <SubHeading text='任务状态' />
        <StatusTable rows={taskStatuses} />
        <ResponseFieldTable rows={taskResponseFields} />
      </Block>
      <Block>
        <SubHeading text='任务失败响应' />
        <Paragraph text='任务失败时，查询接口仍返回 HTTP `200`，响应体中 `status` 为 `failed`，并带上 `error.code` 与 `error.message`，示例如下。请以 `status` 判断任务结果，不要依赖 HTTP 状态码。任务不存在或不属于当前用户时返回 HTTP `404`（`{"error": {"code": "task_not_found", "message": "Task not found"}}`）。' />
        <CodeSample language='json' code={failedTaskJson} />
        <Paragraph text='建议间隔数秒轮询，直到 `status` 为 `succeeded`、`failed`、`cancelled` 或 `expired`。' />
      </Block>
    </EndpointCard>
  )
}

const createAssetGroupJson = `{
  "ResponseMetadata": {
    "RequestId": "20260328000000000000000000000000",
    "Action": "CreateAssetGroup",
    "Version": "2024-01-01",
    "Service": "ark",
    "Region": "cn-beijing"
  },
  "Result": {
    "Id": "group-2026**********-*****"
  }
}`

const createAssetGroupUpstreamErrorJson = `{
  "ResponseMetadata": {
    "RequestId": "20260328000000000000000000000001",
    "Action": "CreateAssetGroup",
    "Version": "2024-01-01",
    "Service": "ark",
    "Region": "cn-beijing",
    "Error": {
      "Code": "InvalidParameter",
      "Message": "asset group name is required"
    }
  }
}`

const assetLocalErrorJson = `{
  "ResponseMetadata": {
    "Action": "CreateAsset",
    "Error": {
      "Code": "InvalidParameter",
      "Message": "record not found"
    },
    "Region": "cn-beijing",
    "RequestId": "202610191145235772347488268d9d6m3WtwfaB",
    "Service": "ark",
    "Version": "2024-01-01"
  }
}`

function CreateAssetGroup() {
  const baseUrl = useApiBaseUrl()
  return (
    <EndpointCard
      anchor='create-asset-group-endpoint'
      number='3'
      title='创建素材组'
      method='POST'
      path='/api/assets/?Action=CreateAssetGroup&Version=2024-01-01'
      contentType='application/json'
      remark='素材上传前必须先创建素材组。请求和成功响应遵循火山方舟 CreateAssetGroup 协议。素材组会绑定到本次请求选中的渠道，后续该组的素材操作固定走该渠道。'
    >
      <Block>
        <SubHeading text='请求示例' />
        <CodeSample
          language='bash'
          code={`curl --request POST \\
  "${baseUrl}/api/assets/?Action=CreateAssetGroup&Version=2024-01-01" \\
  --header "Authorization: Bearer $KT_API_KEY" \\
  --header "Content-Type: application/json" \\
  --data '{
    "Name": "虚拟人像素材组",
    "Description": "角色参考素材",
    "GroupType": "AIGC",
    "ProjectName": "default"
  }'`}
        />
      </Block>
      <Block>
        <SubHeading text='请求字段' />
        <FieldTable
          rows={[
            {
              name: 'Name',
              type: 'string',
              required: '是',
              description: '去除首尾空白后不能为空，最长 64 个 Unicode 字符。',
            },
            {
              name: 'Description',
              type: 'string',
              required: '否',
              description: '素材组描述，最长 300 个 Unicode 字符。',
            },
            {
              name: 'GroupType',
              type: 'string',
              required: '否',
              description: '素材组类型，当前仅支持 `AIGC`，默认 `AIGC`。',
            },
            {
              name: 'ProjectName',
              type: 'string',
              required: '否',
              description: '资源所属项目，默认 `default`。',
            },
          ]}
        />
        <Paragraph text='为兼容已有调用方，接口仍接受 `name`、`description`、`group_type` 和 `project_name`。' />
      </Block>
      <Block>
        <SubHeading text='成功响应' />
        <CodeSample language='json' code={createAssetGroupJson} />
        <ResponseFieldTable
          rows={[
            {
              name: 'ResponseMetadata.RequestId',
              type: 'string',
              description: '上游请求 ID；渠道未提供时为空字符串。',
            },
            {
              name: 'ResponseMetadata.Action',
              type: 'string',
              description: '固定为 `CreateAssetGroup`。',
            },
            {
              name: 'ResponseMetadata.Version',
              type: 'string',
              description: '固定为 `2024-01-01`。',
            },
            {
              name: 'ResponseMetadata.Service',
              type: 'string',
              description: '固定为 `ark`。',
            },
            {
              name: 'ResponseMetadata.Region',
              type: 'string',
              description: '固定为 `cn-beijing`。',
            },
            {
              name: 'Result.Id',
              type: 'string',
              description: '素材组 ID，上传和查询素材时使用。',
            },
          ]}
        />
      </Block>
      <Block>
        <SubHeading text='失败响应' />
        <Paragraph text='火山上游错误的 HTTP 状态码和响应体会被原样返回：' />
        <CodeSample language='json' code={createAssetGroupUpstreamErrorJson} />
        <Paragraph text='本平台校验失败（如素材组不存在或不属于当前用户）、余额不足或上游不可用时，由本平台生成同样信封的错误，`Error.Code` 为 `InvalidParameter`、`InsufficientQuota` 或 `UpstreamError` 等，`RequestId` 为本平台请求 ID；API Key 无效时返回 HTTP `401`，使用 `{"error": {"message", "type", "code"}}` 信封。' />
        <CodeSample language='json' code={assetLocalErrorJson} />
      </Block>
    </EndpointCard>
  )
}

const createAssetJson = `{
  "ResponseMetadata": {
    "RequestId": "20260328000000000000000000000000",
    "Action": "CreateAsset",
    "Version": "2024-01-01",
    "Service": "ark",
    "Region": "cn-beijing"
  },
  "Result": {
    "Id": "asset-2026**********-*****"
  }
}`

const createAssetErrorJson = `{
  "ResponseMetadata": {
    "RequestId": "20260328000000000000000000000000",
    "Action": "CreateAsset",
    "Version": "2024-01-01",
    "Service": "ark",
    "Region": "cn-beijing",
    "Error": {
      "Code": "InvalidParameter.URL",
      "Message": "URL must be an absolute HTTP(S) URL"
    }
  }
}`

function CreateAsset() {
  const baseUrl = useApiBaseUrl()
  return (
    <EndpointCard
      anchor='create-asset-endpoint'
      number='4'
      title='创建素材资产'
      method='POST'
      path='/api/assets/?Action=CreateAsset&Version=2024-01-01'
      contentType='application/json'
      remark='通过服务端可访问的素材 URL 创建素材。`GroupId` 必须是当前认证用户创建的素材组，`ProjectName` 必须与素材组所在项目一致。素材会绑定到创建素材组时选定的渠道，不能跨渠道用于视频任务。'
    >
      <Block>
        <SubHeading text='请求示例' />
        <CodeSample
          language='bash'
          code={`curl --request POST \\
  "${baseUrl}/api/assets/?Action=CreateAsset&Version=2024-01-01" \\
  --header "Authorization: Bearer $KT_API_KEY" \\
  --header "Content-Type: application/json" \\
  --data '{
    "GroupId": "group-20260813103927-cgxm6",
    "URL": "https://media.example.com/reference.png",
    "AssetType": "Image",
    "Name": "女主参考图",
    "ProjectName": "default"
  }'`}
        />
      </Block>
      <Block>
        <SubHeading text='请求字段' />
        <FieldTable
          rows={[
            {
              name: 'GroupId',
              type: 'string',
              required: '是',
              description: '当前用户拥有的素材组 ID；不能为空。',
            },
            {
              name: 'URL',
              type: 'string',
              required: '是',
              description: '服务端可直接访问的绝对 HTTP(S) 素材 URL。',
            },
            {
              name: 'AssetType',
              type: 'string',
              required: '是',
              description: '`Image`、`Video` 或 `Audio`，大小写不敏感。',
            },
            {
              name: 'Name',
              type: 'string',
              required: '是',
              description:
                '素材名称，必填且不能为空，最长 64 个 Unicode 字符。',
            },
            {
              name: 'ProjectName',
              type: 'string',
              required: '否',
              description:
                '所属项目，默认 `default`，必须与素材组所在项目一致。',
            },
          ]}
        />
        <Paragraph text='服务端校验 `URL` 为绝对 HTTP(S) 地址，并校验 `AssetType` 是否受支持。' />
      </Block>
      <Block>
        <SubHeading text='Multipart 表单兼容' />
        <Paragraph text='接口同时兼容 multipart 表单。使用 `curl --form` 时不要手动设置 `Content-Type`，curl 会自动生成包含 boundary 的请求头：' />
        <CodeSample
          language='bash'
          code={`curl --request POST \\
  "${baseUrl}/api/assets/?Action=CreateAsset&Version=2024-01-01" \\
  --header "Authorization: Bearer $KT_API_KEY" \\
  --form "GroupId=group-20260813103927-cgxm6" \\
  --form "URL=https://media.example.com/reference.png" \\
  --form "AssetType=Image" \\
  --form "Name=女主参考图" \\
  --form "ProjectName=default"`}
        />
      </Block>
      <Block>
        <SubHeading text='成功响应' />
        <CodeSample language='json' code={createAssetJson} />
        <ResponseFieldTable
          rows={[
            {
              name: 'ResponseMetadata.RequestId',
              type: 'string',
              description: '上游请求 ID；渠道未提供时为空字符串。',
            },
            {
              name: 'ResponseMetadata.Action',
              type: 'string',
              description: '固定为 `CreateAsset`。',
            },
            {
              name: 'ResponseMetadata.Version',
              type: 'string',
              description: '固定为 `2024-01-01`。',
            },
            {
              name: 'ResponseMetadata.Service',
              type: 'string',
              description: '固定为 `ark`。',
            },
            {
              name: 'ResponseMetadata.Region',
              type: 'string',
              description: '固定为 `cn-beijing`。',
            },
            {
              name: 'Result.Id',
              type: 'string',
              description: '新建素材 ID，用于后续查询处理状态。',
            },
          ]}
        />
      </Block>
      <Block>
        <SubHeading text='失败响应' />
        <CodeSample language='json' code={createAssetErrorJson} />
        <Paragraph text='所有素材 Action 的错误均使用火山方舟 `ResponseMetadata.Error` 信封。火山渠道错误会保留上游 HTTP 状态码、`Content-Type` 和完整响应体；`Error.Data` 为可选字段。常见失败包括素材组不存在或不属于当前用户、名称超过长度限制、类型不支持、URL 无效，以及素材服务调用失败。' />
      </Block>
    </EndpointCard>
  )
}

const listAssetGroupsJson = `{
  "ResponseMetadata": {
    "RequestId": "20260328000000000000000000000000",
    "Action": "ListAssetGroups",
    "Version": "2024-01-01",
    "Service": "ark",
    "Region": "cn-beijing"
  },
  "Result": {
    "TotalCount": 1,
    "Items": [
      {
        "Id": "group-20260813103927-cgxm6",
        "Name": "人物素材组",
        "Title": "人物素材组",
        "Description": "人物参考素材",
        "GroupType": "AIGC",
        "ProjectName": "default",
        "CreateTime": "2026-08-13T10:39:27Z",
        "UpdateTime": "2026-08-17T10:08:49Z"
      }
    ],
    "PageNumber": 1,
    "PageSize": 10
  }
}`

function ListAssetGroups() {
  const baseUrl = useApiBaseUrl()
  return (
    <EndpointCard
      anchor='list-asset-groups-endpoint'
      number='5'
      title='查询素材资产组合列表'
      method='POST'
      path='/api/assets/?Action=ListAssetGroups&Version=2024-01-01'
      contentType='application/json'
      remark='根据素材组 ID、类型和名称等条件分页查询当前用户的素材组。请求和响应遵循火山方舟 ListAssetGroups 协议。'
    >
      <Block>
        <SubHeading text='请求示例' />
        <CodeSample
          language='bash'
          code={`curl --request POST \\
  "${baseUrl}/api/assets/?Action=ListAssetGroups&Version=2024-01-01" \\
  --header "Authorization: Bearer $KT_API_KEY" \\
  --header "Content-Type: application/json" \\
  --data '{
    "Filter": {
      "GroupIds": ["group-20260813103927-cgxm6"],
      "GroupType": "AIGC",
      "Name": "人物"
    },
    "PageNumber": 1,
    "PageSize": 10,
    "SortBy": "CreateTime",
    "SortOrder": "Desc",
    "ProjectName": "default"
  }'`}
        />
      </Block>
      <Block>
        <SubHeading text='请求字段' />
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className='w-[28%]'>字段</TableHead>
              <TableHead className='w-[16%]'>类型</TableHead>
              <TableHead className='w-[10%]'>必填</TableHead>
              <TableHead className='w-[10%]'>默认值</TableHead>
              <TableHead>说明</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  Filter
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                object
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                是
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                无
              </TableCell>
              <TableCell className='text-sm'>搜索过滤条件。</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  Filter.GroupIds
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string[]
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                否
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                无
              </TableCell>
              <TableCell className='text-sm'>
                素材组 ID 列表；不属于当前用户的 ID
                会被忽略，省略时查询当前用户该类型的全部素材组。
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  Filter.GroupType
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                是
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                无
              </TableCell>
              <TableCell className='text-sm'>
                <RichText text='素材组类型：`AIGC` 或 `LivenessFace`。' />
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  Filter.Name
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                否
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                无
              </TableCell>
              <TableCell className='text-sm'>
                素材组名称，最多 64 个字符。
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  PageNumber
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                integer
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                是
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                无
              </TableCell>
              <TableCell className='text-sm'>
                页码，必须大于或等于 1。
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  PageSize
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                integer
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                是
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                无
              </TableCell>
              <TableCell className='text-sm'>
                每页数量，范围为 1 至 100。
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  SortBy
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                否
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                CreateTime
              </TableCell>
              <TableCell className='text-sm'>
                <RichText text='排序字段：`CreateTime` 或 `UpdateTime`。' />
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  SortOrder
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                否
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                Desc
              </TableCell>
              <TableCell className='text-sm'>
                <RichText text='排序方向：`Desc` 或 `Asc`。' />
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  ProjectName
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                否
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                default
              </TableCell>
              <TableCell className='text-sm'>资源所属项目。</TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </Block>
      <Block>
        <SubHeading text='成功响应' />
        <CodeSample language='json' code={listAssetGroupsJson} />
        <ResponseFieldTable
          rows={[
            {
              name: 'ResponseMetadata',
              type: 'object',
              description: '请求元数据，字段含义与其他素材 Action 接口一致。',
            },
            {
              name: 'Result.TotalCount',
              type: 'integer',
              description: '符合条件的素材组总数。',
            },
            {
              name: 'Result.Items',
              type: 'array',
              description: '符合条件的素材组列表。',
            },
            {
              name: 'Result.Items[].Id',
              type: 'string',
              description: '素材组 ID。',
            },
            {
              name: 'Result.Items[].Name',
              type: 'string',
              description: '素材组名称。',
            },
            {
              name: 'Result.Items[].Title',
              type: 'string',
              description: '上游可能返回的兼容标题字段。',
            },
            {
              name: 'Result.Items[].Description',
              type: 'string',
              description: '素材组描述。',
            },
            {
              name: 'Result.Items[].GroupType',
              type: 'string',
              description: '素材组类型。',
            },
            {
              name: 'Result.Items[].ProjectName',
              type: 'string',
              description: '素材组所属项目。',
            },
            {
              name: 'Result.Items[].CreateTime',
              type: 'string',
              description: '创建时间。',
            },
            {
              name: 'Result.Items[].UpdateTime',
              type: 'string',
              description: '更新时间。',
            },
            {
              name: 'Result.PageNumber',
              type: 'integer',
              description: '当前页码。',
            },
            {
              name: 'Result.PageSize',
              type: 'integer',
              description: '当前每页数量。',
            },
          ]}
        />
        <TimestampNote />
      </Block>
      <Block>
        <SubHeading text='失败响应' />
        <Paragraph text='`Filter.GroupType` 缺失（`InvalidParameter`：`Filter.GroupType is required`）、分页参数无效或上游查询失败时请求失败。本接口只接受 JSON 请求体。' />
      </Block>
    </EndpointCard>
  )
}

const getAssetGroupJson = `{
  "ResponseMetadata": {
    "RequestId": "20260328000000000000000000000000",
    "Action": "GetAssetGroup",
    "Version": "2024-01-01",
    "Service": "ark",
    "Region": "cn-beijing"
  },
  "Result": {
    "Id": "group-20260813103927-cgxm6",
    "Name": "人物素材组",
    "Description": "人物参考素材",
    "GroupType": "AIGC",
    "ProjectName": "default",
    "CreateTime": "2026-08-13T10:39:27Z",
    "UpdateTime": "2026-08-24T10:08:49Z"
  }
}`

function GetAssetGroup() {
  const baseUrl = useApiBaseUrl()
  return (
    <EndpointCard
      anchor='get-asset-group-endpoint'
      number='6'
      title='查询单个素材资产组合信息'
      method='POST'
      path='/api/assets/?Action=GetAssetGroup&Version=2024-01-01'
      contentType='application/json'
      remark='根据素材组 ID 获取详情。请求和响应遵循火山方舟 GetAssetGroup 协议。`Id` 必须是当前认证用户创建的素材组 ID；本平台校验归属后，将请求转发到创建该素材组时使用的渠道，响应原样返回。'
    >
      <Block>
        <SubHeading text='请求示例' />
        <CodeSample
          language='bash'
          code={`curl --request POST \\
  "${baseUrl}/api/assets/?Action=GetAssetGroup&Version=2024-01-01" \\
  --header "Authorization: Bearer $KT_API_KEY" \\
  --header "Content-Type: application/json" \\
  --data '{
    "Id": "group-20260813103927-cgxm6",
    "ProjectName": "default"
  }'`}
        />
      </Block>
      <Block>
        <SubHeading text='请求字段' />
        <FieldTable
          rows={[
            {
              name: 'Id',
              type: 'string',
              required: '是',
              description: '当前用户拥有的素材组 ID。',
            },
            {
              name: 'ProjectName',
              type: 'string',
              required: '否',
              description: '素材组所属项目，默认 `default`。',
            },
          ]}
        />
      </Block>
      <Block>
        <SubHeading text='成功响应' />
        <CodeSample language='json' code={getAssetGroupJson} />
        <ResponseFieldTable
          rows={[
            {
              name: 'ResponseMetadata',
              type: 'object',
              description: '请求元数据，字段含义与其他素材 Action 接口一致。',
            },
            {
              name: 'Result.Id',
              type: 'string',
              description: '素材组 ID。',
            },
            {
              name: 'Result.Name',
              type: 'string',
              description: '素材组名称。',
            },
            {
              name: 'Result.Description',
              type: 'string',
              description: '素材组描述。',
            },
            {
              name: 'Result.GroupType',
              type: 'string',
              description: '`AIGC` 或 `LivenessFace`。',
            },
            {
              name: 'Result.ProjectName',
              type: 'string',
              description: '素材组所属项目。',
            },
            {
              name: 'Result.CreateTime',
              type: 'string',
              description: '创建时间。',
            },
            {
              name: 'Result.UpdateTime',
              type: 'string',
              description: '更新时间。',
            },
          ]}
        />
        <TimestampNote />
      </Block>
      <Block>
        <SubHeading text='失败响应' />
        <Paragraph text='素材组不存在、不属于当前用户或上游查询失败时请求失败。' />
      </Block>
    </EndpointCard>
  )
}

const listAssetsJson = `{
  "ResponseMetadata": {
    "RequestId": "20260328000000000000000000000000",
    "Action": "ListAssets",
    "Version": "2024-01-01",
    "Service": "ark",
    "Region": "cn-beijing"
  },
  "Result": {
    "Items": [
      {
        "Id": "asset-20260813150950-t2tpn",
        "Name": "女主参考图",
        "URL": "https://example.com/temporary-asset-url",
        "GroupId": "group-20260813103927-cgxm6",
        "AssetType": "Image",
        "Status": "Active",
        "Moderation": {"Strategy": "Default"},
        "CreateTime": "2026-08-13T15:09:50Z",
        "UpdateTime": "2026-08-13T15:09:54Z",
        "LastInferenceTime": "2026-08-17T10:08:49+08:00",
        "ProjectName": "default"
      }
    ],
    "TotalCount": 1,
    "PageNumber": 1,
    "PageSize": 20
  }
}`

const listAssetsErrorJson = `{
  "ResponseMetadata": {
    "Action": "ListAssets",
    "Error": {
      "Code": "InvalidParameter",
      "Message": "Filter.GroupIds must contain at least one asset group id"
    },
    "Region": "cn-beijing",
    "RequestId": "202610191145236113335288268d9d6vf31oFj2",
    "Service": "ark",
    "Version": "2024-01-01"
  }
}`

function ListAssets() {
  const baseUrl = useApiBaseUrl()
  return (
    <EndpointCard
      anchor='list-assets-endpoint'
      number='7'
      title='查询素材资产列表'
      method='POST'
      path='/api/assets/?Action=ListAssets&Version=2024-01-01'
      contentType='application/json'
      remark='按素材组、素材组类型、处理状态和名称等条件分页查询素材及其最新处理状态（由上游实时返回）。请求和响应遵循火山方舟 ListAssets 协议。'
    >
      <Block>
        <SubHeading text='请求示例' />
        <CodeSample
          language='bash'
          code={`curl --request POST \\
  "${baseUrl}/api/assets/?Action=ListAssets&Version=2024-01-01" \\
  --header "Authorization: Bearer $KT_API_KEY" \\
  --header "Content-Type: application/json" \\
  --data '{
    "Filter": {
      "GroupIds": ["group-20260813103927-cgxm6"],
      "GroupType": "AIGC",
      "Statuses": ["Active", "Processing", "Failed"],
      "Name": "女主"
    },
    "PageNumber": 1,
    "PageSize": 20,
    "SortBy": "CreateTime",
    "SortOrder": "Desc",
    "ProjectName": "default"
  }'`}
        />
      </Block>
      <Block>
        <SubHeading text='请求字段' />
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className='w-[28%]'>字段</TableHead>
              <TableHead className='w-[16%]'>类型</TableHead>
              <TableHead className='w-[10%]'>必填</TableHead>
              <TableHead className='w-[10%]'>默认值</TableHead>
              <TableHead>说明</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  Filter
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                object
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                是
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                无
              </TableCell>
              <TableCell className='text-sm'>搜索过滤条件。</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  Filter.GroupIds
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string[]
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                是
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                无
              </TableCell>
              <TableCell className='text-sm'>
                当前用户创建的素材组 ID 列表，至少包含一个
                ID；请求会路由到第一个素材组所在的渠道。
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  Filter.GroupType
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                否
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                无
              </TableCell>
              <TableCell className='text-sm'>
                <RichText text='素材组类型：`AIGC` 或 `LivenessFace`。' />
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  Filter.Statuses
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string[]
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                否
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                无
              </TableCell>
              <TableCell className='text-sm'>
                <RichText text='素材状态列表：`Active`、`Processing`、`Failed`。' />
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  Filter.Name
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                否
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                无
              </TableCell>
              <TableCell className='text-sm'>
                素材名称，最多 64 个字符，支持模糊搜索。
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  PageNumber
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                integer
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                否
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>1</TableCell>
              <TableCell className='text-sm'>
                页码，必须大于或等于 1。
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  PageSize
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                integer
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                否
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                20
              </TableCell>
              <TableCell className='text-sm'>
                每页数量，范围为 1 至 100。
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  SortBy
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                否
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                CreateTime
              </TableCell>
              <TableCell className='text-sm'>
                <RichText text='排序字段：`CreateTime`、`UpdateTime` 或 `GroupId`。' />
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  SortOrder
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                否
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                Desc
              </TableCell>
              <TableCell className='text-sm'>
                <RichText text='排序方向：`Desc` 或 `Asc`。' />
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  ProjectName
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                否
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                default
              </TableCell>
              <TableCell className='text-sm'>资源所属项目。</TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </Block>
      <Block>
        <SubHeading text='成功响应' />
        <CodeSample language='json' code={listAssetsJson} />
        <ResponseFieldTable
          rows={[
            {
              name: 'ResponseMetadata',
              type: 'object',
              description: '请求元数据，字段含义与其他素材 Action 接口一致。',
            },
            {
              name: 'Result.Items',
              type: 'array',
              description: '符合条件的素材列表，字段定义与素材详情一致。',
            },
            {
              name: 'Result.TotalCount',
              type: 'integer',
              description: '符合条件的素材总数。',
            },
            {
              name: 'Result.PageNumber',
              type: 'integer',
              description: '当前页码。',
            },
            {
              name: 'Result.PageSize',
              type: 'integer',
              description: '当前每页数量。',
            },
          ]}
        />
        <AssetDetailNote />
      </Block>
      <Block>
        <SubHeading text='失败响应' />
        <Paragraph text='`Filter.GroupIds` 缺失或为空（`Filter.GroupIds must contain at least one asset group id`）、素材组不存在或不属于当前用户（`record not found`），以及页码不合法时请求失败。请求会转发到第一个素材组所在的渠道，请勿在一次查询中混用不同渠道的素材组。例如：' />
        <CodeSample language='json' code={listAssetsErrorJson} />
      </Block>
    </EndpointCard>
  )
}

const getAssetJson = `{
  "ResponseMetadata": {
    "RequestId": "20260328000000000000000000000000",
    "Action": "GetAsset",
    "Version": "2024-01-01",
    "Service": "ark",
    "Region": "cn-beijing"
  },
  "Result": {
    "Id": "asset-20260813150950-t2tpn",
    "Name": "女主参考图",
    "URL": "https://example.com/temporary-asset-url",
    "AssetType": "Image",
    "GroupId": "group-20260813103927-cgxm6",
    "Status": "Active",
    "Moderation": {"Strategy": "Default"},
    "CreateTime": "2026-08-13T15:09:50Z",
    "UpdateTime": "2026-08-13T15:09:54Z",
    "LastInferenceTime": "2026-08-17T10:08:49+08:00",
    "ProjectName": "default"
  }
}`

const assetStatuses: StatusRow[] = [
  {
    name: 'processing / Processing',
    description: '素材正在处理。',
  },
  {
    name: 'Active',
    description: '素材可用。',
  },
  {
    name: 'Failed',
    description: '素材处理失败，可读取 `error`。',
  },
]

function GetAsset() {
  const baseUrl = useApiBaseUrl()
  return (
    <EndpointCard
      anchor='get-asset-endpoint'
      number='8'
      title='查询素材资产信息'
      method='POST'
      path='/api/assets/?Action=GetAsset&Version=2024-01-01'
      contentType='application/json'
      remark='根据素材 ID 获取详情及最新处理状态（由上游实时返回）。请求和响应遵循火山方舟 GetAsset 协议。'
    >
      <Block>
        <SubHeading text='请求示例' />
        <CodeSample
          language='bash'
          code={`curl --request POST \\
  "${baseUrl}/api/assets/?Action=GetAsset&Version=2024-01-01" \\
  --header "Authorization: Bearer $KT_API_KEY" \\
  --header "Content-Type: application/json" \\
  --data '{
    "Id": "asset-20260813150950-t2tpn",
    "ProjectName": "default"
  }'`}
        />
      </Block>
      <Block>
        <SubHeading text='请求字段' />
        <FieldTable
          rows={[
            {
              name: 'Id',
              type: 'string',
              required: '是',
              description: '当前用户拥有的素材 ID。',
            },
            {
              name: 'ProjectName',
              type: 'string',
              required: '否',
              description: '素材所属项目，默认 `default`。',
            },
          ]}
        />
        <Paragraph text='接口同时兼容 multipart 表单；使用 `curl --form` 时不要手动设置 `Content-Type`。' />
      </Block>
      <Block>
        <SubHeading text='成功响应' />
        <CodeSample language='json' code={getAssetJson} />
      </Block>
      <Block>
        <SubHeading text='响应字段' />
        <ResponseFieldTable
          rows={[
            {
              name: 'ResponseMetadata.RequestId',
              type: 'string',
              description: '上游请求 ID；渠道未提供时为空字符串。',
            },
            {
              name: 'ResponseMetadata.Action',
              type: 'string',
              description: '固定为 `GetAsset`。',
            },
            {
              name: 'ResponseMetadata.Version',
              type: 'string',
              description: '固定为 `2024-01-01`。',
            },
            {
              name: 'ResponseMetadata.Service',
              type: 'string',
              description: '固定为 `ark`。',
            },
            {
              name: 'ResponseMetadata.Region',
              type: 'string',
              description: '固定为 `cn-beijing`。',
            },
            {
              name: 'Result.Id',
              type: 'string',
              description: '素材 ID。',
            },
            {
              name: 'Result.Name',
              type: 'string',
              description: '素材名称。',
            },
            {
              name: 'Result.URL',
              type: 'string',
              description: '素材临时访问地址。',
            },
            {
              name: 'Result.GroupId',
              type: 'string',
              description: '素材组 ID。',
            },
            {
              name: 'Result.AssetType',
              type: 'string',
              description: 'Image / Video / Audio',
            },
            {
              name: 'Result.Status',
              type: 'string',
              description: '`Processing`、`Active` 或 `Failed`。',
            },
            {
              name: 'Result.Moderation.Strategy',
              type: 'string',
              description: '内容审核策略。',
            },
            {
              name: 'Result.Error.Code',
              type: 'string',
              description: '没有错误时省略。',
            },
            {
              name: 'Result.Error.Message',
              type: 'string',
              description: '没有错误时省略。',
            },
            {
              name: 'Result.CreateTime',
              type: 'string',
              description: '创建时间。',
            },
            {
              name: 'Result.UpdateTime',
              type: 'string',
              description: '更新时间。',
            },
            {
              name: 'Result.LastInferenceTime',
              type: 'string',
              description: '最近用于推理的时间，可能不返回。',
            },
            {
              name: 'Result.ProjectName',
              type: 'string',
              description: '素材所属项目。',
            },
          ]}
        />
        <AssetDetailNote />
      </Block>
      <Block>
        <SubHeading text='常见素材状态' />
        <StatusTable rows={assetStatuses} />
        <Paragraph text='素材处理失败仍属于成功查询，`Result.Status` 为 `Failed`，具体原因位于 `Result.Error.Code` 和 `Result.Error.Message`。素材不存在、不属于当前用户或查询服务失败时返回错误响应。' />
      </Block>
    </EndpointCard>
  )
}

const updateAssetGroupJson = `{
  "ResponseMetadata": {
    "RequestId": "20260328000000000000000000000000",
    "Action": "UpdateAssetGroup",
    "Version": "2024-01-01",
    "Service": "ark",
    "Region": "cn-beijing"
  },
  "Result": {
    "Id": "group-20260813103927-cgxm6"
  }
}`

function UpdateAssetGroup() {
  const baseUrl = useApiBaseUrl()
  return (
    <EndpointCard
      anchor='update-asset-group-endpoint'
      number='9'
      title='更新素材资产组合信息'
      method='POST'
      path='/api/assets/?Action=UpdateAssetGroup&Version=2024-01-01'
      contentType='application/json'
      remark='更新当前用户拥有的素材组名称或描述。请求和响应遵循火山方舟 UpdateAssetGroup 协议。'
    >
      <Block>
        <SubHeading text='请求示例' />
        <CodeSample
          language='bash'
          code={`curl --request POST \\
  "${baseUrl}/api/assets/?Action=UpdateAssetGroup&Version=2024-01-01" \\
  --header "Authorization: Bearer $KT_API_KEY" \\
  --header "Content-Type: application/json" \\
  --data '{
    "Id": "group-20260813103927-cgxm6",
    "Name": "人物素材组-更新",
    "Description": "更新后的人物参考素材",
    "ProjectName": "default"
  }'`}
        />
      </Block>
      <Block>
        <SubHeading text='请求字段' />
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className='w-[28%]'>字段</TableHead>
              <TableHead className='w-[16%]'>类型</TableHead>
              <TableHead className='w-[10%]'>必填</TableHead>
              <TableHead className='w-[10%]'>默认值</TableHead>
              <TableHead>说明</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  Id
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                是
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                无
              </TableCell>
              <TableCell className='text-sm'>
                当前用户拥有的素材组 ID。
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  Name
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                条件必填
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                无
              </TableCell>
              <TableCell className='text-sm'>
                <RichText text='新名称，最多 64 个 Unicode 字符。与 `Description` 至少填写一个。' />
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  Description
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                条件必填
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                无
              </TableCell>
              <TableCell className='text-sm'>
                <RichText text='新描述，最多 300 个 Unicode 字符。与 `Name` 至少填写一个。' />
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  ProjectName
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                否
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                default
              </TableCell>
              <TableCell className='text-sm'>
                素材组所属项目；本字段用于定位资源，不会改变项目归属。
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
        <Paragraph text='当前请求仅更新非空的 `Name` 和 `Description`。如需将名称或描述清空，该协议暂不支持。' />
      </Block>
      <Block>
        <SubHeading text='成功响应' />
        <CodeSample language='json' code={updateAssetGroupJson} />
        <ResponseFieldTable
          rows={[
            {
              name: 'ResponseMetadata.RequestId',
              type: 'string',
              description: '上游请求 ID；渠道未提供时为空字符串。',
            },
            {
              name: 'ResponseMetadata.Action',
              type: 'string',
              description: '固定为 `UpdateAssetGroup`。',
            },
            {
              name: 'ResponseMetadata.Version',
              type: 'string',
              description: '固定为 `2024-01-01`。',
            },
            {
              name: 'ResponseMetadata.Service',
              type: 'string',
              description: '固定为 `ark`。',
            },
            {
              name: 'ResponseMetadata.Region',
              type: 'string',
              description: '固定为 `cn-beijing`。',
            },
            {
              name: 'Result.Id',
              type: 'string',
              description: '更新成功的上游素材组 ID。',
            },
          ]}
        />
      </Block>
      <Block>
        <SubHeading text='失败响应' />
        <Paragraph text='素材组不存在、不属于当前用户、名称或描述超长、两个可更新字段均为空、项目名不匹配或上游更新失败时请求失败。' />
      </Block>
    </EndpointCard>
  )
}

const updateAssetJson = `{
  "ResponseMetadata": {
    "RequestId": "20260328000000000000000000000000",
    "Action": "UpdateAsset",
    "Version": "2024-01-01",
    "Service": "ark",
    "Region": "cn-beijing"
  },
  "Result": {
    "Id": "asset-20260817175008-llds7"
  }
}`

function UpdateAsset() {
  const baseUrl = useApiBaseUrl()
  return (
    <EndpointCard
      anchor='update-asset-endpoint'
      number='10'
      title='更新素材资产信息'
      method='POST'
      path='/api/assets/?Action=UpdateAsset&Version=2024-01-01'
      contentType='application/json'
      remark='更新当前用户拥有的素材名称。请求和响应遵循火山方舟 UpdateAsset 协议。'
    >
      <Block>
        <SubHeading text='请求示例' />
        <CodeSample
          language='bash'
          code={`curl --request POST \\
  "${baseUrl}/api/assets/?Action=UpdateAsset&Version=2024-01-01" \\
  --header "Authorization: Bearer $KT_API_KEY" \\
  --header "Content-Type: application/json" \\
  --data '{
    "Id": "asset-20260817175008-llds7",
    "Name": "更新后的素材名称",
    "ProjectName": "default"
  }'`}
        />
      </Block>
      <Block>
        <SubHeading text='请求字段' />
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className='w-[28%]'>字段</TableHead>
              <TableHead className='w-[16%]'>类型</TableHead>
              <TableHead className='w-[10%]'>必填</TableHead>
              <TableHead className='w-[10%]'>默认值</TableHead>
              <TableHead>说明</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  Id
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                是
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                无
              </TableCell>
              <TableCell className='text-sm'>当前用户拥有的素材 ID。</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  Name
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                是
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                无
              </TableCell>
              <TableCell className='text-sm'>
                素材的新名称，不能为空，最多 64 个 Unicode 字符。
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell>
                <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
                  ProjectName
                </code>
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                string
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                否
              </TableCell>
              <TableCell className='text-muted-foreground text-xs'>
                default
              </TableCell>
              <TableCell className='text-sm'>
                素材所属项目；本字段用于定位资源，不会改变项目归属。
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </Block>
      <Block>
        <SubHeading text='成功响应' />
        <CodeSample language='json' code={updateAssetJson} />
        <ResponseFieldTable
          rows={[
            {
              name: 'ResponseMetadata.RequestId',
              type: 'string',
              description: '上游请求 ID；渠道未提供时为空字符串。',
            },
            {
              name: 'ResponseMetadata.Action',
              type: 'string',
              description: '固定为 `UpdateAsset`。',
            },
            {
              name: 'ResponseMetadata.Version',
              type: 'string',
              description: '固定为 `2024-01-01`。',
            },
            {
              name: 'ResponseMetadata.Service',
              type: 'string',
              description: '固定为 `ark`。',
            },
            {
              name: 'ResponseMetadata.Region',
              type: 'string',
              description: '固定为 `cn-beijing`。',
            },
            {
              name: 'Result.Id',
              type: 'string',
              description: '更新成功的上游素材 ID。',
            },
          ]}
        />
      </Block>
      <Block>
        <SubHeading text='失败响应' />
        <Paragraph text='素材不存在、不属于当前用户、名称为空或超过 64 个字符、项目名不匹配或上游更新失败时请求失败。' />
      </Block>
    </EndpointCard>
  )
}

const deleteAssetJson = `{
  "ResponseMetadata": {
    "RequestId": "20260328000000000000000000000000",
    "Action": "DeleteAsset",
    "Version": "2024-01-01",
    "Service": "ark",
    "Region": "cn-beijing"
  },
  "Result": {}
}`

function DeleteAsset() {
  const baseUrl = useApiBaseUrl()
  return (
    <EndpointCard
      anchor='delete-asset-endpoint'
      number='11'
      title='删除单个素材资产'
      method='POST'
      path='/api/assets/?Action=DeleteAsset&Version=2024-01-01'
      contentType='application/json'
      remark='删除当前用户拥有的指定素材。请求和响应遵循火山方舟 DeleteAsset 协议。删除操作不可逆。'
    >
      <Block>
        <SubHeading text='请求示例' />
        <CodeSample
          language='bash'
          code={`curl --request POST \\
  "${baseUrl}/api/assets/?Action=DeleteAsset&Version=2024-01-01" \\
  --header "Authorization: Bearer $KT_API_KEY" \\
  --header "Content-Type: application/json" \\
  --data '{
    "Id": "asset-20260817175008-llds7",
    "ProjectName": "default"
  }'`}
        />
      </Block>
      <Block>
        <SubHeading text='请求字段' />
        <FieldTable
          rows={[
            {
              name: 'Id',
              type: 'string',
              required: '是',
              description: '当前用户拥有的素材 ID。',
            },
            {
              name: 'ProjectName',
              type: 'string',
              required: '否',
              description: '素材所属项目，默认 `default`。',
            },
          ]}
        />
      </Block>
      <Block>
        <SubHeading text='成功响应' />
        <CodeSample language='json' code={deleteAssetJson} />
        <ResponseFieldTable
          rows={[
            {
              name: 'ResponseMetadata.RequestId',
              type: 'string',
              description: '上游请求 ID；渠道未提供时为空字符串。',
            },
            {
              name: 'ResponseMetadata.Action',
              type: 'string',
              description: '固定为 `DeleteAsset`。',
            },
            {
              name: 'ResponseMetadata.Version',
              type: 'string',
              description: '固定为 `2024-01-01`。',
            },
            {
              name: 'ResponseMetadata.Service',
              type: 'string',
              description: '固定为 `ark`。',
            },
            {
              name: 'ResponseMetadata.Region',
              type: 'string',
              description: '固定为 `cn-beijing`。',
            },
            {
              name: 'Result',
              type: 'object',
              description: '空对象；接口没有业务返回参数。',
            },
          ]}
        />
      </Block>
      <Block>
        <SubHeading text='失败响应' />
        <Paragraph text='素材不存在、不属于当前用户、项目名不匹配或上游删除失败时请求失败。上游删除失败时，本地素材记录会保留，以便后续重试。' />
      </Block>
    </EndpointCard>
  )
}

const deleteAssetGroupJson = `{
  "ResponseMetadata": {
    "RequestId": "20260328000000000000000000000000",
    "Action": "DeleteAssetGroup",
    "Version": "2024-01-01",
    "Service": "ark",
    "Region": "cn-beijing"
  },
  "Result": {}
}`

function DeleteAssetGroup() {
  const baseUrl = useApiBaseUrl()
  return (
    <EndpointCard
      anchor='delete-asset-group-endpoint'
      number='12'
      title='删除指定素材组'
      method='POST'
      path='/api/assets/?Action=DeleteAssetGroup&Version=2024-01-01'
      contentType='application/json'
      remark='删除当前用户拥有的指定素材组，并批量删除组内所有素材。请求和响应遵循火山方舟 DeleteAssetGroup 协议。'
    >
      <Block>
        <SubHeading text='请求示例' />
        <CodeSample
          language='bash'
          code={`curl --request POST \\
  "${baseUrl}/api/assets/?Action=DeleteAssetGroup&Version=2024-01-01" \\
  --header "Authorization: Bearer $KT_API_KEY" \\
  --header "Content-Type: application/json" \\
  --data '{
    "Id": "group-20260817174619-glh22",
    "ProjectName": "default"
  }'`}
        />
      </Block>
      <Block>
        <SubHeading text='请求字段' />
        <FieldTable
          rows={[
            {
              name: 'Id',
              type: 'string',
              required: '是',
              description: '当前用户拥有的素材组 ID。',
            },
            {
              name: 'ProjectName',
              type: 'string',
              required: '否',
              description: '素材组所属项目，默认 `default`。',
            },
          ]}
        />
      </Block>
      <Block>
        <SubHeading text='风险说明' />
        <ul className='text-muted-foreground list-disc space-y-1 pl-6 text-sm leading-relaxed'>
          <li>删除素材组会同时删除组内所有素材，操作不可恢复。</li>
          <li>素材组包含较多素材时，上游删除可能耗费一定时间。</li>
          <li>真人素材组可能受上游授权状态限制，无法在授权有效期内删除。</li>
        </ul>
      </Block>
      <Block>
        <SubHeading text='成功响应' />
        <CodeSample language='json' code={deleteAssetGroupJson} />
        <ResponseFieldTable
          rows={[
            {
              name: 'ResponseMetadata.RequestId',
              type: 'string',
              description: '上游请求 ID；渠道未提供时为空字符串。',
            },
            {
              name: 'ResponseMetadata.Action',
              type: 'string',
              description: '固定为 `DeleteAssetGroup`。',
            },
            {
              name: 'ResponseMetadata.Version',
              type: 'string',
              description: '固定为 `2024-01-01`。',
            },
            {
              name: 'ResponseMetadata.Service',
              type: 'string',
              description: '固定为 `ark`。',
            },
            {
              name: 'ResponseMetadata.Region',
              type: 'string',
              description: '固定为 `cn-beijing`。',
            },
            {
              name: 'Result',
              type: 'object',
              description: '空对象；接口没有业务返回参数。',
            },
          ]}
        />
      </Block>
      <Block>
        <SubHeading text='失败响应' />
        <Paragraph text='素材组不存在、不属于当前用户、项目名不匹配、上游拒绝删除或素材服务异常时请求失败。' />
      </Block>
    </EndpointCard>
  )
}

const visualSessionJson = `{
  "BytedToken": "********22041234D89A4DDA08**************",
  "H5Link": "https://ark.volcengine.com/region:cn-beijing/mobile/livenees-face-manage/authorization?pl=********&uid=********",
  "CallbackURL": "https://www.example.com/callback"
}`

const visualCallbackParams: FieldRow[] = [
  {
    name: 'bytedToken',
    type: 'string',
    required: '是',
    description: '本次认证凭证，应与创建会话返回的 token 对应。',
  },
  {
    name: 'resultCode',
    type: 'string',
    required: '是',
    description: '认证结果码；`10000` 表示认证成功。',
  },
  {
    name: 'algorithmBaseRespCode',
    type: 'string',
    required: '否',
    description: '服务端子错误码。',
  },
  {
    name: 'reqMeasureInfoValue',
    type: 'string',
    required: '否',
    description: '是否计费：`0` 不计费，`1` 计费。',
  },
  {
    name: 'verify_type',
    type: 'string',
    required: '否',
    description: '认证类型，当前固定为 `real_time`。',
  },
]

function CreateVisualValidateSession() {
  const baseUrl = useApiBaseUrl()
  return (
    <EndpointCard
      anchor='create-visual-validate-session-endpoint'
      number='13'
      title='拉起真人认证 H5'
      method='POST'
      path='/api/assets/?Action=CreateVisualValidateSession&Version=2024-01-01'
      contentType='application/json'
      remark='创建真人认证会话，获取供终端客户打开的 H5 真人认证页面链接。'
    >
      <Block>
        <SubHeading text='请求示例' />
        <CodeSample
          language='bash'
          code={`curl --request POST \\
  "${baseUrl}/api/assets/?Action=CreateVisualValidateSession&Version=2024-01-01" \\
  --header "Authorization: Bearer $KT_API_KEY" \\
  --header "Content-Type: application/json" \\
  --data '{
    "CallbackURL": "https://www.example.com/callback",
    "ProjectName": "default"
  }'`}
        />
      </Block>
      <Block>
        <SubHeading text='请求字段' />
        <FieldTable
          rows={[
            {
              name: 'CallbackURL',
              type: 'string',
              required: '是',
              description: '认证结束后跳转的绝对 HTTP(S) URL。',
            },
            {
              name: 'ProjectName',
              type: 'string',
              required: '否',
              description: '资源所属项目，默认 `default`，大小写敏感。',
            },
          ]}
        />
      </Block>
      <Block>
        <SubHeading text='成功响应' />
        <Paragraph text='接口直接返回真人认证会话业务字段，不使用 `ResponseMetadata` 或 `Result` 包装：' />
        <CodeSample language='json' code={visualSessionJson} />
        <ResponseFieldTable
          rows={[
            {
              name: 'BytedToken',
              type: 'string',
              description:
                '本次认证的唯一凭证，用于后续查询认证结果；有效期为 30 分钟。',
            },
            {
              name: 'H5Link',
              type: 'string',
              description: '终端客户使用的真人认证页面链接，使用一次后失效。',
            },
            {
              name: 'CallbackURL',
              type: 'string',
              description: '认证完成后跳转的回调地址。',
            },
          ]}
        />
      </Block>
      <Block>
        <SubHeading text='H5 页面语言' />
        <Paragraph text='可在 `H5Link` 后追加或修改 `lng` 查询参数：' />
        <KeyValueTable
          kLabel='lng'
          vLabel='语言'
          rows={[
            {
              k: 'zh',
              v: '简体中文，默认值。',
            },
            {
              k: 'en',
              v: '英文。',
            },
            {
              k: 'zh-Hant',
              v: '繁体中文。',
            },
          ]}
        />
      </Block>
      <Block>
        <SubHeading text='回调参数' />
        <Paragraph text='认证完成后，H5 页面会跳转到类似以下地址：' />
        <CodeSample
          language='bash'
          code='https://www.example.com/callback?bytedToken=<TOKEN>&resultCode=10000&algorithmBaseRespCode=0&reqMeasureInfoValue=1&verify_type=real_time'
        />
        <FieldTable rows={visualCallbackParams} />
        <Note text='仅当 `resultCode=10000` 时，才应使用 `BytedToken` 查询认证结果。`BytedToken` 仅支持一次认证流程，并应在生成后 30 分钟内使用。' />
      </Block>
      <Block>
        <SubHeading text='失败响应' />
        <Paragraph text='`CallbackURL` 缺失或不是绝对 HTTP(S) URL、项目不可用、素材渠道不支持真人认证，或上游创建会话失败时请求失败。' />
      </Block>
    </EndpointCard>
  )
}

const visualResultJson = `{
  "GroupId": "group-20260331145705-*****"
}`

function GetVisualValidateResult() {
  const baseUrl = useApiBaseUrl()
  return (
    <EndpointCard
      anchor='get-visual-validate-result-endpoint'
      number='14'
      title='获取真人认证结果'
      method='POST'
      path='/api/assets/?Action=GetVisualValidateResult&Version=2024-01-01'
      contentType='application/json'
      remark='真人认证成功后，使用创建会话返回的 `BytedToken` 查询本次认证创建的真人素材组 ID。火山上游可能需要约 3 分钟生成素材组；在回调 `resultCode=10000` 后如尚未返回 `GroupId`，请间隔短暂时间重试。查询必须在 Token 生成后 30 分钟内完成。服务端只保存 Token 的 SHA-256 摘要，用于将查询路由回创建会话的同一素材渠道；不会保存或记录明文 Token。查询返回 `GroupId` 后，会将该素材组登记为当前用户的 `LivenessFace` 素材组，并删除会话摘要；尚未返回 `GroupId` 时会话保留，可继续查询。'
    >
      <Block>
        <SubHeading text='请求示例' />
        <CodeSample
          language='bash'
          code={`curl --request POST \\
  "${baseUrl}/api/assets/?Action=GetVisualValidateResult&Version=2024-01-01" \\
  --header "Authorization: Bearer $KT_API_KEY" \\
  --header "Content-Type: application/json" \\
  --data '{
    "BytedToken": "<BYTED_TOKEN>",
    "ProjectName": "default"
  }'`}
        />
      </Block>
      <Block>
        <SubHeading text='请求字段' />
        <FieldTable
          rows={[
            {
              name: 'BytedToken',
              type: 'string',
              required: '是',
              description:
                '`CreateVisualValidateSession` 返回的认证凭证，有效期为 30 分钟。',
            },
            {
              name: 'ProjectName',
              type: 'string',
              required: '否',
              description:
                '创建会话时的项目；必须与创建认证会话时使用的项目名称一致，大小写敏感。',
            },
          ]}
        />
        <Note text='不要在日志、文档或版本库中持久化真实 `BytedToken`。' />
      </Block>
      <Block>
        <SubHeading text='成功响应' />
        <Paragraph text='接口直接返回业务字段，不使用 `ResponseMetadata` 或 `Result` 包装：' />
        <CodeSample language='json' code={visualResultJson} />
        <ResponseFieldTable
          rows={[
            {
              name: 'GroupId',
              type: 'string',
              description: '本次真人认证创建的真人人像素材组 ID。',
            },
          ]}
        />
      </Block>
      <Block>
        <SubHeading text='调用流程' />
        <ol className='text-muted-foreground list-decimal space-y-2 pl-6 text-sm leading-relaxed'>
          <li>
            {'调用 '}
            <code className='bg-muted rounded px-1 py-0.5 text-xs'>
              CreateVisualValidateSession
            </code>
            {'，保存返回的 '}
            <code className='bg-muted rounded px-1 py-0.5 text-xs'>
              BytedToken
            </code>
            {'，并将 '}
            <code className='bg-muted rounded px-1 py-0.5 text-xs'>H5Link</code>
            {' 提供给终端客户。'}
          </li>
          <li>
            {'终端客户完成 H5 真人认证，页面跳转至 '}
            <code className='bg-muted rounded px-1 py-0.5 text-xs'>
              CallbackURL
            </code>
            。
          </li>
          <li>
            {'服务端校验回调参数中的 '}
            <code className='bg-muted rounded px-1 py-0.5 text-xs'>
              resultCode
            </code>
            。
          </li>
          <li>
            {'当 '}
            <code className='bg-muted rounded px-1 py-0.5 text-xs'>
              resultCode=10000
            </code>
            {' 时，等待约 3 分钟后调用本接口；未返回 '}
            <code className='bg-muted rounded px-1 py-0.5 text-xs'>
              GroupId
            </code>
            {' 时可间隔短暂时间继续轮询，直至 Token 过期。'}
          </li>
          <li>
            {'使用返回的 '}
            <code className='bg-muted rounded px-1 py-0.5 text-xs'>
              GroupId
            </code>
            {' 上传或查询真人素材。'}
          </li>
        </ol>
      </Block>
      <Block>
        <SubHeading text='失败响应' />
        <Paragraph text='`BytedToken` 缺失、已过期、已使用或不属于当前用户（`BytedToken is missing, expired, used, or does not belong to current user`），绑定渠道已停用，或上游查询失败时请求失败。' />
      </Block>
    </EndpointCard>
  )
}

function UserIsolation() {
  return (
    <BulletList
      items={[
        '视频任务 `id` 为本平台生成的 `task_` 开头 ID，上游任务 ID 不会返回给调用方；Draft 任务引用（`draft_task.id`）也使用本平台任务 ID。',
        '素材组 `group-*` 与素材 `asset-*` ID 由上游渠道生成并原样返回；本平台记录每个 ID 的归属用户和所在渠道。',
        '每次查询、更新、删除或在视频任务中引用（`asset://<ASSET_ID>`）时，都会校验资源是否属于当前 API Key 对应的用户，并固定路由到资源所在渠道；不属于当前用户时返回 `record not found` 或 `asset not found or not owned by you`。',
        '同一个视频任务引用的素材必须属于同一渠道，否则返回 `assets must belong to the same channel`。',
        '用户不能查询或使用其他用户的任务、素材组或素材。',
      ]}
    />
  )
}

const commonErrors: ErrorRow[] = [
  {
    scenario: 'API Key 无效（视频与素材接口）',
    message:
      'HTTP 401 / {"error": {"code": "", "message": "Invalid token (request id: ...)", "type": "new_api_error"}}',
  },
  {
    scenario: '视频请求 JSON 无效或缺少字段',
    message: 'HTTP 400 / invalid_request（例如 content is required）',
  },
  {
    scenario: '视频模型不可用或填写了 Endpoint ID',
    message: 'HTTP 400 / model "..." is not served by this plugin',
  },
  {
    scenario: 'callback_url 不合法',
    message:
      'HTTP 400 / invalid_request: callback_url is not allowed: ...（不扣费、不创建任务）',
  },
  {
    scenario: '视频任务被上游拒绝',
    message: 'HTTP 400 / invalid_request，message 为上游原始错误响应',
  },
  {
    scenario: '视频任务不存在或不属于当前用户',
    message: 'HTTP 404 / task_not_found: Task not found',
  },
  {
    scenario: '引用的素材不存在或不属于当前用户',
    message: 'HTTP 400 / asset_not_found: asset not found or not owned by you',
  },
  {
    scenario: 'Draft 任务不存在或不属于当前用户',
    message: 'HTTP 400 / origin task not found or not owned by you',
  },
  {
    scenario: '素材接口 Action 或 Version 不支持',
    message:
      'HTTP 400 / InvalidParameter: Action is not supported / Version must be 2024-01-01',
  },
  {
    scenario: '素材请求 JSON 无效',
    message: 'HTTP 400 / InvalidParameter: invalid JSON body',
  },
  {
    scenario: '素材组或素材不存在或不属于当前用户',
    message: 'HTTP 400 / InvalidParameter: record not found',
  },
  {
    scenario: '缺少素材组 ID',
    message: 'HTTP 400 / InvalidParameter: GroupId is required',
  },
  {
    scenario: '素材组筛选条件缺失',
    message: 'HTTP 400 / InvalidParameter: Filter.GroupType is required',
  },
  {
    scenario: '素材列表筛选条件缺失',
    message:
      'HTTP 400 / InvalidParameter: Filter.GroupIds must contain at least one asset group id',
  },
  {
    scenario: '素材所在渠道已停用',
    message: 'InvalidParameter: asset channel is disabled',
  },
  {
    scenario: '账户余额不足',
    message: 'InsufficientQuota',
  },
  {
    scenario: '上游不可用',
    message: 'HTTP 502 / UpstreamError',
  },
  {
    scenario: '真人认证 Token 无效/过期',
    message:
      'InvalidParameter: BytedToken is missing, expired, used, or does not belong to current user',
  },
  {
    scenario: '素材参数缺失或不合法（上游校验）',
    message:
      'HTTP 400 / MissingParameter.<Field> / InvalidParameter.<Field>（上游原样返回）',
  },
]

function CommonErrors() {
  return (
    <Block>
      <Paragraph text='视频接口错误使用 `{"error": {"code", "message"}}` 信封；素材接口错误使用 `ResponseMetadata.Error` 信封（本平台生成的错误 `Code` 为 `InvalidParameter`、`InsufficientQuota` 或 `UpstreamError`，上游错误原样返回）。API Key 认证失败时两类接口都返回 HTTP `401` 的 `{"error": {...}}` 信封。' />
      <ErrorTable rows={commonErrors} />
    </Block>
  )
}

const recommendations = [
  {
    title: '创建视频生成任务',
    body: '提交成功后保存响应中的 `id`。创建接口只表示任务已被接受，不表示视频已经可用；应立即进入任务查询流程。',
  },
  {
    title: '获取视频任务状态',
    body: '优先使用 `callback_url` 接收状态推送，并以周期性查询兜底（建议间隔不少于 5 秒）。收到 `succeeded` 后尽快下载 `content.video_url`（签名地址会过期）；收到 `failed` 或 `expired` 后读取 `error.code` 和 `error.message`，停止轮询。对网络或服务失败使用有限次数的退避重试。',
  },
  {
    title: '回调接收',
    body: '回调接收端应在 5 秒内返回 2xx，并按 `X-KtApi-Event-Id` 去重；开启签名时校验 `X-KtApi-Signature`。回调只是通知，最终状态以查询接口为准。',
  },
  {
    title: '计费',
    body: '视频任务在提交时预扣费用，完成后按实际用量结算；任务失败或超时会退还预扣费用。素材接口按次计费，调用失败时退还。',
  },
  {
    title: '创建素材组',
    body: '在首次为某一组素材上传前调用一次，并保存 `Result.Id`。同一业务素材集合可复用该 ID，不需要为每个素材重复创建组。',
  },
  {
    title: '创建素材资产',
    body: '创建成功后保存 `Result.Id`。创建接口是异步接口，应继续调用素材详情接口查询处理状态，而不是立即依赖素材可用。',
  },
  {
    title: '查询素材资产组合列表',
    body: '需要按条件筛选或展示素材组列表时使用。`Filter.GroupType` 为必填项；使用分页参数稳定翻页。',
  },
  {
    title: '查询单个素材资产组合信息',
    body: '需要获取单个素材组详情时使用。本平台校验归属后路由到素材组所在渠道查询。',
  },
  {
    title: '查询素材资产列表',
    body: '需要展示或批量检查素材时使用。`Filter.GroupIds` 至少包含一个 ID，且同一次查询中的素材组应属于同一渠道；列表用于概览，失败详情以素材详情接口为准。',
  },
  {
    title: '查询素材资产信息',
    body: '在需要确认单个素材是否可用时使用。`Result.Status` 为 `Active` 时可继续后续业务；为 `Processing` 时稍后再查；为 `Failed` 时读取 `Result.Error.Code` 和 `Result.Error.Message`。',
  },
  {
    title: '更新素材资产组合信息',
    body: '仅在需要修改素材组名称或描述时使用。`Name` 和 `Description` 至少填写一个；更新失败时本地记录保持不变。',
  },
  {
    title: '更新素材资产信息',
    body: '仅用于修改素材名称。`Name` 不能为空且不超过 64 字符；更新失败时本地记录保持不变。',
  },
  {
    title: '删除单个素材资产',
    body: '删除操作不可逆。上游删除失败时，本地素材记录会保留，以便后续重试。',
  },
  {
    title: '删除指定素材组',
    body: '删除素材组会同时删除组内所有素材，操作不可恢复。建议在删除前先查询确认素材组内容。',
  },
  {
    title: '拉起真人认证 H5',
    body: '创建会话后及时将 `H5Link` 交付给终端客户。认证完成后使用 `BytedToken` 和 `resultCode` 查询结果；`BytedToken` 有效期 30 分钟且仅支持一次认证流程。',
  },
  {
    title: '获取真人认证结果',
    body: '在终端客户完成真人认证且 `resultCode=10000` 后 30 分钟内调用。上游可能需要约 3 分钟生成素材组，未返回 `GroupId` 时可间隔短暂时间轮询。使用返回的 `GroupId` 可直接上传或查询真人素材；`BytedToken` 仅支持一次查询，切勿重复使用。',
  },
]

export function ApiDocsPage() {
  return (
    <DocsPageLayout
      sections={SECTIONS}
      title='接口文档'
      description='本文档描述了包括视频生成和素材创建等任务的相关接口协议'
      headerExtra={
        <p className='text-muted-foreground text-sm leading-relaxed'>
          使用 OpenAI 兼容格式（
          <code className='bg-muted rounded px-1.5 py-0.5 text-xs'>
            /v1/videos
          </code>
          ）调用视频生成，请参阅{' '}
          <Link
            to='/docs/kt'
            className='text-primary underline-offset-4 hover:underline'
          >
            KT API 视频生成接口文档
          </Link>
          。
        </p>
      }
    >
      <DocsSection id='base-url' title='BASE_URL'>
        <BaseUrlBlock />
      </DocsSection>
      <DocsSection id='api-overview' title='接口列表'>
        <EndpointOverviewTable rows={ENDPOINTS} />
      </DocsSection>
      <DocsSection id='create-video-task' title='1. 创建视频生成任务'>
        <CreateVideoTask />
      </DocsSection>
      <DocsSection id='get-video-task' title='2. 获取视频任务状态'>
        <GetVideoTask />
      </DocsSection>
      <DocsSection id='create-asset-group' title='3. 创建素材组'>
        <CreateAssetGroup />
      </DocsSection>
      <DocsSection id='create-asset' title='4. 创建素材资产'>
        <CreateAsset />
      </DocsSection>
      <DocsSection id='list-asset-groups' title='5. 查询素材资产组合列表'>
        <ListAssetGroups />
      </DocsSection>
      <DocsSection id='get-asset-group' title='6. 查询单个素材资产组合信息'>
        <GetAssetGroup />
      </DocsSection>
      <DocsSection id='list-assets' title='7. 查询素材资产列表'>
        <ListAssets />
      </DocsSection>
      <DocsSection id='get-asset' title='8. 查询素材资产信息'>
        <GetAsset />
      </DocsSection>
      <DocsSection id='update-asset-group' title='9. 更新素材资产组合信息'>
        <UpdateAssetGroup />
      </DocsSection>
      <DocsSection id='update-asset' title='10. 更新素材资产信息'>
        <UpdateAsset />
      </DocsSection>
      <DocsSection id='delete-asset' title='11. 删除单个素材资产'>
        <DeleteAsset />
      </DocsSection>
      <DocsSection id='delete-asset-group' title='12. 删除指定素材组'>
        <DeleteAssetGroup />
      </DocsSection>
      <DocsSection
        id='create-visual-validate-session'
        title='13. 拉起真人认证 H5'
      >
        <CreateVisualValidateSession />
      </DocsSection>
      <DocsSection id='get-visual-validate-result' title='14. 获取真人认证结果'>
        <GetVisualValidateResult />
      </DocsSection>
      <DocsSection id='user-isolation' title='用户隔离与 ID 规则'>
        <UserIsolation />
      </DocsSection>
      <DocsSection id='common-errors' title='常见错误'>
        <CommonErrors />
      </DocsSection>
      <DocsSection id='usage-recommendations' title='调用建议'>
        <RecommendationList items={recommendations} />
      </DocsSection>
    </DocsPageLayout>
  )
}
