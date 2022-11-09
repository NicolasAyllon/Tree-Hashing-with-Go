# This script was generously provided by Zesen Huang
# https://piazza.com/class/l5wvxtgrn3l51e/post/335

import sys
import os
import hashlib
from functools import reduce
result = []
with open(sys.argv[1]) as f:
    for line in f.readlines():
        line = line.strip()
        if line.startswith('group '):
            nums = line.split(':')
            nums = [int(x) for x in nums[1].strip().split(' ')]
            result.append(reduce(lambda x, y: x + y, nums))
result.sort()
text = reduce(lambda x, y: str(x) + str(y), result)
print(hashlib.md5(text.encode(encoding='UTF-8')).hexdigest())

# How to use:
# Step 1:
# go run  src/*.go -input=input/simple.txt -hash-workers=8 -data-workers=8
# -comp-workers=8 > output.txt
# Step 2:
# python output_check.py output.txt
