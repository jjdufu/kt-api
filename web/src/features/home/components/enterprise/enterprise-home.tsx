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
import { ArrowRight, Check, ChevronRight } from 'lucide-react'
import type { ReactNode } from 'react'
import { useTranslation } from 'react-i18next'

import { BrandLogo } from '@/components/brand-logo'
import { PublicLayout } from '@/components/layout'
import {
  Container,
  CtaLink,
  HeroBackdrop,
  Kicker,
  Section,
  SectionHeading,
  Surface,
} from '@/features/marketing/primitives'
import { useDocsHref } from '@/features/marketing/use-docs-href'
import { useStatus } from '@/hooks/use-status'
import { useSystemConfig } from '@/hooks/use-system-config'
import { getLobeIcon } from '@/lib/lobe-icon'
import { cn } from '@/lib/utils'
import { useAuthStore } from '@/stores/auth-store'

import { AI_MODELS, getDefaultStats } from '../../constants'
import { CAPABILITIES, STEPS } from '../comic/content'
import { splitHeadline } from '../comic/headline'
import { CodeSamples } from './code-samples'
import { ProductMock } from './product-mock'

/** Page shell shared by the enterprise home and about pages. */
export function EnterpriseShell(props: { children: ReactNode }) {
  return (
    <div data-enterprise='' className='bg-background text-foreground'>
      <PublicLayout showMainContainer={false}>
        <main className='isolate overflow-x-clip'>{props.children}</main>
        <EnterpriseFooter />
      </PublicLayout>
    </div>
  )
}

/* ── Hero ─────────────────────────────────────────────────────────────── */
function Hero(props: { isAuthenticated: boolean }) {
  const { t } = useTranslation()
  const docsHref = useDocsHref()
  return (
    <section className='relative isolate'>
      <HeroBackdrop />
      <Container className='grid items-center gap-16 pt-16 pb-24 sm:pt-24 lg:grid-cols-[minmax(0,1fr)_minmax(0,1.15fr)] lg:gap-12 lg:pt-28 lg:pb-32'>
        <div>
          <span className='bg-background/70 text-muted-foreground inline-flex items-center gap-2 rounded-full border px-3 py-1 text-xs font-medium backdrop-blur'>
            <span className='bg-primary size-1.5 rounded-full' />
            {t('A friend to every creator')}
          </span>
          <h1 className='text-foreground mt-6 text-4xl leading-[1.08] font-semibold tracking-tight text-balance sm:text-5xl lg:text-6xl'>
            {splitHeadline(t('Define Everything with LLMs!')).map((line) => (
              <span key={line} className='block'>
                {line}
              </span>
            ))}
          </h1>
          <p className='text-muted-foreground mt-6 max-w-xl text-lg leading-relaxed text-pretty'>
            {t(
              'KTAPI is a large-model gateway for builders and creators: one key, one endpoint, every model you need.'
            )}
          </p>
          <div className='mt-9 flex flex-wrap gap-3'>
            <CtaLink to={props.isAuthenticated ? '/dashboard' : '/sign-up'}>
              {t('Start Now')}
              <ArrowRight className='size-4' />
            </CtaLink>
            <CtaLink to={docsHref} variant='outline'>
              {t('Read the Docs')}
            </CtaLink>
          </div>
        </div>
        <div className='lg:ps-6'>
          <ProductMock />
        </div>
      </Container>
    </section>
  )
}

/* ── Logo strip ───────────────────────────────────────────────────────── */
const PROVIDERS = AI_MODELS.map((name) => name.replace(/\.Color$/, ''))

function LogoStrip() {
  const { t } = useTranslation()
  return (
    <section className='border-y'>
      <Container className='py-10'>
        <p className='text-muted-foreground text-center text-sm'>
          {t('Plugged into the models you love')}
        </p>
        <ul className='mt-7 flex flex-wrap items-center justify-center gap-x-10 gap-y-6 sm:gap-x-14'>
          {PROVIDERS.map((name) => (
            <li
              key={name}
              className='text-muted-foreground hover:text-foreground flex items-center gap-2.5 opacity-70 grayscale transition hover:opacity-100'
            >
              {getLobeIcon(name, 24)}
              <span className='text-sm font-medium'>{name}</span>
            </li>
          ))}
        </ul>
      </Container>
    </section>
  )
}

