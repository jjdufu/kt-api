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
import { Fragment, useEffect, useState, type ReactNode } from 'react'
import { useTranslation } from 'react-i18next'

import {
  CodeBlock,
  CodeBlockCopyButton,
} from '@/components/ai-elements/code-block'
import { PublicLayout } from '@/components/layout'
import { Card, CardContent } from '@/components/ui/card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { cn } from '@/lib/utils'

import '@/styles/comic/comic-base.css'

import { useApiBaseUrl } from '../hooks/use-api-base-url'

import './docs-comic.css'

export type DocsSectionItem = { id: string; title: string }

export type FieldRow = {
  name: string
  type: string
  required: string
  description: string
}

export type ResponseFieldRow = {
  name: string
  type: string
  description: string
}

export type StatusRow = { name: string; description: string }

export type ErrorRow = { scenario: string; message: string }

export type KeyValueRow = { k: string; v: string }

const INLINE_CODE_CLASS =
  'dc-inline-code bg-muted rounded px-1.5 py-0.5 text-xs'

function InlineCode(props: { children: ReactNode }) {
  return <code className={INLINE_CODE_CLASS}>{props.children}</code>
}

/** Renders `code` spans in plain doc text as inline code. */
type RichTextSegment = { key: number; text: string; code: boolean }

function splitRichText(text: string): RichTextSegment[] {
  const parts = text.split('`')
  const segments: RichTextSegment[] = []
  let offset = 0
  parts.forEach((part, index) => {
    const code = index % 2 === 1 && index < parts.length - 1
    segments.push({ key: offset, text: part, code })
    offset += part.length + 1
  })
  return segments
}

export function RichText(props: { text: string }) {
  if (props.text.split('`').length < 3) return props.text
  return splitRichText(props.text).map((segment) =>
    segment.code ? (
      <code key={segment.key} className={INLINE_CODE_CLASS}>
        {segment.text}
      </code>
    ) : (
      <Fragment key={segment.key}>{segment.text}</Fragment>
    )
  )
}

export function DocsSection(props: {
  id: string
  title: ReactNode
  description?: ReactNode
  children: ReactNode
  className?: string
}) {
  return (
    <section
      id={props.id}
      data-docs-section
      className={cn('dc-section scroll-mt-24 space-y-6', props.className)}
    >
      <div className='space-y-2'>
        <h2 className='dc-h2 flex items-center gap-2 text-2xl font-semibold tracking-tight'>
          <span>{props.title}</span>
        </h2>
        {props.description ? (
          <p className='text-muted-foreground text-sm leading-relaxed'>
            {props.description}
          </p>
        ) : null}
      </div>
      <div className='space-y-6'>{props.children}</div>
    </section>
  )
}

function DocsToc(props: { sections: DocsSectionItem[]; activeId: string }) {
  const { t } = useTranslation()
  return (
    <nav className='dc-toc space-y-1 text-sm'>
      <p className='dc-toc-label text-muted-foreground mb-3 text-xs font-medium tracking-wider uppercase'>
        {t('On this page')}
      </p>
      {props.sections.map((section) => (
        <a
          key={section.id}
          href={`#${section.id}`}
          data-active={props.activeId === section.id ? '' : undefined}
          className={cn(
            'dc-toc-link block rounded-md px-3 py-1.5 transition-colors',
            props.activeId === section.id
              ? 'bg-primary/10 text-primary font-medium'
              : 'text-muted-foreground hover:bg-muted hover:text-foreground'
          )}
        >
          {section.title}
        </a>
      ))}
    </nav>
  )
}

