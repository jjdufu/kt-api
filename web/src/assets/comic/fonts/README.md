# KT Comic theme fonts

SIL Open Font License 1.1 fonts used only under the KT Comic theme preset.

## Full CJK fonts (unicode-range slices)

| Family           | Source font                                                    | License               |
| ---------------- | -------------------------------------------------------------- | --------------------- |
| `KTHome Display` | Smiley Sans Oblique 得意黑 v2.0.1 (atelier-anchor/smiley-sans) | `OFL-SmileySans.txt`  |
| `KTHome KuaiLe`  | ZCOOL KuaiLe Regular (google/fonts `ofl/zcoolkuaile`)          | `OFL-ZCOOLKuaiLe.txt` |

Smiley Sans declares the Reserved Font Names "Smiley" and "得意黑", so the
Modified Version is renamed to "KTHome Display" in its name table.

Both are split into woff2 slices under `public/fonts/kt-comic/` with
`unicode-range` @font-face rules in `src/styles/comic/comic-fonts-sliced.css`.
Slice 000 holds ASCII/Latin-1/punctuation plus every character in
`glyphs.txt` and `glyphs-console.txt` (home, auth and console copy) and is
preloaded by `index.html` when KT Comic is active; then GB2312 level-1 hanzi
in frequency-ish order, then the rest of the font by code point. Regenerate
(deterministic output) with:

    pip install 'fonttools[woff]' brotli
    python3 scripts/build-comic-fonts.py \
      --display SmileySans-Oblique.ttf --body ZCOOLKuaiLe-Regular.ttf

Characters missing from both fonts fall back to system fonts.

## Latin fonts (subsets)

| File                   | Source font               | License                 |
| ---------------------- | ------------------------- | ----------------------- |
| `bangers.woff2`        | Bangers                   | `OFL-Bangers.txt`       |
| `bungee.woff2`         | Bungee                    | `OFL-Bungee.txt`        |
| `jetbrains-mono.woff2` | JetBrains Mono (variable) | `OFL-JetBrainsMono.txt` |

Subset with
`pyftsubset --unicodes=U+0020-007E,U+00A0-00FF,U+2018-201D,U+2022,U+2026,U+2190-2193 --layout-features='*' --no-hinting --flavor=woff2`.
