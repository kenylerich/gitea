#!/usr/bin/env python3
"""Restore Go source files where an upstream tool replaced a multi-byte UTF-8
sequence with a 2-byte prefix followed by a literal '?' (0x3F). Go 1.27's
source-encoding check rejects these. The fixer:

* emits the original character (em-dash U+2014 by default) in place of the orphan
  3-byte prefix;
* inserts a missing close-quote when the orphan was followed by an
  end-of-expression token or whitespace-then-newline.

Usage: python3 tools/fix-orphan-utf8.py path/...
"""
import os, sys, glob

EM_DASH = b'\xe2\x80\x94'
CLOSE_QUOTE = b'"'


def fix_file(path):
    with open(path, 'rb') as f:
        data = f.read()
    try:
        data.decode('utf-8', errors='strict')
    except UnicodeDecodeError:
        pass  # file already has orphan UTF-8; nothing to do here
    out = bytearray()
    i = 0
    n = len(data)
    fixed_a = 0  # orphan -> em-dash
    fixed_b = 0  # open-quote + em-dash + insert close-quote
    END_BYTES = {0x29, 0x2C, 0x5D, 0x3B, 0x7D}  # ) , ] ; }
    while i < n:
        b = data[i]
        if 0xE0 <= b <= 0xEF and i + 2 < n:
            b1, b2 = data[i+1], data[i+2]
            if 0x80 <= b1 <= 0xBF and b2 == 0x3F:
                j = i + 3
                while j < n and data[j] in (0x20, 0x09):
                    j += 1
                next_byte = data[j] if j < n else 0
                out.extend(EM_DASH)
                if next_byte in END_BYTES:
                    out.append(0x22)
                i += 3
                fixed_a += 1
                continue
        if 0xF0 <= b <= 0xF7 and i + 3 < n:
            b1, b2, b3 = data[i+1], data[i+2], data[i+3]
            if 0x80 <= b1 <= 0xBF and 0x80 <= b2 <= 0xBF and b3 == 0x3F:
                out.append(0x3F)
                i += 4
                fixed_a += 1
                continue
        out.append(b)
        i += 1

    out2 = bytearray()
    i = 0
    n = len(out)
    while i < n:
        b = out[i]
        if b == 0x22 and i + 4 < n and out[i+1:i+4] == EM_DASH:
            j = i + 4
            next_byte = out[j] if j < n else 0
            if next_byte != 0x22 and next_byte != 0x5C:
                out2.append(b)
                out2.extend(EM_DASH)
                out2.append(0x22)
                i += 4
                fixed_b += 1
                continue
        out2.append(b)
        i += 1

    if fixed_a or fixed_b:
        with open(path, 'wb') as f:
            f.write(bytes(out2))
    return fixed_a + fixed_b


def main(argv):
    if not argv:
        print('usage: fix-orphan-utf8.py path [path ...]', file=sys.stderr)
        return 2
    total = 0
    files = 0
    for arg in argv:
        if os.path.isdir(arg):
            candidates = glob.glob(os.path.join(arg, '**/*.go'), recursive=True)
        else:
            candidates = [arg]
        for path in candidates:
            if not path.endswith('.go'):
                continue
            n = fix_file(path)
            if n:
                total += n
                files += 1
                print(f'{path}: {n}', flush=True)
    print(f'TOTAL fixed {total} in {files} files', flush=True)
    return 0


if __name__ == '__main__':
    sys.exit(main(sys.argv[1:]))