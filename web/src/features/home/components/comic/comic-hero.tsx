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

import heroGirl from '@/assets/comic/hero-girl.webp'

import { ComicTerminal } from './comic-terminal'
import { Burst, SpeechBubble } from './decorations'
import { splitHeadline } from './headline'
import { ComicButton } from './primitives'
import { useDocsHref } from './use-docs-href'

type ComicHeroProps = {
  isAuthenticated: boolean
}

export function ComicHero({ isAuthenticated }: ComicHeroProps) {
  const { t } = useTranslation()
  const docsHref = useDocsHref()
  // Start the hero image (LCP) right away instead of waiting for CSS/layout.
  // The comic font slices are preloaded by index.html under KT Comic.
  preload(heroGirl, { as: 'image', fetchPriority: 'high' })
  const headline = t('Define Everything with LLMs!')
  const lines = splitHeadline(headline)

  return (
    <section className='hc-hero'>
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
