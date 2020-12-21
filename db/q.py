from datetime import datetime
from elasticsearch import Elasticsearch
from elasticsearch import helpers
import json
import time
es = Elasticsearch()
f = open('wlist.txt','r')
a = []
for i in f:
    i = i.strip('\n')
    a.append(i)
s = time.time()
for i in range(20000):
    body = {
    "query":{
        "query_string":{
            "query":a[i],
        },
    },
    "size":1000,
}
    res = es.search(index="youtube",body=body)
e = time.time()
print("{}s".format(e-s))