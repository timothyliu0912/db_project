import pymongo
from datetime import datetime
import json
import time

myclient = pymongo.MongoClient("mongodb://localhost:27017/")
mydb = myclient["youtube"]
mycol = mydb["videos"]
 
mylist = []
with open('a.json') as f:
    json_from_file = json.load(f)


cnt = 1
data_cnt = 0
s = time.time()
for i in json_from_file['videos']:
    # print(i)
    mylist.append(i)
    data_cnt+=1
    cnt+=1
    if data_cnt == 1000:
        x = mycol.insert_many(mylist)
        mylist = []
        data_cnt = 0
        print('cnt',cnt)
if len(mylist) != 0:
    x = mycol.insert_many(mylist)
e = time.time()
print("{}s".format(e-s))

f.close()