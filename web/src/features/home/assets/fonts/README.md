# Home page fonts

Subsets (pyftsubset, woff2) of SIL Open Font License 1.1 fonts, limited to the
glyphs the home page uses. Loaded only by `components/comic/home-comic.css`.

| File                   | Source font                                             | License                 |
| ---------------------- | ------------------------------------------------------- | ----------------------- |
| `kthome-display.woff2` | Smiley Sans Oblique 得意黑 (atelier-anchor/smiley-sans) | `OFL-SmileySans.txt`    |
| `zcool-kuaile.woff2`   | ZCOOL KuaiLe                                            | `OFL-ZCOOLKuaiLe.txt`   |
| `bangers.woff2`        | Bangers                                                 | `OFL-Bangers.txt`       |
| `bungee.woff2`         | Bungee                                                  | `OFL-Bungee.txt`        |
| `jetbrains-mono.woff2` | JetBrains Mono (variable)                               | `OFL-JetBrainsMono.txt` |

Smiley Sans declares the Reserved Font Names "Smiley" and "得意黑", so the
subset (a Modified Version under the OFL) is renamed to "KTHome Display".

To regenerate after changing home copy, rebuild the glyph list from the home
i18n strings in every locale and rerun pyftsubset with `--flavor=woff2`.
