import subprocess

INPUTS =  [ 
          # "input/simple.txt"   #,
          #  "input/coarse.txt"   #,
            "input/fine.txt"  
          ]

inputSizes = {
  "input/simple.txt": 10,
  "input/coarse.txt": 100,
  "input/fine.txt": 100000
}

HASH_WORKERS = [1, 2, 4, 8, 16]
# DATA_WORKERS = [78]
# COMP_WORKERS = [66]

# n_hash_workers = 21
n_data_workers = 1
n_comp_workers = 66

ITERATIONS = 1

for filename in INPUTS:
  # Insert total number of trees (N) as a possible number of threads
  # HASH_WORKERS.append(inputSizes[filename])
  for n_hash_workers in HASH_WORKERS:
    for i in range(ITERATIONS):
        subprocess.call([
          "./BST",
          "-input={}".format(filename),
          "-hash-workers={}".format(n_hash_workers) , 
          "-data-workers={}".format(n_data_workers) ,
          "-comp-workers={}".format(n_comp_workers) ,
        ])
