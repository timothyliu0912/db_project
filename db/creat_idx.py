from elasticsearch import Elasticsearch
es = Elasticsearch()

mappings = {                       #type_doc_test为doc_type
                    "properties": {
                        "id": {
                            "type": "keyword",
                            "index": "false"
                        },
                        "url": {
                            "type": "keyword",  # keyword不会进行分词,text会分词
                            "index": "false"  # 不建索引
                        },
                        "vid": {
                            "type": "keyword",  # keyword不会进行分词,text会分词
                            "index": "false"  # 不建索引
                        },
                        "published": {
                            "type": "text",
                            "index": "true"
                        },
                        "title": {
                            "type": "text",
                            "index": "true"
                        },
                        "content": {
                            "type": "text",
                            "index": "true"
                        },
                        "author": {
                            "type": "text",
                            "index": "true"
                        },
                        "favoritecount": {
                            "type": "long",
                            "index": "false"
                        },
                        "viewcount":{
                            "type": "long",
                            "index": "false"
                        },
                        "res":{
                            "type": "long",
                            "index": "false"
                        },
                        "duration":{
                            "type": "long",
                            "index": "false"
                        },
                        "category":{
                            "type": "text",
                            "index": "true"
                        }
                    }
        }
#es.indices.delete(index='youtube',ignore=[400,404])
res = es.indices.create(index='youtube',ignore=400)
#print(res)
res = es.indices.put_mapping(index = 'youtube', body=mappings, doc_type='video', include_type_name=True)