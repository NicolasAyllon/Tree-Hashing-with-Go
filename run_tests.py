import subprocess

INPUTS =  [ 
          #  "input/simple.txt"   #,
          #  "input/coarse.txt"  #,
            "input/fine.txt"  
          ]

# values of i
HASH_WORKERS = [
                 1, 
                 2, 
                 4, 
                 8, 
                16, 
                -1 # means N: the number of trees
               ] 
# DATA_WORKERS = [78]               # values of j
# COMP_WORKERS = [66]

# n_hash_workers = 21
n_data_workers = 1
n_comp_workers = 66

ITERATIONS = 1

for filename in INPUTS:
  for n_hash_workers in HASH_WORKERS:
    for i in range(ITERATIONS):
        subprocess.call([
          "./BST",
          "-input={}".format(filename),
          "-hash-workers={}".format(n_hash_workers) , 
          "-data-workers={}".format(n_data_workers) ,
          "-comp-workers={}".format(n_comp_workers) ,
        ])