/** Public docs page shell: sticky table of contents plus the content column. */
export function DocsPageLayout(props: {
  sections: DocsSectionItem[]
  title: ReactNode
  description: ReactNode
  headerExtra?: ReactNode
  children: ReactNode
}) {
  const { sections } = props
  const [activeId, setActiveId] = useState(sections[0]?.id ?? '')

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        for (const entry of entries) {
          if (entry.isIntersecting) setActiveId(entry.target.id)
        }
      },
      { rootMargin: '-20% 0px -70% 0px', threshold: 0 }
    )
    for (const section of sections) {
      const element = document.querySelector(`#${CSS.escape(section.id)}`)
      if (element) observer.observe(element)
    }
    return () => observer.disconnect()
  }, [sections])

  return (
    <div data-docs=''>
      <PublicLayout showMainContainer={false}>
        <div className='dc-page bg-background text-foreground pt-20'>
          <div className='mx-auto flex w-full max-w-7xl gap-8 px-4 py-8 md:px-6'>
            <aside className='hidden w-60 shrink-0 md:block'>
              <div className='sticky top-24'>
                <DocsToc sections={sections} activeId={activeId} />
              </div>
            </aside>
            <main className='min-w-0 flex-1 space-y-12 pb-24'>
              <header className='dc-hero space-y-3'>
                <span className='dc-kicker' aria-hidden='true'>
                  DOCS!
                </span>
                <h1 className='dc-h1 text-3xl font-semibold tracking-tight'>
                  {props.title}
                </h1>
                <p className='dc-hero-desc text-muted-foreground text-sm leading-relaxed'>
                  {props.description}
                </p>
                {props.headerExtra}
              </header>
              {props.children}
            </main>
          </div>
        </div>
      </PublicLayout>
    </div>
  )
}

export function CodeSample(props: { code: string; language: string }) {
  return (
    <div className='dc-code'>
      <CodeBlock
        code={props.code}
        language={props.language}
        showLineNumbers
        showToolbar={false}
      >
        <CodeBlockCopyButton />
      </CodeBlock>
    </div>
  )
}

export function BaseUrlBlock() {
  const baseUrl = useApiBaseUrl()
  return <CodeSample language='bash' code={`BASE_URL="${baseUrl}"`} />
}

