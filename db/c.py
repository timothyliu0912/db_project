from datetime import datetime
from elasticsearch import Elasticsearch
from elasticsearch import helpers
import json
import time
es = Elasticsearch()
f = open("yt_data.rst")
lines = f.readlines()
cnt = 1
data_cnt = 0
actions = []
s = time.time()
for line in lines:
    data = json.loads(line)
    action = {
        "_index": "youtube",
        "_id":cnt,
        "_source": data
    }
    actions.append(action)
    data_cnt+=1
    if data_cnt == 20000:
        a = helpers.bulk(es, actions)
        actions = []
        data_cnt = 0
    cnt+=1
    print(cnt)
a = helpers.bulk(es, actions)
e = time.time()
print("{}s".format(e-s))

f.close()