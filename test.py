import sys

for line in sys.stdin:
    print(int(line) + 1)
    sys.stdout.flush()