/* ── Bento features ───────────────────────────────────────────────────── */
function PreviewEndpoints() {
  return (
    <Surface className='divide-y font-mono text-xs'>
      {[
        ['POST', '/v1/chat/completions'],
        ['POST', '/v1/images/generations'],
        ['POST', '/v1/videos'],
        ['POST', '/v1/messages'],
      ].map(([m, p]) => (
        <div key={p} className='flex items-center gap-3 px-3 py-2'>
          <span className='bg-primary/10 text-primary rounded px-1.5 py-0.5 text-[10px] font-semibold'>
            {m}
          </span>
          <span className='text-foreground truncate'>{p}</span>
        </div>
      ))}
    </Surface>
  )
}

function PreviewRouting() {
  const targets = [
    { name: 'gpt-4.1', w: '62%' },
    { name: 'claude-sonnet-4', w: '28%' },
    { name: 'gemini-2.5-pro', w: '10%' },
  ]
  return (
    <Surface className='space-y-3 p-3 font-mono text-xs'>
      <div className='text-muted-foreground flex items-center gap-2'>
        <span className='bg-foreground/80 size-1.5 rounded-full' />
        /v1/chat/completions
      </div>
      {targets.map((x) => (
        <div key={x.name} className='space-y-1'>
          <div className='flex justify-between'>
            <span className='text-foreground'>{x.name}</span>
            <span className='text-muted-foreground tabular-nums'>{x.w}</span>
          </div>
          <div className='bg-muted h-1.5 overflow-hidden rounded-full'>
            <div
              className='bg-primary h-full rounded-full'
              style={{ width: x.w }}
            />
          </div>
        </div>
      ))}
    </Surface>
  )
}

const BILLING_BARS = [30, 46, 38, 64, 52, 70, 58].map((h, i) => ({
  id: `d${i}`,
  h,
}))

function PreviewBilling() {
  return (
    <Surface className='flex h-20 items-end gap-1.5 p-3'>
      {BILLING_BARS.map(({ id, h }) => (
        <div
          key={id}
          className='bg-primary/70 flex-1 rounded-sm'
          style={{ height: `${h}%` }}
        />
      ))}
    </Surface>
  )
}

function PreviewKeys() {
  return (
    <Surface className='divide-y font-mono text-xs'>
      {['sk-prod-••••4f2a', 'sk-dev-••••91c0'].map((k, i) => (
        <div key={k} className='flex items-center justify-between px-3 py-2'>
          <span className='text-foreground'>{k}</span>
          <span
            className={cn(
              'size-1.5 rounded-full',
              i === 0 ? 'bg-success' : 'bg-muted-foreground/40'
            )}
          />
        </div>
      ))}
    </Surface>
  )
}

function PreviewTasks() {
  return (
    <Surface className='space-y-2 p-3 font-mono text-xs'>
      {[
        { id: 'task_8f21', p: 100 },
        { id: 'task_8f22', p: 64 },
      ].map((x) => (
        <div key={x.id} className='flex items-center gap-3'>
          <span className='text-foreground w-20 shrink-0'>{x.id}</span>
          <div className='bg-muted h-1.5 flex-1 overflow-hidden rounded-full'>
            <div
              className={cn(
                'h-full rounded-full',
                x.p === 100 ? 'bg-success' : 'bg-primary'
              )}
              style={{ width: `${x.p}%` }}
            />
          </div>
        </div>
      ))}
    </Surface>
  )
}

function PreviewProtocols() {
  return (
    <div className='flex flex-wrap gap-2'>
      {[
        'OpenAI',
        'Claude',
        'Gemini',
        'Responses',
        'Realtime',
        'Embeddings',
      ].map((p) => (
        <span
          key={p}
          className='bg-background/60 text-foreground rounded-md border px-2.5 py-1 font-mono text-xs'
        >
          {p}
        </span>
      ))}
    </div>
  )
}

const PREVIEWS = [
  PreviewEndpoints,
  PreviewRouting,
  PreviewBilling,
  PreviewKeys,
  PreviewTasks,
  PreviewProtocols,
]
// Bento placement on lg (3 columns): wide, tall, then four smaller tiles
// with the last one wide.
const SPANS = ['lg:col-span-2', 'lg:row-span-2', '', '', '', 'lg:col-span-2']

