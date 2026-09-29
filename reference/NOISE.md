# Textured, half-bright retained dancers

The side-track 108 controller at `$40000` reuses the six-mask ring with five
displayed pose planes. Its sixth display plane selects one of four original
352-by-290 monochrome materials from the payload loaded at `$6cb10`.
The sixth bit selects the half-bright version of a 32-color RGB12 palette.
This follows BPLCON0 `$6200` and the
[Commodore playfield register definition](https://www.theflatnet.de/pub/cbm/amiga/AmigaDevDocs/hard_3.html).

The material transport advances an eight-position counter every software
callback and selects `counter / 2`, giving two callbacks per material.
Startup advances that transport once, then prepares seven retained poses
from the embedded 24-frame table, beginning at source frame 7. Subsequent
playback loops frames 0 through 21 with the original reader's shortened
DBF count. The controller's minimum wait is 350 VBL ticks.

The first display mask addresses the previous working pose; the remaining
four address the working ring two through five positions later. Their
RGB12 colors follow the source palette permutation and signed interpolation.
The 32 half-bright colors retain the hardware's per-component bit shift.

## Verification

`original-noise-clock.csv` contains all 350 independently executed states:
software count, frame cursor, remaining frames, working-mask pointer,
material pointer, five displayed-mask pointers and all 32 base colors.
Every value matches the Go controller. Hardware raster writers alone are
omitted from this controller fixture; the source palette/transport and
animation programs remain unchanged.

The material renderer shares `maskRing` and the four-bit RGBA packing pass
with the trail and block effects. A graphics regression test independently
decodes six synthetic bitplanes on the CPU, renders the same inputs through
the actual GPU shader and compares every output byte. Its sixth material
resides on Ebitengine's image atlas, while the retained masks are unmanaged
surfaces. This catches source-coordinate mistakes that make textured fields
appear flat. Kage samples all source images in source-0 coordinates; each
source's atlas offset is already applied internally.

```sh
GOWORK=off go run ./cmd/preview -bank noise -music
GOWORK=off go run ./cmd/preview -bank noise -frame 200 -capture .local/noise
```

The original textures are rendered natively. No movie frames are embedded
or sampled during playback. Full-production loading/cue alignment and the
original line/fill boundary comparison remain part of the director work.

The [continuous native director](DIRECTOR.md) now assembles this unit with
the remaining source compositions, pages, loading programs and music handoffs.
