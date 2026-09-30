# Spaceballs: State of the Art Go

Native Go/Ebitengine conversion of the supplied Amiga production, using
**Demo Construction Kit v1.0.11**

Run the complete production with its recovered music:

```sh
GOWORK=off go run .
```

Escape exits. The final page remains stationary. `-silent`, `-start`, `-ticks`
and `-capture` support native validation. Both music formats are selected by DCK.

## Conversion scope

- Recover the original disk payloads, soundtrack, artwork and animated shapes.
- Reconstruct the complete continuous sequence, preserving the reference's
  50-frame/s cadence, transitions, choreography and music synchronization.
- Use DCK for composition, bounded rendering, timelines, deformable image
  layers and automatically selected music playback.
- Keep original-specific data readers and integer animation programs in Go.
- Render native effects; the reference movie is not a playback implementation.

## Recovered source and native units

The provided production is **State of the Art**. The recording lasts 249.754
seconds at 50 frames/s, with a 1920 by 1080 picture and stereo audio. Initial
scene observations identify the close-up hands, credits, dragon illustration,
animated dancer silhouettes, several deformed patterned backdrops, trails,
outline passages and the final dancer composition.

The packed-picture reader reproduces the original CPU output. The opening
animation bank has 142 frames and the first dancer bank 335; all 318 prepared
morph contours match independently executed CPU coordinates, including unequal
vertex counts. `cmd/preview` provides native DCK views of those source assets.
Its colors and cycling are diagnostic; it is not the complete production.
The loader chain and resident program now recover 33 complete animation tables.
The hires dragon
illustration is decoded with its original plane offsets and palette. The same
preview can select individual recovered banks by their inventory names.

The first complete musical composition is also available as `-bank pattern`.
It combines the original 640-by-512 material, two independent motion views,
the 26/210-frame contour banks and the authored palette program. Its 50 Hz
background and 25 Hz silhouette clocks, all palette steps and bank cursors
match 480 independently executed original CPU states. These units are now also assembled in the complete native director.

The composed opening is available as `-bank opening-effect`. It preserves the
foreground hands, the held first display plane, the second-plane animation
and the three working buffers. Its bank switches, buffer addresses and
palette program match all 961 independently executed original VBL states.

The retained-mask trail composition is available as `-bank trails`. Its six
working masks, five display planes, authored cue sequence, zoom lookup and
32-color program match 566 independent CPU states. All 110 executed morph
contours in that passage also match the native reader. Add `-music` to an
interactive composed-effect preview to play its recovered module through DCK.

The black dancer over changing color blocks is available as `-bank blocks-effect`.
It retains the six-mask ring, the older-pose shadow and the original corner
color sequences. All 422 controller states and 162,180 generated RGB12
gradient words match the independently executed source program. The complete director also retains both animated loading handoffs.

The textured retained-dancer passage is available as `-bank noise`. It reads
the four original material planes and half-bright palette and shares mask
rendering with the trail and block effects. Its 350 controller states match
the original CPU, and a graphics test compares every shader output byte
against an independent six-plane decoder, including atlas-backed textures.

The second spirograph passage is available as `-bank wave`. It shares the
first musical scene's material renderer but uses its own motion tables,
44 pose cues, negated-color cycles and fade. All 846 recorded clock,
motion, palette and displayed-buffer states match the original program.
The converted units now use one DCK effect registry in the preview host.

The short animated tile loading passage is available as `-bank tiles`.
Its 21 original images, three-band placement, cue transport and alternating
palettes reproduce the source. The complete unpacker output and all 143
controller states, including hashes of each displayed bitmap, match the
independently executed original program.

The later material-trail and angular-pose passage is available as `-bank angular`.
Its five banks, 32 cues, bank-dependent trail age and retained material tail
match 715 independently executed controller states. Its GPU composition is
shared with the earlier textured-dancer passage.

The sliced-dancer unit is available as `-bank sliced`. It retains 71 cues,
counter rewinds, the normal/close-up lookups, per-edge endpoint deformation
and the source palette program. Its 1,618 controller states match the
independently executed original code. The parity-edge geometry uses the
same bounded DCK renderer and retained material as neighboring units.

The angular ribbon passage is available as `-bank ribbons`. It retains
four working pose pairs, forward/backward playback, mirrored coordinate
tables and the original 16-color flash/blend program. All 816 recorded
controller states match the original CPU.

The outlined body composition is available as `-bank duet`. Its two banks,
51 cues, pointer-dependent trails, palette changes and final white material
exit match all 1,041 independently executed source states.

The returning colored-block field with fine contours is available as
`-bank outline`. Its six banks, cue rewinds, mirrored coordinate lookup and
background colors match all 1,385 controller states and 530,145 RGB12
gradient words. It shares the block/material renderer with the earlier unit
while keeping the original disabled-fill and one-row outline sampling.

The final spirograph is available as `-bank finale`. Its changing shape
cadence, 94 cues, palette arithmetic, blinking display mode, visible poses
and material motion match all 1,716 recorded source states. The original
stationary artwork is also decoded and available as `-bank closing`.

