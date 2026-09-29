# Final spirograph, variable cue cadence and closing picture

The last controller reuses the spirograph material and two-view motion,
with three final contour banks and 94 source cues. Four authored palette
sets supply the early color cycle. The display mask begins four bytes
before its working buffer, retaining the same fetch compensation as the
second spirograph passage.

At VBL count 720 the program changes its shape interrupt interval from
two ticks to four and alternates the bitplane-enable mode. The material
continues moving each VBL. Cue spans also change at the authored counters:
11, 8, 7, 6, 5, 4, 3 and finally 2 shape callbacks. The native controller
keeps this changing pace rather than replaying every bank at one fixed rate.

Later palette changes use the original word additions and `$777` masks.
VBL counts 1701–1716 interpolate the source-selected palette entries toward
white while the other three entries remain white. The controller then
unpacks a two-plane, 352-by-280 closing image and fades its four original
colors in before remaining stationary.

## Independent verification

`original-finale-clock.csv` records all 1,716 executed states: VBL/shape
counts, cue-bank cursor, remaining frames, cue span, bitplane-enable mode,
working/display pointers, the pose actually displayed, both material views,
fine scroll and all eight live colors. Every field matches the Go program.

The shared spirograph renderer uses the verified visible-pose program and
the correct zero mask-sampling shift. `ClosingPicture` also decodes the
source closing artwork and its original copper palette; it is available
as a separate diagnostic image while the full director is being assembled.

```sh
GOWORK=off go run ./cmd/preview -bank finale -music
GOWORK=off go run ./cmd/preview -bank closing
```

The uninterrupted director still needs the title/credit and illustration
handoffs, both music starts and the final closing-image fade. Individual
effect/controller tests alone do not prove whole-production alignment.

The [continuous native director](DIRECTOR.md) now assembles this unit with
the remaining source compositions, pages, loading programs and music handoffs.