export function FieldTable(props: { rows: FieldRow[] }) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className='w-[28%]'>字段</TableHead>
          <TableHead className='w-[16%]'>类型</TableHead>
          <TableHead className='w-[10%]'>必填</TableHead>
          <TableHead>校验与说明</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {props.rows.map((row) => (
          <TableRow key={row.name}>
            <TableCell>
              <InlineCode>{row.name}</InlineCode>
            </TableCell>
            <TableCell className='text-muted-foreground text-xs'>
              {row.type}
            </TableCell>
            <TableCell>
              <span className='text-muted-foreground text-xs'>
                {row.required}
              </span>
            </TableCell>
            <TableCell className='text-sm leading-relaxed whitespace-normal'>
              <RichText text={row.description} />
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}

export function ResponseFieldTable(props: { rows: ResponseFieldRow[] }) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className='w-[28%]'>字段</TableHead>
          <TableHead className='w-[16%]'>类型</TableHead>
          <TableHead>说明</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {props.rows.map((row) => (
          <TableRow key={row.name}>
            <TableCell>
              <InlineCode>{row.name}</InlineCode>
            </TableCell>
            <TableCell className='text-muted-foreground text-xs'>
              {row.type}
            </TableCell>
            <TableCell className='text-sm leading-relaxed whitespace-normal'>
              <RichText text={row.description} />
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}

export function StatusTable(props: { rows: StatusRow[] }) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className='w-[28%]'>状态</TableHead>
          <TableHead>说明</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {props.rows.map((row) => (
          <TableRow key={row.name}>
            <TableCell>
              <InlineCode>{row.name}</InlineCode>
            </TableCell>
            <TableCell className='text-sm leading-relaxed whitespace-normal'>
              <RichText text={row.description} />
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}

export function ErrorTable(props: { rows: ErrorRow[] }) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className='w-[40%]'>场景</TableHead>
          <TableHead>
            <code className='text-xs'>message</code> 示例
          </TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {props.rows.map((row) => (
          <TableRow key={row.scenario}>
            <TableCell className='text-sm whitespace-normal'>
              <RichText text={row.scenario} />
            </TableCell>
            <TableCell className='whitespace-normal'>
              <InlineCode>{row.message}</InlineCode>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}

export function KeyValueTable(props: {
  kLabel?: string
  vLabel?: string
  rows: KeyValueRow[]
}) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead className='w-[36%]'>{props.kLabel ?? '字段'}</TableHead>
          <TableHead>{props.vLabel ?? '说明'}</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {props.rows.map((row) => (
          <TableRow key={row.k}>
            <TableCell>
              <InlineCode>{row.k}</InlineCode>
            </TableCell>
            <TableCell className='text-sm leading-relaxed whitespace-normal'>
              <RichText text={row.v} />
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}

export function EndpointOverviewTable(props: {
  rows: { feature: string; method: string; path: string }[]
}) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>功能</TableHead>
          <TableHead className='w-[80px]'>方法</TableHead>
          <TableHead>路径</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {props.rows.map((row) => (
          <TableRow key={`${row.method} ${row.path}`}>
            <TableCell className='text-sm'>{row.feature}</TableCell>
            <TableCell>
              <InlineCode>{row.method}</InlineCode>
            </TableCell>
            <TableCell>
              <code className='text-foreground text-xs'>{row.path}</code>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}

export function EndpointCard(props: {
  anchor?: string
  number: string
  title: string
  method: string
  path: string
  contentType?: string
  remark?: string
  children: ReactNode
}) {
  return (
    <Card
      id={props.anchor}
      className='dc-endpoint bg-muted/20 scroll-mt-24 border-0 ring-0'
    >
      <CardContent className='space-y-6 pt-6'>
        <div className='space-y-3'>
          <h3 className='dc-h3 text-xl font-semibold tracking-tight'>
            <span className='dc-num'>{props.number}</span>
            {props.title}
          </h3>
          <div className='flex flex-wrap items-center gap-2'>
            <span
              data-method={props.method}
              className={cn(
                'dc-method inline-flex h-6 min-w-[58px] items-center justify-center rounded-md border px-2 font-mono text-[11px] tracking-wide',
                props.method === 'POST'
                  ? 'border-emerald-500/40 text-emerald-600 dark:text-emerald-400'
                  : 'border-sky-500/40 text-sky-600 dark:text-sky-400'
              )}
            >
              {props.method}
            </span>
            <code className='dc-path text-foreground text-sm break-all'>
              {props.path}
            </code>
          </div>
          {props.contentType ? (
            <p className='text-muted-foreground text-xs'>
              <InlineCode>Content-Type</InlineCode>: {props.contentType}
            </p>
          ) : null}
          {props.remark ? (
            <p className='text-muted-foreground text-sm leading-relaxed'>
              <RichText text={props.remark} />
            </p>
          ) : null}
        </div>
        {props.children}
      </CardContent>
    </Card>
  )
}

export function SubHeading(props: { children?: ReactNode; text?: string }) {
  return (
    <h4 className='dc-h4 text-base font-semibold tracking-tight'>
      {props.text ?? props.children}
    </h4>
  )
}

export function MinorHeading(props: { children?: ReactNode; text?: string }) {
  return (
    <h5 className='text-sm font-semibold tracking-tight'>
      {props.text ?? props.children}
    </h5>
  )
}

export function Paragraph(props: { children?: ReactNode; text?: string }) {
  const raw = props.text ?? props.children
  const content = typeof raw === 'string' ? <RichText text={raw} /> : raw
  return (
    <p className='text-muted-foreground text-sm leading-relaxed'>{content}</p>
  )
}

export function Note(props: { children?: ReactNode; text?: string }) {
  const raw = props.text ?? props.children
  const content = typeof raw === 'string' ? <RichText text={raw} /> : raw
  return (
    <blockquote className='dc-note bg-muted/40 text-muted-foreground border-primary/60 rounded-md border-l-2 px-3 py-2 text-sm leading-relaxed'>
      {content}
    </blockquote>
  )
}

export function BulletList(props: { items: string[]; ordered?: boolean }) {
  const items = props.items.map((item) => (
    <li key={item}>
      <RichText text={item} />
    </li>
  ))
  return props.ordered ? (
    <ol className='text-muted-foreground list-decimal space-y-2 pl-6 text-sm leading-relaxed'>
      {items}
    </ol>
  ) : (
    <ul className='text-muted-foreground list-disc space-y-1 pl-6 text-sm leading-relaxed'>
      {items}
    </ul>
  )
}

export function ParamDetail(props: { name: string; children: ReactNode }) {
  return (
    <div className='space-y-3'>
      <h4 className='text-base font-semibold tracking-tight'>
        <code className='text-foreground bg-muted rounded px-1.5 py-0.5 text-sm'>
          {props.name}
        </code>
      </h4>
      {props.children}
    </div>
  )
}

export function Block(props: { children: ReactNode }) {
  return <div className='space-y-3'>{props.children}</div>
}

export function RecommendationList(props: {
  items: { title: string; body: string }[]
}) {
  return (
    <div className='space-y-4'>
      {props.items.map((item) => (
        <div className='space-y-2' key={item.title}>
          <h3 className='text-base font-semibold tracking-tight'>
            {item.title}
          </h3>
          <p className='text-muted-foreground text-sm leading-relaxed'>
            <RichText text={item.body} />
          </p>
        </div>
      ))}
    </div>
  )
}
