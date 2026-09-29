# Cued dancer over the second spirograph passage

The side-track 110 controller reuses the same packed 640-by-512 material as
the first musical composition. It obtains new 500/477-word motion tables
from the later animation payload and drives three contour banks through an
authored 44-entry cue table. Each entry selects a bank and starting frame.

Its shape counter advances after the first two software callbacks. Multiples
of eleven switch the cue; remainder one negates the eight live RGB12 colors.
Shape count three installs the original base palette. VBL counts 828 through
844 also apply the programmed fade to four palette entries. The standalone
preview includes that complete fade. The original director's earlier minimum
shape wait is reached after 755 VBL ticks; loading/handoff timing determines
the production's actual exit.

The palette, motion transport, bitmap and parity-mask shader are shared with
the first musical composition. The source addresses this passage's display
mask four bytes before its working buffer, canceling the fetch lead. The
same renderer therefore accepts a mask-sampling offset: zero here, 32 pixels
in the first passage. The controller preserves the source's one-callback
delay between preparing and displaying a pose.

## Independent verification

`original-wave-clock.csv` stores 846 executed source states, including clocks,
bank/frame/remaining count, working/display pointers, the actual pose in the
displayed buffer, both material views, fine-scroll word and all eight live
colors. Every field matches the Go controller. Raster-writing helpers alone
are omitted; buffer ownership is tracked at the original draw/clear calls.

The first musical fixture now also records its displayed buffer's pose.
That check corrected its previously assumed two-callback delay to the actual
one-callback delay. The resident opening keeps its different buffer rotation
and two-callback behavior.

The last cue in this source table points beyond its bank header. It is not
reached by this passage's programmed duration; the native bounded preview
retains the table and stops before that unused cue rather than looping the
cue table beyond the original passage.

```sh
GOWORK=off go run ./cmd/preview -bank wave -music
GOWORK=off go run ./cmd/preview -bank wave -frame 360 -capture .local/wave
```

`NewEffect` exposes all converted units through DCK's common effect interface.
The preview host now uses that registry; per-effect controls stay in their
source adapters. The complete production director and remaining scenes are
still in progress.
