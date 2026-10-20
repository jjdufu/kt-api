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
  ArrowUpRight,
  BookOpenText,
  Gauge,
  KeyRound,
  Scale,
  type LucideIcon,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { useDocsHref } from '@/features/home/components/comic/use-docs-href'
import { EnterpriseShell } from '@/features/home/components/enterprise/enterprise-home'
import {
  Container,
  CtaLink,
  HeroBackdrop,
  Kicker,
  Section,
  SectionHeading,
} from '@/features/home/components/enterprise/primitives'
import { useAuthStore } from '@/stores/auth-store'

import { PILLARS } from './content'

type Entry = { title: string; body: string; to: string; icon: LucideIcon }

function EntryCard({ entry }: { entry: Entry }) {
  const { t } = useTranslation()
  const Icon = entry.icon
  const className =
    'bg-card group flex flex-col rounded-xl border p-6 transition-colors hover:border-[color-mix(in_oklch,var(--primary)_35%,var(--border))]'
  const body = (
    <>
      <div className='flex items-center justify-between'>
        <span className='bg-primary/10 text-primary flex size-9 items-center justify-center rounded-lg'>
          <Icon className='size-4.5' />
        </span>
        <ArrowUpRight className='text-muted-foreground group-hover:text-foreground size-4 transition-colors' />
      </div>
      <span className='text-foreground mt-5 text-base font-semibold'>
        {t(entry.title)}
      </span>
      <span className='text-muted-foreground mt-1.5 text-sm'>
        {t(entry.body)}
      </span>
    </>
  )
  if (/^https?:\/\//.test(entry.to)) {
    return (
      <a
        href={entry.to}
        target='_blank'
        rel='noopener noreferrer'
        className={className}
      >
        {body}
      </a>
    )
  }
  return (
    <Link to={entry.to} className={className}>
      {body}
    </Link>
  )
}

/** About page for every preset except KT Comic. */
export function AboutEnterprise() {
  const { t } = useTranslation()
  const docsHref = useDocsHref()
  const isAuthenticated = !!useAuthStore((s) => s.auth.user)
  const entries: Entry[] = [
    {
      title: 'Console',
      body: 'Usage, balance and logs at a glance.',
      to: isAuthenticated ? '/dashboard' : '/sign-in',
      icon: Gauge,
    },
    {
      title: 'API Keys',
      body: 'Create and manage keys for your apps.',
      to: isAuthenticated ? '/keys' : '/sign-in',
      icon: KeyRound,
    },
    {
      title: 'Docs',
      body: 'Endpoints, parameters and examples.',
      to: docsHref,
      icon: BookOpenText,
    },
  ]

  return (
    <EnterpriseShell>
      {/* Hero: left-aligned */}
      <section className='relative isolate border-b'>
        <HeroBackdrop />
        <Container className='pt-20 pb-20 sm:pt-28 sm:pb-24'>
          <Kicker>{t('About')}</Kicker>
          <h1 className='text-foreground mt-4 max-w-3xl text-4xl leading-[1.08] font-semibold tracking-tight text-balance sm:text-5xl lg:text-6xl'>
            {t('About KTAPI')}
          </h1>
          <p className='text-muted-foreground mt-6 max-w-2xl text-lg leading-relaxed text-pretty'>
            {t(
              'KTAPI is a large-model gateway for builders and creators: one key, one endpoint, every model you need.'
            )}
          </p>
        </Container>
      </section>

      {/* Mission: two columns */}
      <Section>
        <div className='grid gap-10 lg:grid-cols-[minmax(0,0.9fr)_minmax(0,1.1fr)] lg:gap-16'>
          <div>
            <Kicker>{t('What is KTAPI')}</Kicker>
            <h2 className='text-foreground mt-3 text-3xl font-semibold tracking-tight text-balance sm:text-4xl'>
              {t('Open the sky with one request')}
            </h2>
          </div>
          <div className='space-y-6'>
            <p className='text-foreground border-primary border-s-2 ps-5 text-xl leading-relaxed font-medium sm:text-2xl'>
              {t('One API, every model')}
            </p>
            <p className='text-muted-foreground text-base leading-relaxed'>
              {t(PILLARS[0].body)}
            </p>
          </div>
        </div>
      </Section>

      {/* Pillars */}
      <Section className='bg-muted/30 border-y'>
        <div className='grid gap-4 md:grid-cols-3'>
          {PILLARS.map((p, i) => {
            const Icon = p.icon
            return (
              <article key={p.title} className='bg-card rounded-xl border p-6'>
                <div className='flex items-center justify-between'>
                  <span className='bg-primary/10 text-primary flex size-9 items-center justify-center rounded-lg'>
                    <Icon className='size-4.5' />
                  </span>
                  <span className='text-muted-foreground font-mono text-xs'>
                    {String(i + 1).padStart(2, '0')}
                  </span>
                </div>
                <h3 className='text-foreground mt-5 text-base font-semibold'>
                  {t(p.title)}
                </h3>
                <p className='text-muted-foreground mt-2 text-sm leading-relaxed'>
                  {t(p.body)}
                </p>
              </article>
            )
          })}
        </div>
      </Section>

      {/* Next steps */}
      <Section>
        <SectionHeading kicker={t('Jump In')} title={t('Where to go next')} />
        <div className='grid gap-4 md:grid-cols-3'>
          {entries.map((e) => (
            <EntryCard key={e.title} entry={e} />
          ))}
        </div>
        {!isAuthenticated && (
          <div className='mt-10 flex flex-wrap justify-center gap-3'>
            <CtaLink to='/sign-up'>{t('Start Now')}</CtaLink>
            <CtaLink to={docsHref} variant='outline'>
              {t('Read the Docs')}
            </CtaLink>
          </div>
        )}
        <div className='bg-muted/30 text-muted-foreground mx-auto mt-14 flex max-w-2xl items-start gap-3 rounded-lg border px-4 py-3 text-sm'>
          <Scale className='mt-0.5 size-4 shrink-0' />
          <span>
            {t(
              'KTAPI is built on QuantumNous/new-api and licensed under AGPL-3.0.'
            )}
          </span>
        </div>
      </Section>
    </EnterpriseShell>
  )
}
