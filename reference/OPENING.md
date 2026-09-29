# Resident opening composition

The original opening is not a single animation on a flat background. It uses
two bitplanes, three working buffers and palettes that can make one of the
planes invisible. Treating the 142-frame bank as a standalone picture loses
the foreground hands and the intermediate morph composition.

An additional **66-frame table** starts inside the resident program at
`$3fec6` (file offset `$12c6`). Its first 57 frames supply the initial hands.
The recovered inventory therefore contains 33 tables and 3,386 frames.

The VBL routine runs at 50 Hz and requests a level-1 software interrupt on odd VBL
counts. Each software callback prepares one contour at 25 Hz. The source exits
the resident opening at VBL count 961, before the title/credit director.

| Shape counts | Source table | Working/display-plane behavior |
| --- | --- | --- |
| 1–57 | Embedded hands, frames 0–56 | Rotate both display planes; clear and fill both working planes |
| 58–199 | 142-frame bank at `$80000` | Hold the first display plane; rotate, clear and fill the second plane only |
| 200–479 | Bank at `$83dfe`, starting at frame 55 | Rotate both planes; stop clearing/filling the second plane after VBL count 580 |
| 480–481 | Exhausted source table | Retain the final source buffer/plane behavior until the director handoff |

The middle passage uses a second coordinate lookup: X is scaled by 351/256
and shifted by -15; Y is scaled by 350/204, shifted by +8 and capped at 289.
Other poses use the 351/256 and 289/204 lookups. The visible opening window
is 352 by 290 pixels, with a 44-byte plane row and no extra fetch cropping.

The native effect retains six persistent plane targets. DCK batches the
parity-filled contours; the late retained plane receives contour edges, and
an Ebitengine shader performs the original two-bit palette lookup. Seeking
backward clears the working targets and reconstructs the same state.

## Independent checks

`internal/source/testdata/original-opening-clock.csv` stores all 961 source
states. Ghidra executes the actual VBL/software routines in the **initial**
memory image, before the dragon director overwrites resident code. Hardware
raster helpers and the music update are omitted; source bank switching,
buffer rotation, held-plane pointers and palette interpolation stay intact.
Tests compare every state with the Go controller.

Native captures now contain the foreground hands and their internal masked
composition. Source clock/palette verification is complete for this opening;
the GPU polygon/edge rasterizer remains a close reconstruction rather than
an independently proven, pixel-identical replacement for Amiga line/fill DMA.
The movie's initial disk-loading delay is not part of this standalone effect.

```sh
GOWORK=off go run ./cmd/preview -bank opening-effect
GOWORK=off go run ./cmd/preview -bank opening-effect -frame 230 -capture .local/opening
GOWORK=off go run ./cmd/preview -bank hands -frame 20
```

The [continuous native director](DIRECTOR.md) now assembles this unit with
the remaining source compositions, pages, loading programs and music handoffs.
