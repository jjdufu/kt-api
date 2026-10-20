# Home page fonts

Subsets (pyftsubset, woff2) of SIL Open Font License 1.1 fonts, limited to the
glyphs the comic pages use. Declared in `src/styles/comic/comic-base.css`,
which only the home, about, docs and auth pages import; the display fonts are
used for headings, bubbles and buttons only (body text keeps the app font).

| File                   | Source font                                             | License                 |
| ---------------------- | ------------------------------------------------------- | ----------------------- |
| `kthome-display.woff2` | Smiley Sans Oblique 得意黑 (atelier-anchor/smiley-sans) | `OFL-SmileySans.txt`    |
| `zcool-kuaile.woff2`   | ZCOOL KuaiLe                                            | `OFL-ZCOOLKuaiLe.txt`   |
| `bangers.woff2`        | Bangers                                                 | `OFL-Bangers.txt`       |
| `bungee.woff2`         | Bungee                                                  | `OFL-Bungee.txt`        |
| `jetbrains-mono.woff2` | JetBrains Mono (variable)                               | `OFL-JetBrainsMono.txt` |

Smiley Sans declares the Reserved Font Names "Smiley" and "得意黑", so the
subset (a Modified Version under the OFL) is renamed to "KTHome Display".

`glyphs.txt` is the CJK glyph set (home/about copy, docs headings and auth
headings in every locale). After changing that copy, add the new characters
and rerun, for example:

    pyftsubset SmileySans-Oblique.ttf --text-file=glyphs.txt \
      --layout-features='*' --no-hinting --desubroutinize --flavor=woff2 \
      --output-file=kthome-display.woff2   # then rename the name table
    pyftsubset ZCOOLKuaiLe-Regular.ttf --text-file=glyphs.txt ... \
      --output-file=zcool-kuaile.woff2

Latin fonts (Bangers, Bungee, JetBrains Mono) use
`--unicodes=U+0020-007E,U+00A0-00FF,U+2018-201D,U+2022,U+2026,U+2190-2193`.
Missing glyphs (e.g. kana, Cyrillic) fall back to system fonts.
