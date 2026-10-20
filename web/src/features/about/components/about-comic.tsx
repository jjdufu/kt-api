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
  BookOpenText,
  Gauge,
  KeyRound,
  Layers,
  Sparkles,
  type LucideIcon,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { PublicLayout } from '@/components/layout'
import { Footer } from '@/components/layout/components/footer'
import {
  Burst,
  CapabilityWall,
  ComicBackdrop,
  ComicButton,
  ModelMarquee,
  Reveal,
  SectionHeading,
  Skyline,
  SpeechBubble,
  useDocsHref,
} from '@/features/home/components/comic'
import { useAuthStore } from '@/stores/auth-store'

import '../about-comic.css'

type Pillar = {
  title: string
  body: string
  icon: LucideIcon
  tone: 'cyan' | 'magenta' | 'yellow'
}

const PILLARS: Pillar[] = [
  {
    title: 'One gateway',
    body: 'KTAPI puts every major model behind one OpenAI-compatible endpoint, so switching models is a one-line change.',
    icon: Layers,
    tone: 'cyan',
  },
  {
    title: 'Built for creators',
    body: 'Text, images, video and assets share the same keys, billing and logs, so you spend time creating instead of wiring vendors together.',
    icon: Sparkles,
    tone: 'magenta',
  },
  {
    title: 'Clear and in control',
    body: 'Usage-based billing, per-key quotas and detailed logs show exactly where every request and every credit went.',
    icon: Gauge,
    tone: 'yellow',
  },
]

type Entry = { title: string; body: string; to: string; icon: LucideIcon }

function EntryCard({ entry }: { entry: Entry }) {
  const { t } = useTranslation()
  const Icon = entry.icon
  const body = (
    <>
      <span className='ab-entry-icon'>
        <Icon strokeWidth={2.4} />
      </span>
      <span className='ab-entry-title'>{t(entry.title)}</span>
      <span className='ab-entry-body'>{t(entry.body)}</span>
    </>
  )
  if (/^https?:\/\//.test(entry.to)) {
    return (
      <a
        href={entry.to}
        target='_blank'
        rel='noopener noreferrer'
        className='ab-entry'
      >
        {body}
      </a>
    )
  }
  return (
    <Link to={entry.to} className='ab-entry'>
      {body}
    </Link>
  )
}

/** Built-in about page, shown when the admin has not set About content. */
export function AboutComic() {
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
    <div data-home=''>
      <PublicLayout showMainContainer={false}>
        <main className='hc-main'>
          <section className='hc-hero ab-hero'>
            <ComicBackdrop />
            <Skyline className='ab-skyline' />
            <div className='hc-container ab-hero-inner'>
              <span className='hc-kicker'>{t('About')}</span>
              <h1 className='hc-headline ab-headline'>
                <span className='hc-glitch' data-text={t('About KTAPI')}>
                  {t('About KTAPI')}
                </span>
              </h1>
              <SpeechBubble className='hc-hero-bubble'>
                {t('One API, every model')}
              </SpeechBubble>
              <p className='ab-lead'>
                {t(
                  'KTAPI is a large-model gateway for builders and creators: one key, one endpoint, every model you need.'
                )}
              </p>
              <Burst label='HELLO!' className='ab-hello' />
            </div>
          </section>

          <section className='hc-section'>
            <div className='hc-container'>
              <SectionHeading
                kicker={t('What is KTAPI')}
                title={t('Open the sky with one request')}
              />
              <div className='ab-pillars'>
                {PILLARS.map((p, i) => {
                  const Icon = p.icon
                  return (
                    <Reveal
                      key={p.title}
                      delay={i * 0.1}
                      className='ab-pillar-cell'
                    >
                      <article
                        className={`hc-card hc-tone-${p.tone} ab-pillar`}
                      >
                        <span className='hc-card-no'>{`#0${i + 1}`}</span>
                        <span className='hc-card-icon'>
                          <Icon strokeWidth={2.2} />
                        </span>
                        <h3 className='hc-card-title'>{t(p.title)}</h3>
                        <p className='hc-card-desc'>{t(p.body)}</p>
                      </article>
                    </Reveal>
                  )
                })}
              </div>
            </div>
          </section>

          <CapabilityWall />
          <ModelMarquee />

          <section className='hc-section'>
            <div className='hc-container'>
              <SectionHeading
                kicker={t('Jump In')}
                title={t('Where to go next')}
              />
              <div className='ab-entries'>
                {entries.map((entry) => (
                  <EntryCard key={entry.title} entry={entry} />
                ))}
              </div>
              {!isAuthenticated && (
                <div className='hc-hero-actions ab-actions'>
                  <ComicButton to='/sign-up'>{t('Start Now')}</ComicButton>
                  <ComicButton to={docsHref} variant='outline'>
                    {t('Read the Docs')}
                  </ComicButton>
                </div>
              )}
              <p className='ab-credit'>
                {t(
                  'KTAPI is built on QuantumNous/new-api and licensed under AGPL-3.0.'
                )}
              </p>
            </div>
          </section>
        </main>
        <Footer />
      </PublicLayout>
    </div>
  )
}
