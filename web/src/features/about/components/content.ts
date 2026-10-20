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
import { Gauge, Layers, Sparkles, type LucideIcon } from 'lucide-react'

/* Shared about-page content (i18n keys) for the comic and enterprise pages. */

export type Pillar = {
  title: string
  body: string
  icon: LucideIcon
  tone: 'cyan' | 'magenta' | 'yellow'
}

export const PILLARS: Pillar[] = [
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
