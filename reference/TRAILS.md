# Retained-mask trail composition

The side-track 85 controller at `$45000` displays five independently
addressed bitplanes drawn from a ring of six one-bit masks. Their 32-color
palette controls the apparent layered depth. The source has 27 cues, each
selecting a contour bank and an initial frame. No trail shape, color sequence
or camera movement needs to be invented for the native effect.

The recovered controller starts its shape counter at **40**. Initialization
selects/clears seven successive ring masks, then builds the authored palette
through an explicit permutation table and signed RGB12 interpolation.
The VBL routine requests a level-1 software interrupt every second tick.
The first three software callbacks only prepare the mask ring. Subsequent
callbacks advance the shapes, switch cues at multiples of eleven and draw
the next contour into the new working mask. The controller hands off at VBL
count 566 (11.32 seconds).

Each display update initially addresses all five planes at the same mask.
Its modulo-22 phase progressively redirects up to four of them toward older
masks. At phase zero, later passages invert all 32 RGB12 colors. The bank
at `$47212` also selects a clipped close-up lookup when the shape counter's
modulo-88 phase lies between 22 and 43. That lookup scales X by 580/256,
shifts it by -30, and scales Y by 450/204 with a -100 shift.

## Independent verification

`original-trails-clock.csv` records all 566 original controller states:
VBL/software/shape counts, selected bank and frame, remaining frames,
working-mask address, active coordinate lookup, five display-mask addresses
and all 32 live colors. The Go controller matches every field. The separate
`original-trails-morphs.csv` fixture records 110 executed morph contours,
including repeated cues and close-up passages; their prepared coordinates
and masks match the Go reader.

These fixtures execute the source routines in the memory stage that loads
the side-track 85 block. Only raster-writing helpers are omitted. The
native renderer uses the recovered contours and DCK's parity batches, six
persistent mask targets and two small shader passes. Four mask bits are
packed into RGBA so the final five-plane palette lookup fits Ebitengine's
four-source limit. No CPU pixel readback occurs during playback.

```sh
GOWORK=off go run ./cmd/preview -bank trails -music
GOWORK=off go run ./cmd/preview -bank trails -frame 260 -capture .local/trails
```

Music is opened by DCK from the recovered `.mod` bytes. The standalone
preview starts it at the converted effect's local main-module offset; whole
production handoff/loading-delay alignment remains part of the director
work. The original line/fill raster boundary comparison also remains open.

The [continuous native director](DIRECTOR.md) now assembles this unit with
the remaining source compositions, pages, loading programs and music handoffs.
