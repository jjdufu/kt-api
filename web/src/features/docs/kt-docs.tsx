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
  BaseUrlBlock,
  Block,
  CodeSample,
  DocsPageLayout,
  DocsSection,
  EndpointCard,
  EndpointOverviewTable,
  ErrorTable,
  FieldTable,
  Note,
  Paragraph,
  RecommendationList,
  ResponseFieldTable,
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
  { id: 'download-video', title: '3. 下载视频文件' },
  { id: 'common-errors', title: '常见错误' },
  { id: 'usage-recommendations', title: '调用建议' },
]

const modelList = SEEDANCE_MODELS.map((model) => `\`${model}\``).join('、')

const createFields: FieldRow[] = [
  {
    name: 'model',
    type: 'string',
    required: '是',
    description: `视频模型 ID，当前可用：${modelList}。具体可用模型以「模型广场」为准；未开通的模型会返回 HTTP 503（\`model_not_found\`）。`,
  },
  {
    name: 'prompt',
    type: 'string',
    required: '否',
    description:
      '视频内容描述文本，会作为一个 `text` 内容项提交。纯文本生成视频时必填；携带参考素材时可省略。',
  },
  {
    name: 'images',
    type: 'array<string>',
    required: '否',
    description:
      '参考图片 URL 列表（也可使用 `asset://<ASSET_ID>`），每一项会转换为一个 `image_url` 内容项。携带任意参考素材时任务按图/参考生视频处理，否则按文生视频处理。',
  },
  {
    name: 'seconds',
    type: 'number',
    required: '否',
    description:
      '视频时长（秒），`duration` 为同义字段。网关只校验取值为 `-1` 或 1–3600，实际可用范围由模型决定（Seedance 2.5：4–30 或 `-1`；Seedance 2.0 系列：4–15 或 `-1`），超出范围由上游返回错误。省略时使用模型默认时长。',
  },
  {
    name: 'metadata',
    type: 'object',
    required: '否',
    description:
      '原样合并进 Seedance 原生请求体，用于传递 `ratio`（画面比例）、`resolution`（清晰度）、`generate_audio`、`watermark`、`seed` 以及视频/音频参考（`content` 数组）等字段，字段含义与取值见「接口文档」中的创建视频生成任务。',
  },
]

const createSuccessJson = `{
  "id": "task_GlvqIS2axgKamyWpKzSP4yuV0hdf0q78",
  "object": "video",
  "model": "doubao-seedance-2-0-mini-260615",
  "status": "queued",
  "progress": 0,
  "created_at": 1792410250
}`

const createResponseFields: ResponseFieldRow[] = [
  {
    name: 'id',
    type: 'string',
    description:
      '任务 ID（`task_` 开头，由本平台生成）。请原样保存该值，后续查询状态和下载视频都需要使用它。',
  },
  {
    name: 'object',
    type: 'string',
    description: '固定为 `video`。',
  },
  {
    name: 'model',
    type: 'string',
    description: '请求中传入的模型 ID。',
  },
  {
    name: 'status',
    type: 'string',
    description: '任务创建后的初始状态，为 `queued`。',
  },
  {
    name: 'progress',
    type: 'integer',
    description: '任务进度百分比，创建时为 `0`。',
  },
  {
    name: 'created_at',
    type: 'integer',
    description: '任务创建时间（Unix 秒）。',
  },
]

const createErrorJson = `{
  "code": "fail_to_fetch_task",
  "message": "{\\"error\\": {\\"code\\": \\"InvalidParameter\\", \\"message\\": \\"The parameter \`ratio\` specified in the request is not valid.\\", ...}}",
  "data": null
}`

