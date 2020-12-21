from elasticsearch import Elasticsearch
import json
import numpy as np
from elasticsearch import helpers
import hashlib
from datetime import datetime
import time

md5 = hashlib.md5()
es = Elasticsearch(['localhost:9200'])
with open('a.json') as f:
    json_from_file = json.load(f)
print('success open') 
actions = []
cnt = 0
data_cnt = 0
s = time.time()
for i in json_from_file['videos']:
    md5.update(i['url'].encode("utf-8"))
    hash_md5 = md5.hexdigest()
    action = {
        "_index": "youtube",
        "_id":hash_md5,
        "_source": i
    }
    cnt += 1
    data_cnt+=1
    actions.append(action)
    if cnt == 1000:
        res = es.index(index="youtube",body = i)
        a = helpers.bulk(es, actions)
        actions = []
        cnt = 0
        print("success",data_cnt)

a = helpers.bulk(es, actions)
e = time.time()
print("{}s".format(e-s))