```sh
GOWORK=off go run ./cmd/preview -bank intro
GOWORK=off go run ./cmd/preview -bank opening-effect
GOWORK=off go run ./cmd/preview -bank hands
GOWORK=off go run ./cmd/preview -bank first-animation
GOWORK=off go run ./cmd/preview -bank picture
GOWORK=off go run ./cmd/preview -bank dragon
GOWORK=off go run ./cmd/preview -bank final-b
GOWORK=off go run ./cmd/preview -bank pattern
GOWORK=off go run ./cmd/preview -bank trails -music
GOWORK=off go run ./cmd/preview -bank blocks-effect -music
GOWORK=off go run ./cmd/preview -bank noise -music
GOWORK=off go run ./cmd/preview -bank wave -music
GOWORK=off go run ./cmd/preview -bank tiles -music
GOWORK=off go run ./cmd/preview -bank angular -music
GOWORK=off go run ./cmd/preview -bank sliced -music
GOWORK=off go run ./cmd/preview -bank ribbons -music
GOWORK=off go run ./cmd/preview -bank duet -music
GOWORK=off go run ./cmd/preview -bank outline -music
GOWORK=off go run ./cmd/preview -bank finale -music
GOWORK=off go run ./cmd/preview -bank closing
GOWORK=off go run ./cmd/extract -disk /path/to/unpacked-source.adf
```

The continuous director now includes all recovered compositions, title/credit
pages, illustration handoffs, both loading-card sequences and the closing page.
A full update/draw traversal captured 62 segment milestones without errors.
Original-controller fixtures and the material shader pixel oracle pass. The
[director verification](reference/DIRECTOR.md) records timing sources and the
remaining differences at individual Amiga/GPU raster boundaries.

Opening, Trails, Noise, Angular, Sliced, Ribbons and Duet use DCK's shared
`composite.BitplanePalette`. The production supplies its retained poses, palette
words and channel selection; DCK owns the one/four-plane draw or five/six-plane
packing and palette passes. `composite.ContourBank` owns the six-slot and paired
working masks, selective clears and reusable geometry batches. Shared fan,
stroke and parity-edge builders preserve the original coordinate maps and
buffer-pointer program. Pattern, Wave and Finale also use the retained bank.
Their material now uses `BitplanePalette.DrawOffsets`: one texture supplies two
independently shifted planes, while the contour keeps its own fetch offset.
Fine-scroll delay, RGB12 palette and display crop remain original data; no local
material shader or additional working image is needed.

Blocks and Outline use `composite.PaletteGrid` for their background/shadow
colors and black body material. The intro supplies its 24-pixel columns,
36-pixel first row, 16-pixel later rows, skipped palette row and source words.
DCK owns the small color-bank image and lookup pass. Palette storage drops
from 408,320 to 2,040 bytes without changing the upload size or adding a pass.

Tiles and Vote use `effects.IndexedImageBank`: 21 indexed textures and two small
palettes replace 42 precolored textures. Three retained slots keep original
working/display selection and placement. The texture payload per instance drops
from 1,731,072 to 865,536 bytes, with three direct lookup draws and no conversion
surface. Every pixel of all 293 native controller ticks matches independent CPU
composition; the complete director retains all 1,218 frame samples.

The State/Of/The/Art, credit, dragon and closing pages use
`effects.IndexedImage`: DCK owns the indexed-color conversion, signed RGB12
palette transitions, absolute clock and output resources. The demo supplies
decoded artwork, original palette words, fade endpoints, crop and placement.
The complete director retains all 1,218 sampled RGBA frames over 12,414 drawn
updates, including the first 24 ticks at every segment boundary. The shared
renderer keeps the preceding conversion surface and two-pass image path.
Use `go run ./cmd/checkframes -production -output /path/to/director.json`
to reproduce the full traversal with every intervening update drawn.

The published DCK 1.0.9 Android build passed ELF/ZIP 16 KiB alignment checks
and completed the director on Pixel 10a with normal audio. Five-second samples
after startup measured 49.2–50.9 TPS and 58.0–60.1 FPS, with 66.6 MiB peak Go
heap. These runtime measurements exclude native/GPU memory and are separate
from an Android visual-reference comparison.

All 690 sampled full RGBA frames match the preceding renderer across 11,209
rendered frames in all fourteen effect units, with unchanged pass counts and
50 Hz timing. The independent six-plane decoder also checks opaque monochrome
textures placed in Ebitengine's atlas. Frame fingerprints can be reproduced with
`go run ./cmd/checkframes -output /path/to/frames.json`.

[Scene inventory](SCENES.md) · [Reference manifest](reference/sources.json)
· [Original loader](reference/LOADER.md) · [Verified animation banks](reference/ANIMATION.md)
· [First musical composition](reference/PATTERN.md)
· [Resident opening composition](reference/OPENING.md)
· [Retained-mask trail composition](reference/TRAILS.md)
· [Block-gradient composition](reference/BLOCKS.md)
· [Textured retained-dancer composition](reference/NOISE.md)
· [Cued spirograph composition](reference/WAVE.md)
· [Loading passage inventory](reference/HANDOFFS.md)
· [Animated tile composition](reference/TILES.md)
· [Authored material trails](reference/ANGULAR.md)
· [Sliced dancer composition](reference/SLICED.md)
· [Four-plane ribbon composition](reference/RIBBONS.md)
· [Outline duet composition](reference/DUET.md)
· [Returning block-field contours](reference/OUTLINE.md)
· [Final spirograph and closing image](reference/FINALE.md)
· [Continuous production director](reference/DIRECTOR.md)