function CreateVideoTask() {
  const baseUrl = useApiBaseUrl()
  return (
    <EndpointCard
      anchor='create-video-task-endpoint'
      number='1'
      title='创建视频生成任务'
      method='POST'
      path='/v1/videos'
      contentType='application/json（也接受 multipart/form-data 文本字段）'
      remark='该接口用于提交一个异步视频生成任务，创建成功只表示任务已被接受，不代表视频已经生成完成，请配合下方“获取视频任务状态”接口轮询结果。'
    >
      <Block>
        <SubHeading>请求参数</SubHeading>
        <FieldTable rows={createFields} />
        <Note text='视频、音频参考素材请放在 `metadata.content` 中，使用 Seedance 原生内容项，例如 `{"type": "video_url", "video_url": {"url": "..."}, "role": "reference_video"}`；图片参考可使用 `images` 字段。multipart 请求只接受文本字段，不支持直接上传文件。' />
        <Paragraph text='本接口不支持顶层 `size` 字段（会被忽略），请使用 `metadata.ratio` 与 `metadata.resolution` 指定画面比例与清晰度；`metadata.callback_url` 在本接口上不会生效，需要任务回调时请使用「接口文档」中的 `/api/v3/contents/generations/tasks`。' />
      </Block>
      <Block>
        <SubHeading>请求示例（纯文本生成视频）</SubHeading>
        <CodeSample
          language='bash'
          code={`curl --location --request POST '${baseUrl}/v1/videos' \\
  --header "Authorization: Bearer $KT_API_KEY" \\
  --header 'Content-Type: application/json' \\
  --data-raw '{
    "model": "${SEEDANCE_MODELS[0]}",
    "prompt": "a cat walking in the rain"
  }'`}
        />
      </Block>
      <Block>
        <SubHeading>请求示例（携带参考素材生成视频）</SubHeading>
        <CodeSample
          language='bash'
          code={`curl --location --request POST '${baseUrl}/v1/videos' \\
  --header "Authorization: Bearer $KT_API_KEY" \\
  --header 'Content-Type: application/json' \\
  --data-raw '{
    "model": "${SEEDANCE_MODELS[0]}",
    "prompt": "animate these references",
    "images": ["https://example.com/ref.png"],
    "seconds": 10,
    "metadata": {
      "ratio": "16:9",
      "resolution": "720p",
      "content": [
        {"type": "video_url", "video_url": {"url": "https://example.com/clip.mp4"}, "role": "reference_video"},
        {"type": "audio_url", "audio_url": {"url": "https://example.com/track.mp3"}, "role": "reference_audio"}
      ]
    }
  }'`}
        />
      </Block>
      <Block>
        <SubHeading>响应示例</SubHeading>
        <CodeSample language='json' code={createSuccessJson} />
        <ResponseFieldTable rows={createResponseFields} />
      </Block>
      <Block>
        <SubHeading>创建失败响应</SubHeading>
        <Paragraph text='请求未通过本平台校验时返回 HTTP `400`（缺少 `model`、`seconds` 越界等，错误信封为 `{"error": {"message", "type", "code"}}`）；上游判定请求不合法时同样返回 HTTP `400`，`message` 中包含上游给出的原始错误文本，示例如下。不同请求的具体内容可能不同，不建议按固定字符串做匹配。' />
        <CodeSample language='json' code={createErrorJson} />
      </Block>
    </EndpointCard>
  )
}

const runningJson = `{
  "created_at": 1792410250,
  "id": "task_GlvqIS2axgKamyWpKzSP4yuV0hdf0q78",
  "model": "doubao-seedance-2-0-mini-260615",
  "object": "video",
  "progress": 50,
  "status": "in_progress"
}`

const completedJson = `{
  "completed_at": 1792410260,
  "created_at": 1792410250,
  "id": "task_GlvqIS2axgKamyWpKzSP4yuV0hdf0q78",
  "model": "doubao-seedance-2-0-mini-260615",
  "object": "video",
  "progress": 100,
  "status": "completed"
}`

const failedJson = `{
  "created_at": 1792410250,
  "error": {
    "code": "InputTextSensitiveContentDetected",
    "message": "The request failed because the input text may contain sensitive information."
  },
  "id": "task_CBkXCGoJ5l8VezyE5oHwlQwFEVJf5BxV",
  "model": "doubao-seedance-2-0-mini-260615",
  "object": "video",
  "progress": 100,
  "status": "failed"
}`

const taskStatuses: StatusRow[] = [
  { name: 'queued', description: '任务排队中，尚未开始处理。' },
  { name: 'in_progress', description: '任务正在生成中。' },
  {
    name: 'completed',
    description: '任务已完成，可通过“下载视频文件”接口获取视频。',
  },
  {
    name: 'failed',
    description:
      '任务失败，原因见 `error` 字段。上游超时（`expired`）、取消（`cancelled`）以及本平台判定的超时在本接口中也显示为 `failed`。',
  },
]

