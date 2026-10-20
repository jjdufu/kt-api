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
import { motion, useReducedMotion } from 'motion/react'
import type { ReactNode } from 'react'

import { cn } from '@/lib/utils'

type ComicButtonProps = {
  to: string
  children: ReactNode
  variant?: 'primary' | 'outline'
  className?: string
}

/** Comic button: 3px ink outline, 6px hard shadow, presses on click. */
export function ComicButton({
  to,
  children,
  variant = 'primary',
  className,
}: ComicButtonProps) {
  const classes = cn('hc-btn', `hc-btn-${variant}`, className)
  if (/^https?:\/\//.test(to)) {
    return (
      <a
        href={to}
        target='_blank'
        rel='noopener noreferrer'
        className={classes}
      >
        {children}
      </a>
    )
  }
  return (
    <Link to={to} className={classes}>
      {children}
    </Link>
  )
}

type RevealProps = {
  children: ReactNode
  className?: string
  delay?: number
  tilt?: number
}

/** Panel "slam-in" on scroll; static when reduced motion is requested. */
export function Reveal({
  children,
  className,
  delay = 0,
  tilt = -2,
}: RevealProps) {
  const reduce = useReducedMotion()
  if (reduce) {
    return <div className={className}>{children}</div>
  }
  return (
    <motion.div
      className={className}
      initial={{ opacity: 0, y: 28, scale: 0.94, rotate: tilt }}
      whileInView={{ opacity: 1, y: 0, scale: 1, rotate: 0 }}
      viewport={{ once: true, margin: '0px 0px -10% 0px' }}
      transition={{ type: 'spring', stiffness: 260, damping: 20, delay }}
    >
      {children}
    </motion.div>
  )
}

type SectionHeadingProps = {
  kicker: string
  title: string
  className?: string
}

export function SectionHeading({
  kicker,
  title,
  className,
}: SectionHeadingProps) {
  return (
    <div className={cn('hc-heading', className)}>
      <span className='hc-kicker'>{kicker}</span>
      <h2 className='hc-title'>{title}</h2>
    </div>
  )
}
