#!/usr/bin/env python
import os
import re

import fire
import json

from datetime import datetime
now = datetime.now() # current date and time

version_filepath = os.path.join('.', 'version.json')
date_time = datetime.now()

def get_dict():
    with open(version_filepath, 'r') as version_file:
        version_datas = json.load(version_file)
    return version_datas


def get_version_str():
    version_datas = get_dict()
        
    version = str(version_datas["Major"] )+ "." + str(version_datas["Minor"]) + "." +  str(version_datas["Patch"])

    return version

def inc_patch(commit=""):
    version_datas = get_dict()
    version_datas["Patch"] = str(int(version_datas["Patch"]) + 1)
    version_datas["Commit"] = commit
    version_datas["BuildTime"] = now.strftime("%Y%m%d_%H:%M:%S")
    with open(version_filepath, 'w') as version_file:
        json.dump(version_datas, version_file)

def inc_minor(commit=""):
    version_datas = get_dict()
    version_datas["Minor"] = str(version_datas["Minor"] + 1)
    version_datas["Commit"] = commit
    version_datas["BuildTime"] = now.strftime("%Y%m%d_%H:%M:%S")    
    with open(version_filepath, 'w') as version_file:
        json.dump(version_datas, version_file)

def inc_major(commit=""):
    version_datas = get_dict()
    version_datas["Major"] = str(version_datas["Major"] + 1)
    version_datas["Commit"] = commit
    version_datas["BuildTime"] = now.strftime("%Y%m%d_%H:%M:%S")    
    with open(version_filepath, 'w') as version_file:
        json.dump(version_datas, version_file)

if __name__ == "__main__":
    fire.Fire({
        'get': get_version_str,
        'inc-patch': inc_patch,
        'inc-minor': inc_minor,
        'inc-major': inc_major
    })