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
import { useTranslation } from 'react-i18next'

import { Burst, ComicBackdrop } from './decorations'
import { ComicButton } from './primitives'
import { useDocsHref } from './use-docs-href'

/** Closing call to action; rendered for signed-out visitors only. */
export function ComicCTA() {
  const { t } = useTranslation()
  const docsHref = useDocsHref()
  return (
    <section className='hc-section'>
      <div className='hc-container'>
        <div>
          <div className='hc-cta'>
            <ComicBackdrop />
            <Burst label='POW!' className='hc-cta-pow' />
            <h2 className='hc-cta-title'>
              {t('Your story starts with one request')}
            </h2>
            <p className='hc-cta-sub'>
              {t('Sign up free and ship your first call in minutes.')}
            </p>
            <div className='hc-hero-actions hc-cta-actions'>
              <ComicButton to='/sign-up'>{t('Start Now')}</ComicButton>
              <ComicButton to={docsHref} variant='outline'>
                {t('Read the Docs')}
              </ComicButton>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
