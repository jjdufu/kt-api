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
import '@/styles/comic/comic-base.css'
import './home-comic.css'

export { CapabilityWall } from './capability-wall'
export { ComicCTA } from './comic-cta'
export { ComicHero } from './comic-hero'
export { ComicStrip } from './comic-strip'
export { ModelMarquee } from './model-marquee'
export { NeonStats } from './neon-stats'

// Building blocks reused by other comic-styled pages (about).
export { Burst, ComicBackdrop, Skyline, SpeechBubble } from './decorations'
export { ComicButton, Reveal, SectionHeading } from './primitives'
export { useDocsHref } from './use-docs-href'
