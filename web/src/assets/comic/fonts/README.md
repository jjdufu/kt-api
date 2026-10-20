# KT Comic theme fonts

KT Comic uses one font set on every page (console, home, about, docs, auth):

| Role                   | Font                                                         | License                                 |
| ---------------------- | ------------------------------------------------------------ | --------------------------------------- |
| Headings, nav, buttons | Chakra Petch (Latin subset, `src/assets/fonts/chakra-petch`) | `src/assets/fonts/chakra-petch/OFL.txt` |
| Body text and all CJK  | System sans (PingFang SC, Microsoft YaHei, Noto Sans CJK…)   | —                                       |
| Code and numbers       | `jetbrains-mono.woff2`, JetBrains Mono (variable)            | `OFL-JetBrainsMono.txt`                 |

JetBrains Mono was subset with
`pyftsubset --unicodes=U+0020-007E,U+00A0-00FF,U+2018-201D,U+2022,U+2026,U+2190-2193 --layout-features='*' --no-hinting --flavor=woff2`.
