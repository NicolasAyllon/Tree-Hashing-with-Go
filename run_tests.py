import subprocess

INPUTS =  [ 
          #  "input/simple.txt"   #,
          #  "input/coarse.txt"   #,
            "input/fine.txt"  
          ]

HASH_WORKERS = [21]
DATA_WORKERS = [78]
COMP_WORKERS = [66]

n_hash_workers = 21
n_data_workers = 78
n_comp_workers = 66

ITERATIONS = 1

for filename in INPUTS:
    for i in range(ITERATIONS):
        subprocess.call([
          "./BST",
          "-input={}".format(filename),
          "-hash-workers={}".format(n_hash_workers) , 
          "-data-workers={}".format(n_data_workers) ,
          "-comp-workers={}".format(n_comp_workers) ,
        ])