function Features() {
  const { t } = useTranslation()
  return (
    <Section>
      <SectionHeading
        kicker={t('Core Capabilities')}
        title={t(
          'One OpenAI-compatible endpoint for chat, images, video and more.'
        )}
      />
      <div className='grid gap-4 sm:grid-cols-2 lg:grid-cols-3'>
        {CAPABILITIES.map((c, i) => {
          const Icon = c.icon
          const Preview = PREVIEWS[i]
          return (
            <article
              key={c.no}
              className={cn(
                'bg-card group flex flex-col gap-6 rounded-xl border p-6 transition-colors hover:border-[color-mix(in_oklch,var(--primary)_35%,var(--border))]',
                SPANS[i]
              )}
            >
              <div>
                <div className='text-primary bg-primary/10 flex size-9 items-center justify-center rounded-lg'>
                  <Icon className='size-4.5' />
                </div>
                <h3 className='text-foreground mt-4 text-base font-semibold'>
                  {t(c.title)}
                </h3>
                <p className='text-muted-foreground mt-1.5 text-sm leading-relaxed'>
                  {t(c.description)}
                </p>
              </div>
              <div className='mt-auto' aria-hidden='true'>
                <Preview />
              </div>
            </article>
          )
        })}
      </div>
    </Section>
  )
}

/* ── Steps ────────────────────────────────────────────────────────────── */
function Steps() {
  const { t } = useTranslation()
  return (
    <Section className='bg-muted/30 border-y'>
      <SectionHeading
        kicker={t('Get Started')}
        title={t('Make your first call in three steps')}
      />
      <ol className='relative grid gap-10 md:grid-cols-3 md:gap-8'>
        <div
          aria-hidden='true'
          className='bg-border absolute top-5 right-[16%] left-[16%] hidden h-px md:block'
        />
        {STEPS.map((s, i) => {
          const Icon = s.icon
          return (
            <li key={s.title} className='relative text-center'>
              <div className='bg-background text-foreground relative mx-auto flex size-10 items-center justify-center rounded-full border font-mono text-sm font-semibold shadow-xs'>
                {String(i + 1).padStart(2, '0')}
              </div>
              <h3 className='text-foreground mt-5 inline-flex items-center gap-2 text-base font-semibold'>
                <Icon className='text-primary size-4' />
                {s.to ? (
                  <Link to={s.to} className='hover:underline'>
                    {t(s.title)}
                  </Link>
                ) : (
                  t(s.title)
                )}
              </h3>
              <p className='text-muted-foreground mx-auto mt-2 max-w-xs text-sm leading-relaxed'>
                {t(s.description)}
              </p>
            </li>
          )
        })}
      </ol>
    </Section>
  )
}

/* ── Stats ────────────────────────────────────────────────────────────── */
function Stats() {
  const { t } = useTranslation()
  const stats = getDefaultStats(t)
  return (
    <Section className='py-16 sm:py-20'>
      <Kicker className='mb-8 text-center'>{t('By the Numbers')}</Kicker>
      <dl className='grid grid-cols-2 divide-y rounded-xl border lg:grid-cols-4 lg:divide-x lg:divide-y-0 [&>*:nth-child(odd)]:border-r lg:[&>*:nth-child(odd)]:border-r-0'>
        {stats.map((s) => (
          <div key={s.description} className='px-6 py-8 text-center'>
            <dd className='text-foreground text-4xl font-semibold tracking-tight tabular-nums'>
              {s.value}
              <span className='text-primary'>{s.suffix}</span>
            </dd>
            <dt className='text-muted-foreground mt-2 text-sm'>
              {s.description}
            </dt>
          </div>
        ))}
      </dl>
    </Section>
  )
}

/* ── Code ─────────────────────────────────────────────────────────────── */
function Code() {
  const { t } = useTranslation()
  const points = CAPABILITIES.filter((c) =>
    ['#01', '#02', '#06'].includes(c.no)
  )
  return (
    <Section>
      <div className='grid items-center gap-12 lg:grid-cols-[minmax(0,0.9fr)_minmax(0,1.1fr)]'>
        <div>
          <SectionHeading
            align='left'
            className='mb-8'
            kicker={t('OpenAI Compatible')}
            title={t('One API, every model')}
          />
          <ul className='space-y-5'>
            {points.map((c) => (
              <li key={c.no} className='flex gap-3'>
                <span className='bg-primary/10 text-primary mt-0.5 flex size-5 shrink-0 items-center justify-center rounded-full'>
                  <Check className='size-3' />
                </span>
                <div>
                  <p className='text-foreground text-sm font-semibold'>
                    {t(c.title)}
                  </p>
                  <p className='text-muted-foreground mt-1 text-sm leading-relaxed'>
                    {t(c.description)}
                  </p>
                </div>
              </li>
            ))}
          </ul>
        </div>
        <CodeSamples />
      </div>
    </Section>
  )
}

