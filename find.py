# -*- coding:UTF-8 -*-

import sys
import mmap
import time
import contextlib

args = str(sys.argv)

if len(sys.argv) <= 1: 
    print("请设置查询参数")
    sys.exit(1)
elif len(sys.argv) <= 2:
    print("请设置查询文件")
    sys.exit(1)
elif len(sys.argv) <= 3:
    print("请设置输出文件")
    sys.exit(1)


word = sys.argv[1]
in_file_path = sys.argv[2]
out_file_path = sys.argv[3]

start =  time.time()

in_file = open(in_file_path, 'r')
out_file = open(out_file_path, 'w', 1024*1024)

size = 0
with contextlib.closing(mmap.mmap(in_file.fileno(), 0,access=mmap.ACCESS_READ)) as m:
    while True: 
        line = m.readline()
        if not line: 
                break 
        size = size + 1
        if line.find(word) >= 0:
            out_file.write(str(size) + "：" + line)
    m.close
in_file.close()
out_file.close()
seconds = time.time() - start
print(seconds)
