# Initial reference scene inventory

The recording is sampled at exact frame numbers rather than an fps resampling
filter. The ranges below are observations at ten-second intervals, not yet
verified source-program boundaries. Recovering the disk data will determine
the actual object banks, animations, palettes and operation counts.

| Approximate interval | Visible production unit | Candidate DCK reuse | Original-specific data to recover |
| --- | --- | --- | --- |
| 0–25 s | Opening fades and close-up hands | Image composition, palette transitions, affine poses | Hand contours, poses, fades and opening clock |
| 25–35 s | Compact author credits | Bitmap layers/pages, cue ranges | Original lettering, layout, colors and timing |
| 35–45 s | Spaceballs dragon illustration | Retained image layers and fades | Planar artwork and source palette |
| 45–60 s | First dancers over twisting colored patterns | Parity triangle batches, image masks, warped backdrops | Authored silhouette frames, patterned texture and motion mapping |
| 60–75 s | Enlarged cyan trails, then black dancer over colored blocks | Layer composition, retained trails, configurable background | Frame/pose sequence, trail offsets, palettes and block texture |
| 75–105 s | Textured dancer field, yellow patterned background and overlapping silhouettes | Masks, repeated materials, deformation layers | Source patterns, layer operation, placements and deformation clocks |
| 105–145 s | Colored or bright dancers over textured dark fields, including sliced deformation | Deformation/composition pipeline, parity masks | Dancer banks, slicing rules, gradients and timing |
| 145–165 s | Expanding or turning multicolor angular trails | Retained geometry and feedback/image transforms | Original geometry, transformations and repeated palette steps |
| 165–185 s | One and then two outlined dancers over a red textured field | Outline/filled geometry, masks and shared texture layers | Silhouette/outline derivation, two-body motion and material palette |
| 185–215 s | Close-up line contours over pastel blocks | Bounded line batches and configurable background composition | Outline frames, camera/scale, fades and block pattern |
| 215–249.754 s | Dancers over twisting monochrome pattern with colored body stripes | Shared warped backdrop, masks, layered materials | Final choreography, patterned layers, palettes and ending behavior |

The dancers and backgrounds recur in multiple compositions. Their source
animation banks and timing must remain shared; independent approximate drawings
or silhouettes traced from the recording would not establish fidelity.

## Composed native units

The opening, first musical pattern and retained-mask trail composition now
use recovered source programs rather than the approximate intervals above.
Their controllers match 961, 480 and 566 independently executed VBL states,
respectively. Their 25 Hz shape updates use software interrupts requested by
the 50 Hz VBL routine; copper display lists supply plane pointers and colors.
GPU parity/edge raster boundaries and whole-production audio alignment still
need final comparison. The remaining units are recovered data inventories,
not finished native compositions.

The block-gradient unit now also retains its six masks, three-pose-delay
shadow and complete corner-color program. Its independently verified local
counter wait is 422 VBL ticks. Native/reference pose comparison demonstrates
that full-production boundaries also include loading/decompression handoff
time; the minimum source waits alone do not prove absolute video alignment.

The following textured unit now shares the six-mask geometry renderer with
the trail and block units. Its original material transport, half-bright
palette and 350 counter states are recovered. A graphics test compares the
GPU material composition byte-for-byte with an independent planar decoder.

The second spirograph passage shares the first one's material renderer with
its own cue/pose/palette program. Its 846 recorded states also retain which
source pose is actually displayed, including the mask-pointer offset.

The recovered [loading handoffs](reference/HANDOFFS.md) include an additional
animated three-plane tile composition between the textured retained dancers
and the second spirograph. It was absent from the initial coarse inventory;
its native composition now matches the original decoded images and all
143 transport states. Its placement in the final sequence remains pending.

The side-track 114 material-trail unit now retains five source pose banks,
32 cues and the original bank-dependent trail reset. Its local 715-state
controller is independently verified. The initial ten-second observations
are intentionally not used as authoritative stage boundaries.

The following side-track 127 unit now retains its two banks, 71 cues,
counter rewinds, camera changes and palette arithmetic. Its second bank
uses an edge-endpoint exchange, requiring the new parity-span geometry
path. The 1,618-state controller is verified separately from the remaining
original line/fill raster-boundary comparison.

Native rendering identifies the side-track 134 unit as the angular ribbon
passage. Its table's earlier `close-up` inventory name remains a diagnostic
asset alias, but the composed effect is `ribbons`. Its 816-state controller
retains reversed playback, mirrored lookups, four pose pairs and palette
flashes.