/* ── CTA ──────────────────────────────────────────────────────────────── */
export function QuietCta() {
  const { t } = useTranslation()
  const docsHref = useDocsHref()
  return (
    <Section className='pt-4'>
      <div className='relative isolate overflow-hidden rounded-2xl border px-6 py-14 text-center sm:px-12'>
        <div
          aria-hidden='true'
          className='absolute inset-0 -z-10 bg-[radial-gradient(60%_120%_at_50%_0%,color-mix(in_oklch,var(--primary)_12%,transparent),transparent)]'
        />
        <h2 className='text-foreground text-3xl font-semibold tracking-tight text-balance sm:text-4xl'>
          {t('Start building with one API')}
        </h2>
        <p className='text-muted-foreground mx-auto mt-4 max-w-xl text-base'>
          {t('Sign up free and ship your first call in minutes.')}
        </p>
        <div className='mt-8 flex flex-wrap justify-center gap-3'>
          <CtaLink to='/sign-up'>
            {t('Start Now')}
            <ArrowRight className='size-4' />
          </CtaLink>
          <CtaLink to={docsHref} variant='outline'>
            {t('Read the Docs')}
          </CtaLink>
        </div>
      </div>
    </Section>
  )
}

/* ── Footer ───────────────────────────────────────────────────────────── */
function FooterLink(props: { to: string; children: ReactNode }) {
  const className =
    'text-muted-foreground hover:text-foreground inline-flex items-center gap-1 text-sm transition-colors'
  if (/^https?:\/\//.test(props.to)) {
    return (
      <a
        href={props.to}
        target='_blank'
        rel='noopener noreferrer'
        className={className}
      >
        {props.children}
      </a>
    )
  }
  return (
    <Link to={props.to} className={className}>
      {props.children}
    </Link>
  )
}

export function EnterpriseFooter() {
  const { t } = useTranslation()
  const { status } = useStatus()
  const { systemName, logo } = useSystemConfig()
  const docsHref = useDocsHref()
  const isAuthenticated = !!useAuthStore((s) => s.auth.user)
  const name = systemName || 'KTAPI'

  const columns: { title: string; links: { label: string; to: string }[] }[] = [
    {
      title: t('Product'),
      links: [
        { label: t('Home'), to: '/' },
        ...(isAuthenticated
          ? [
              { label: t('Console'), to: '/dashboard' },
              { label: t('Model Square'), to: '/pricing' },
            ]
          : [{ label: t('Sign In'), to: '/sign-in' }]),
      ],
    },
    {
      title: t('Resources'),
      links: [
        { label: t('Docs'), to: docsHref },
        { label: t('API Keys'), to: isAuthenticated ? '/keys' : '/sign-in' },
        { label: t('About'), to: '/about' },
      ],
    },
    {
      title: t('Legal'),
      links: [
        ...(status?.user_agreement_enabled
          ? [{ label: t('User Agreement'), to: '/user-agreement' }]
          : []),
        ...(status?.privacy_policy_enabled
          ? [{ label: t('Privacy Policy'), to: '/privacy-policy' }]
          : []),
        {
          label: t('Based on New API (AGPL-3.0)'),
          to: 'https://github.com/QuantumNous/new-api',
        },
      ],
    },
  ]

  return (
    <footer className='border-t'>
      <Container className='grid gap-10 py-14 sm:grid-cols-2 lg:grid-cols-4'>
        <div>
          <Link to='/' className='inline-flex items-center gap-2.5'>
            <BrandLogo src={logo} alt='' className='size-7 rounded-md' />
            <span className='text-foreground font-semibold'>{name}</span>
          </Link>
          <p className='text-muted-foreground mt-3 max-w-xs text-sm'>
            {t('One API, every model')}
          </p>
        </div>
        {columns.map((col) => (
          <div key={col.title}>
            <h3 className='text-foreground text-sm font-semibold'>
              {col.title}
            </h3>
            <ul className='mt-4 space-y-2.5'>
              {col.links.map((l) => (
                <li key={l.label}>
                  <FooterLink to={l.to}>{l.label}</FooterLink>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </Container>
      <div className='border-t'>
        <Container className='text-muted-foreground flex flex-wrap items-center justify-between gap-2 py-6 text-xs'>
          <span>
            © {new Date().getFullYear()} {name}
          </span>
          <FooterLink to={docsHref}>
            {t('Docs')}
            <ChevronRight className='size-3' />
          </FooterLink>
        </Container>
      </div>
    </footer>
  )
}

/** Home page for every preset except KT Comic. */
export function EnterpriseHome() {
  const isAuthenticated = !!useAuthStore((s) => s.auth.user)
  return (
    <EnterpriseShell>
      <Hero isAuthenticated={isAuthenticated} />
      <LogoStrip />
      <Features />
      <Steps />
      <Stats />
      <Code />
      {!isAuthenticated && <QuietCta />}
    </EnterpriseShell>
  )
}
