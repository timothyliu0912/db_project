# -*- coding:utf-8 -*-
import os
import json
import sys
print(sys.getdefaultencoding())     # 打印出目前系統字符編碼


def parse_data(fsplit):
    sub = fsplit.find('url')
    print(sub)


def to_json(fsplit):
    data_json = {}
    for i in range(len(fsplit)):
        if fsplit[i][0] == '@':


with open('ytmp', 'r') as f:
    cnt = 0
    findex = 0
    now = 0
    fsplit = []
    for line in f:
        line = line.replace('\n', '').replace('\r', '')
        if line in '@':
            cnt += 1
            to_json(fsplit)
            fsplit = []
        else:
            fsplit.append(line.decode('utf-8'))
        if(cnt > 50):
            break