const taskResponseFields: ResponseFieldRow[] = [
  {
    name: 'id',
    type: 'string',
    description: '任务 ID，与创建任务时返回的 `id` 相同。',
  },
  { name: 'object', type: 'string', description: '固定为 `video`。' },
  { name: 'model', type: 'string', description: '创建任务时使用的模型 ID。' },
  { name: 'status', type: 'string', description: '本次查询得到的任务状态。' },
  {
    name: 'progress',
    type: 'integer',
    description: '任务进度百分比（排队 0、处理中 50、结束 100）。',
  },
  {
    name: 'error',
    type: 'object',
    description:
      '仅失败时返回，包含 `code` 与 `message`；超时或取消的任务可能不带该字段。',
  },
  {
    name: 'created_at',
    type: 'integer',
    description: '任务创建时间（Unix 秒）。',
  },
  {
    name: 'completed_at',
    type: 'integer',
    description: '任务最近更新时间（Unix 秒），任务结束后即完成时间。',
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
      path='/v1/videos/{task_id}'
      remark='返回本平台后台从上游同步到的最新任务状态，仅能查询当前账号创建的任务。本接口的响应不包含视频地址，请在任务完成后调用“下载视频文件”接口，或使用「接口文档」中的 `GET /api/v3/contents/generations/tasks/{task_id}` 读取 `content.video_url`（两种接口的任务 ID 通用）。'
    >
      <Block>
        <SubHeading>请求示例</SubHeading>
        <CodeSample
          language='bash'
          code={`curl --location --request GET '${baseUrl}/v1/videos/<task_id>' \\
  --header "Authorization: Bearer $KT_API_KEY"`}
        />
      </Block>
      <Block>
        <SubHeading>路径参数</SubHeading>
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
        <SubHeading>处理中响应</SubHeading>
        <CodeSample language='json' code={runningJson} />
      </Block>
      <Block>
        <SubHeading>完成响应</SubHeading>
        <CodeSample language='json' code={completedJson} />
      </Block>
      <Block>
        <SubHeading>任务状态</SubHeading>
        <StatusTable rows={taskStatuses} />
        <ResponseFieldTable rows={taskResponseFields} />
      </Block>
      <Block>
        <SubHeading>任务失败响应</SubHeading>
        <Paragraph text='查询接口对失败任务同样返回 HTTP `200`；请以响应体中的 `status`/`error` 字段判断任务是否失败，不要依赖 HTTP 状态码。' />
        <CodeSample language='json' code={failedJson} />
        <Paragraph text='建议间隔数秒轮询，直到 `status` 为 `completed` 或 `failed`。' />
      </Block>
    </EndpointCard>
  )
}

function DownloadVideo() {
  const baseUrl = useApiBaseUrl()
  return (
    <EndpointCard
      anchor='download-video-endpoint'
      number='3'
      title='下载视频文件'
      method='GET'
      path='/v1/videos/{task_id}/content'
      remark='任务状态为 `completed` 后，通过本接口下载生成的视频文件，由本平台从上游获取后返回给调用方。'
    >
      <Block>
        <SubHeading>请求示例</SubHeading>
        <CodeSample
          language='bash'
          code={`curl --location --request GET '${baseUrl}/v1/videos/<task_id>/content' \\
  --header "Authorization: Bearer $KT_API_KEY" \\
  --output video.mp4`}
        />
      </Block>
      <Block>
        <SubHeading>未完成时的响应</SubHeading>
        <Paragraph>任务尚未完成或已失败时返回 HTTP `400`：</Paragraph>
        <CodeSample
          language='json'
          code={`{
  "error": {
    "message": "Task is not completed yet, current status: IN_PROGRESS",
    "type": "invalid_request_error"
  }
}`}
        />
      </Block>
    </EndpointCard>
  )
}

