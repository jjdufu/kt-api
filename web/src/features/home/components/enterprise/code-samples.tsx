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
import { Check, Copy } from 'lucide-react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import { cn } from '@/lib/utils'

type Lang = 'curl' | 'python' | 'node'

const TABS: { id: Lang; label: string }[] = [
  { id: 'curl', label: 'cURL' },
  { id: 'python', label: 'Python' },
  { id: 'node', label: 'Node.js' },
]

function sample(lang: Lang, base: string): string {
  switch (lang) {
    case 'curl':
      return `curl ${base}/v1/chat/completions \\
  -H "Authorization: Bearer $KTAPI_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "gpt-4.1",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'`
    case 'python':
      return `import os

from openai import OpenAI

client = OpenAI(
    base_url="${base}/v1",
    api_key=os.environ["KTAPI_KEY"],
)

resp = client.chat.completions.create(
    model="claude-sonnet-4",
    messages=[{"role": "user", "content": "Hello!"}],
)
print(resp.choices[0].message.content)`
    case 'node':
      return `import OpenAI from 'openai'

const client = new OpenAI({
  baseURL: '${base}/v1',
  apiKey: process.env.KTAPI_KEY,
})

const resp = await client.chat.completions.create({
  model: 'gemini-2.5-pro',
  messages: [{ role: 'user', content: 'Hello!' }],
})
console.log(resp.choices[0].message.content)`
  }
}

/** Tabbed, copyable code sample (dark editor surface in every preset). */
export function CodeSamples() {
  const { t } = useTranslation()
  const [lang, setLang] = useState<Lang>('curl')
  const [copied, setCopied] = useState(false)
  const base =
    typeof window === 'undefined'
      ? 'https://api.example.com'
      : window.location.origin
  const code = sample(lang, base)

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(code)
      setCopied(true)
      window.setTimeout(() => setCopied(false), 1500)
    } catch {
      // Clipboard blocked: nothing to do.
    }
  }

  return (
    <div className='overflow-hidden rounded-xl border border-white/10 bg-(--code-surface) shadow-[0_24px_60px_-24px_rgb(0_0_0/0.35)]'>
      <div className='flex items-center justify-between border-b border-white/10 pe-2'>
        <div role='tablist' className='flex'>
          {TABS.map((tab) => (
            <button
              key={tab.id}
              type='button'
              role='tab'
              aria-selected={lang === tab.id}
              onClick={() => setLang(tab.id)}
              className={cn(
                '-mb-px border-b-2 px-4 py-2.5 text-xs font-medium transition-colors',
                lang === tab.id
                  ? 'border-[var(--primary)] text-white'
                  : 'border-transparent text-zinc-400 hover:text-zinc-200'
              )}
            >
              {tab.label}
            </button>
          ))}
        </div>
        <button
          type='button'
          onClick={copy}
          aria-label={t('Copy code')}
          className='inline-flex items-center gap-1.5 rounded-md px-2 py-1 text-xs text-zinc-400 transition-colors hover:bg-white/5 hover:text-zinc-100'
        >
          {copied ? (
            <Check className='size-3.5' />
          ) : (
            <Copy className='size-3.5' />
          )}
          {copied ? t('Copied') : t('Copy')}
        </button>
      </div>
      <pre className='overflow-x-auto p-5 font-mono text-[12.5px] leading-6 text-zinc-200'>
        <code>{code}</code>
      </pre>
    </div>
  )
}
