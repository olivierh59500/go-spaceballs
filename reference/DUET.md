# Outline duet and the white material exit

The side-track 139 director sequences the two 256-frame body banks through
51 source cues. Its shape counter begins at 44; the first two software
callbacks only advance the material phase. Startup prepares the six-mask
ring, and active callbacks draw the next outline pose into a working mask.

The source selects different trail lengths depending on both the shape
phase and the cue pointer. Outside one pointer interval it retains all four
older trail masks. Inside the interval it progressively redirects the five
display planes and alternates two RGB12 palette programs. The outline
geometry comes from the banks' outer/inner contours, using shared parity
filling rather than separately drawn generic dancer outlines.

The last source cues repeatedly select the last pose. Once the shape counter
reaches 552, the copper list switches from six planes to a single material
plane. Two original colors interpolate toward white over the final callbacks.
The native effect retains that brief white exit as part of the unit.

## Independent verification

`original-duet-clock.csv` contains all 1,041 source VBL states. Tests compare
software/shape counts, bank/frame/remaining count, cue cursor, working and
display pointers, material phase, the final material pointer, both exit
colors and all 32 base colors. Every field matches the original CPU output.

The material/retained-mask renderer is shared with the preceding native
passages. Captures reproduce the gold outlined pose over the dark red
textured field. Full-production audio/cue alignment and the original
line/fill raster boundaries remain under verification.

```sh
GOWORK=off go run ./cmd/preview -bank duet -music
GOWORK=off go run ./cmd/preview -bank duet -frame 400 -capture .local/duet
```