const commonErrors: ErrorRow[] = [
  {
    scenario: 'API Key 无效或缺失',
    message: 'HTTP 401，Invalid token (request id: ...)',
  },
  {
    scenario: '未传入 `model`',
    message: 'Model name not specified, model name cannot be empty',
  },
  {
    scenario: '模型未开通或当前分组无可用渠道',
    message: 'HTTP 503，model_not_found：No available channel for model ...',
  },
  {
    scenario: '`seconds`/`duration` 不是 `-1` 且不在 1–3600 之间',
    message: 'seconds must be -1 or between 1 and 3600',
  },
  {
    scenario: '`metadata.resolution` 不在所选模型支持范围内',
    message: '... resolution must be one of 480p, 720p ...',
  },
  {
    scenario: '`metadata` 在 multipart 请求中不是合法的 JSON 对象字符串',
    message: 'metadata must be a JSON object string',
  },
  {
    scenario: '创建阶段上游业务校验失败',
    message: 'fail_to_fetch_task：透传上游给出的原始错误文本',
  },
  {
    scenario: '查询的任务不存在或不属于当前账号',
    message: 'HTTP 400，task_not_exist',
  },
  {
    scenario: '任务未完成时下载视频',
    message: 'Task is not completed yet, current status: ...',
  },
]

const recommendations = [
  {
    title: '创建视频生成任务',
    body: '创建成功只表示任务已被提交，不代表视频已经生成完成；请保存响应中的 `id` 并立即进入轮询流程。',
  },
  {
    title: '获取视频任务状态',
    body: '使用创建响应中的 `id` 周期性轮询。收到 `completed` 后调用下载接口获取视频；收到 `failed` 后读取 `error` 并停止轮询。查询结果来自本平台对上游的后台同步，状态变化可能有数秒延迟。',
  },
  {
    title: '参考素材的传递方式',
    body: '图片参考使用 `images` 字段；视频、音频参考通过 `metadata.content` 数组以 Seedance 原生内容项传递。只要请求携带任意参考素材，任务即按参考生视频处理，无需额外声明。',
  },
  {
    title: '画面比例、时长与清晰度',
    body: '画面比例与清晰度通过 `metadata.ratio`、`metadata.resolution` 指定，时长通过 `seconds` 指定；可选值由模型决定，详见「接口文档」。清晰度超出模型支持范围会在创建阶段直接返回 HTTP 400，不会扣费。',
  },
  {
    title: '计费',
    body: '提交时按请求参数预扣费用（未指定时长时按模型最长时长预估），任务成功后按实际用量结算多退少补；任务失败或超时会退还预扣费用；创建失败不扣费。',
  },
]

export function KtDocsPage() {
  return (
    <DocsPageLayout
      sections={SECTIONS}
      title='KT API 视频生成接口文档'
      description='本文档描述了 KT API OpenAI 兼容视频生成接口（/v1/videos）的相关接口协议，底层模型为 Seedance 系列。'
      headerExtra={
        <p className='text-muted-foreground text-sm leading-relaxed'>
          Seedance 原生接口、任务回调与素材接口请参阅{' '}
          <Link
            to='/docs'
            className='text-primary underline-offset-4 hover:underline'
          >
            接口文档
          </Link>
          。
        </p>
      }
    >
      <DocsSection id='base-url' title='BASE_URL'>
        <BaseUrlBlock />
      </DocsSection>
      <DocsSection id='api-overview' title='接口列表'>
        <EndpointOverviewTable
          rows={[
            { feature: '创建视频生成任务', method: 'POST', path: '/v1/videos' },
            {
              feature: '获取视频任务状态',
              method: 'GET',
              path: '/v1/videos/{task_id}',
            },
            {
              feature: '下载视频文件',
              method: 'GET',
              path: '/v1/videos/{task_id}/content',
            },
          ]}
        />
      </DocsSection>
      <DocsSection id='create-video-task' title='1. 创建视频生成任务'>
        <CreateVideoTask />
      </DocsSection>
      <DocsSection id='get-video-task' title='2. 获取视频任务状态'>
        <GetVideoTask />
      </DocsSection>
      <DocsSection id='download-video' title='3. 下载视频文件'>
        <DownloadVideo />
      </DocsSection>
      <DocsSection id='common-errors' title='常见错误'>
        <ErrorTable rows={commonErrors} />
      </DocsSection>
      <DocsSection id='usage-recommendations' title='调用建议'>
        <RecommendationList items={recommendations} />
      </DocsSection>
    </DocsPageLayout>
  )
}
