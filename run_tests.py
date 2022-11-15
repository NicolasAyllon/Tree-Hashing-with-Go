import subprocess

import os
from subprocess import check_output
import re

INPUTS =  [ 
          #  "input/simple.txt"   ,
           "input/coarse.txt"   ,
          #  "input/fine.txt"  
          ]
# values of i
HASH_WORKERS = [
                #  1, 
                #  2, 
                #  4, 
                #  8, 
                 16, 
                # -1 # means N: the number of trees
               ] 
# values of j
DATA_WORKERS = {
#  i : possible values of j
# hash-workers : [data-workers]
                1: [
                    1
                   ],
                2: [
                    1, 
                    2
                   ],
                4: [
                    1, 
                    4
                   ],
                8: [
                    1, 
                    8
                   ],
               16: [
                    #  1,
                    16
                   ],
               -1: [
                     1, 
                    -1
                   ]
}
COMP_WORKERS = [
                # 1, 
                # 2, 
                # 4, 
                8, 
                # 16,
                # -1
               ]
ITERATIONS = 1

for filename in INPUTS:
  for n_hash_workers in HASH_WORKERS:
    for n_data_workers in DATA_WORKERS[n_hash_workers]:
      for n_comp_workers in COMP_WORKERS:
        # csvs = []
        for i in range(ITERATIONS):
          # csv = ["{}:{}/{}/{}".format(filename, n_hash_workers, n_data_workers, n_comp_workers)]
            subprocess.call([
              "./BST",
              "-input={}".format(filename),
              "-hash-workers={}".format(n_hash_workers) , 
              "-data-workers={}".format(n_data_workers) ,
              "-comp-workers={}".format(n_comp_workers) ,
            ])

          # Only hashes (step 1)
          # cmd = "./BST -input={} -hash-workers={}".format(filename, n_hash_workers)

        # All steps: 1, 2, 3
        #   cmd = "./BST -input={} -hash-workers={} -data-workers={} -comp-workers={}".format(filename, n_hash_workers, n_data_workers, n_comp_workers)

        #   out = check_output(cmd, shell=True).decode("utf8")

        #   m_hash = re.search("hashTime = (.*)", out)
        #   if m_hash is not None:
        #     hashTime = m_hash.group(1)
        #     csv.append(hashTime)
        #   m_group = re.search("hashGroupTime = (.*)", out)
        #   if m_group is not None:
        #     hashGroupTime = m_group.group(1)
        #     csv.append(hashGroupTime)
        #   m_comp = re.search("compareTreeTime = (.*)", out)
        #   if m_comp is not None:
        #     compTreeTime = m_comp.group(1)
        #     csv.append(compTreeTime)
      
        #   csvs.append(csv)

        # print("\n")
        # for csv in csvs:
        #   print (", ".join(csv))
