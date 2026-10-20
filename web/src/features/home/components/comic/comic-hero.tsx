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
import { preload } from 'react-dom'
import { useTranslation } from 'react-i18next'

import displayFont from '@/assets/comic/fonts/kthome-display.woff2'
import heroGirl from '@/assets/comic/hero-girl.webp'

import { ComicTerminal } from './comic-terminal'
import { Burst, ComicBackdrop, Skyline, SpeechBubble } from './decorations'
import { ComicButton } from './primitives'
import { useDocsHref } from './use-docs-href'

/**
 * Split a headline into two balanced lines. CJK text has no spaces, so break
 * after a CJK comma when there is one, otherwise in the middle.
 */
function splitHeadline(text: string): string[] {
  const trimmed = text.trim()
  if (/\s/.test(trimmed)) return [trimmed]
  const comma = trimmed.search(/[、，,]/)
  if (comma > 0 && comma < trimmed.length - 1) {
    return [trimmed.slice(0, comma + 1), trimmed.slice(comma + 1)]
  }
  if (trimmed.length < 6) return [trimmed]
  const mid = Math.floor(trimmed.length / 2)
  return [trimmed.slice(0, mid), trimmed.slice(mid)]
}

type ComicHeroProps = {
  isAuthenticated: boolean
}

export function ComicHero({ isAuthenticated }: ComicHeroProps) {
  const { t } = useTranslation()
  const docsHref = useDocsHref()
  // Start the hero image (LCP) and the CJK headline font right away instead of
  // waiting for CSS/layout to discover them. React de-duplicates the hints.
  preload(heroGirl, { as: 'image', fetchPriority: 'high' })
  preload(displayFont, { as: 'font', type: 'font/woff2', crossOrigin: '' })
  const headline = t('Define Everything with LLMs!')
  const lines = splitHeadline(headline)

  return (
    <section className='hc-hero'>
      <ComicBackdrop />
      <Skyline className='hc-hero-skyline' />
      <div className='hc-hero-inner'>
        <div className='hc-hero-copy'>
          <h1 className='hc-headline' aria-label={headline}>
            {lines.map((line) => (
              <span
                key={line}
                className='hc-glitch'
                data-text={line}
                aria-hidden='true'
              >
                {line}
              </span>
            ))}
          </h1>
          <SpeechBubble className='hc-hero-bubble'>
            {t('A friend to every creator')}
          </SpeechBubble>
          <div className='hc-hero-actions'>
            <ComicButton to={isAuthenticated ? '/dashboard' : '/sign-up'}>
              {t('Start Now')}
            </ComicButton>
            <ComicButton to={docsHref} variant='outline'>
              {t('Read the Docs')}
            </ComicButton>
            <Burst label='ZAP!' className='hc-hero-zap' />
          </div>
        </div>

        <div className='hc-hero-stage'>
          <div className='hc-panel-frame'>
            <div className='hc-panel'>
              <div className='hc-panel-lines' aria-hidden='true' />
              <ComicTerminal />
            </div>
          </div>
          <img
            src={heroGirl}
            alt=''
            aria-hidden='true'
            width={640}
            height={887}
            className='hc-hero-girl'
            fetchPriority='high'
            decoding='async'
          />
          <SpeechBubble className='hc-panel-bubble' tail='right'>
            {t('One API, every model')}
          </SpeechBubble>
        </div>
      </div>
    </section>
  )
}
