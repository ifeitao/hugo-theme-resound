#!/usr/bin/env python3
"""Regenerate the site's logo / favicon assets from the SecEvo Labs mark.

  python3 brand/make-icons.py

Sources sit next to this script (the same geometry as secevo.xyz's brand kit:
three nested arcs, all openings at 315°, ember red dot — 200-unit canvas).
Outputs land in static/, which is served as-is by Hugo.

Requirements: headless Google Chrome (rendering), Pillow (ICO assembly).
"""
import os
import struct
import subprocess
import tempfile

from PIL import Image

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
BRAND = os.path.dirname(os.path.abspath(__file__))
STATIC = os.path.join(ROOT, "static")
CHROME = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"

GREEN = "#1F4D3A"
PAPER = "#F5F2EA"


def read(name):
    with open(os.path.join(BRAND, name), encoding="utf-8") as fh:
        return fh.read()


def shoot(svg, out, size, background=None, inset=0.0):
    """Render `svg` into a size×size PNG. inset is a fraction of the canvas."""
    w = h = size
    pad = round(size * inset / 2)
    inner = size - 2 * pad
    body = (
        f'<div style="position:absolute;left:{pad}px;top:{pad}px;'
        f'width:{inner}px;height:{inner}px">{svg}</div>'
    )
    bg = background or "transparent"
    html = (
        '<html><head><style>html,body{margin:0;padding:0;background:%s}'
        "svg{display:block;width:100%%;height:100%%}</style></head>"
        '<body style="width:%dpx;height:%dpx;overflow:hidden">%s</body></html>'
        % (bg, w, h, body)
    )
    fd, path = tempfile.mkstemp(suffix=".html")
    with os.fdopen(fd, "w", encoding="utf-8") as fh:
        fh.write(html)
    args = [
        CHROME, "--headless=new", "--disable-gpu", "--hide-scrollbars",
        "--force-device-scale-factor=1", f"--window-size={w},{h}",
        "--default-background-color=00000000",
        f"--screenshot={out}", "--virtual-time-budget=1200", f"file://{path}",
    ]
    subprocess.run(args, capture_output=True, check=False)
    os.unlink(path)
    image = Image.open(out)
    image.load()
    if image.size != (w, h):  # Chrome pads the viewport by a few px now and then
        image.crop((0, 0, w, h)).save(out)
    return out


def png_ico(paths, out):
    """Write a multi-size .ico whose entries are the PNGs themselves."""
    images = [open(p, "rb").read() for p in paths]
    dims = [Image.open(p).size[0] for p in paths]
    header = struct.pack("<HHH", 0, 1, len(images))
    offset = 6 + 16 * len(images)
    entries, blobs = b"", b""
    for data, d in zip(images, dims):
        entries += struct.pack(
            "<BBBBHHII", d % 256, d % 256, 0, 0, 1, 32, len(data), offset
        )
        blobs += data
        offset += len(data)
    with open(out, "wb") as fh:
        fh.write(header + entries + blobs)


def main():
    tile = read("icon-tile.svg")  # rounded green tile, three rings
    mini = read("icon-tile-mini.svg")  # same tile, single ring: the ≤16px tier
    disc = read("icon-disc.svg")  # full-bleed green disc, used for logo/app

    # Favicons: the tile, transparent outside its rounded corners. 16px gets the
    # one-ring tier — three rings collapse into grey mush below ~20px.
    with open(os.path.join(STATIC, "favicon.svg"), "w", encoding="utf-8") as fh:
        fh.write(tile)

    ico_parts = []
    for size in (16, 32, 96):
        out = os.path.join(STATIC, f"favicon-{size}x{size}.png")
        shoot(mini if size == 16 else tile, out, size)
        if size in (16, 32):
            ico_parts.append(out)
    big = shoot(mini, os.path.join(tempfile.gettempdir(), "ico-48.png"), 48)
    ico_parts.append(big)
    png_ico(ico_parts, os.path.join(STATIC, "favicon.ico"))
    os.unlink(big)

    # Logo: green disc, edge to edge — the masthead clips it into a circle,
    # and the same file doubles as the social-card fallback.
    shoot(disc, os.path.join(STATIC, "logo.png"), 500)

    # Install icons: disc on paper, inset so Android's maskable crop keeps it whole.
    shoot(disc, os.path.join(STATIC, "apple-touch-icon.png"), 180,
          background=PAPER, inset=0.10)
    for size in (192, 512):
        shoot(disc, os.path.join(STATIC, f"web-app-manifest-{size}x{size}.png"),
              size, background=PAPER, inset=0.16)

    for name in sorted(os.listdir(STATIC)):
        p = os.path.join(STATIC, name)
        print(f"{name:34} {os.path.getsize(p):>8,} bytes")


if __name__ == "__main__":
    main()
