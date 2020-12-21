# -*- coding:utf-8 -*-
import os
import json
import sys
import datetime
import time
print(sys.getdefaultencoding())     # 打印出目前系統字符編碼

idx = ['@url','@published','@title','@content','@author','@favoriteCount','@viewCount','@res','@duration','@category']
videos = []
def to_json(fsplit):
    cnt = 0
    for i in fsplit:
        if i[0] != '@':
            return False
        for j in idx:
            if i.find(j) != -1:
                cnt+=1
                continue
    if cnt != 10:
        return False
    l = []
    print(fsplit)
    for i in fsplit:
        o = i.split(':',1)
        if o[0] == '@url':
            s = "https://www.youtube.com/embed/"
            url_s = o[1].split('watch?v=')
            s = s+ url_s[1]
            print(s)
        if (o[0] =='@author')&(o[1]==''):
            print(len(o[1]))
        l.append(o[0][1:])
        l.append(o[1])
    video = dict(zip(l[0::2],l[1::2]))
    print(video)
    # video['duration'] = int(video['duration'])
    # #video['viewcount'] = int(video['viewcount'])
    # video['favoriteCount'] = int(video['favoriteCount'])
    # video['res'] = int(video['res'])
    # print(video)
    return True


with open('ytmp', 'r') as f:
    cnt = 0
    findex = 0
    now = 0
    succ = 0
    input_str = []
    for line in f:
        line = line.strip()
        if (line == "@") & (len(line) == 1):
             cnt += 1
             if cnt > 1:
                ret = to_json(input_str)
                if ret == True:
                    succ+=1
                input_str = []
        elif len(line) == 0:
            continue
        else:
            input_str.append(line)
        if succ > 20:
            print(succ,cnt)
            break
