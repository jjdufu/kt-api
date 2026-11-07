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
import { fireEvent, render, screen } from '@testing-library/react'
import { afterEach, describe, expect, it, vi } from 'vitest'

import { ChangeBadge } from '@/components/change-badge'

afterEach(() => {
  document.body.removeAttribute('data-theme-preset')
  document.body.removeAttribute('data-theme-radius')
  document.documentElement.classList.remove('dark')
})

describe('ChangeBadge theme integration', () => {
  it.each([
    ['default', 'light'],
    ['default', 'dark'],
    ['kt-comic', 'light'],
    ['kt-comic', 'dark'],
  ])('uses semantic tokens for %s in %s mode', (preset, mode) => {
    document.body.dataset.themePreset = preset
    document.body.dataset.themeRadius = 'none'
    document.documentElement.classList.toggle('dark', mode === 'dark')
    const onClick = vi.fn()
    render(<ChangeBadge label='Pending deletion' showDot onClick={onClick} />)
    const badge = screen.getByTitle('Pending deletion')
    expect(badge).toHaveClass(
      'bg-warning/10',
      'text-foreground',
      'rounded-md',
      'text-xs'
    )
    expect(badge.querySelector('[aria-hidden]')).toHaveClass('bg-warning')
    expect(badge.className).not.toMatch(/amber|\[11px\]|cursor-copy/)
    expect(document.body).toHaveAttribute('data-theme-radius', 'none')
    fireEvent.click(badge)
    expect(onClick).toHaveBeenCalledOnce()
  })

  it('shares the warning treatment while accepting layout overrides', () => {
    render(<ChangeBadge label='Modified' className='max-w-40' />)
    expect(screen.getByTitle('Modified')).toHaveClass(
      'max-w-40',
      'ring-warning/25'
    )
  })
})
