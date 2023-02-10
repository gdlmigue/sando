import oss2
import re

from datetime import datetime


endpoint = "oss-cn-shanghai.aliyuncs.com"

access_key_id = "LTAI5tNw8NQjPw7ffBuEujgt"
access_key_secret = "2nTlKpCA0eV0ZrtpoxcXTfrU7dX1IB"

auth = oss2.Auth(access_key_id, access_key_secret)
bucket = oss2.Bucket(auth, endpoint, "cdn-logs-prod-cn-shanghai")

date_pattern = r".+(_2022_11_11_\d{6}_\d{6}).gz"

from sample_common import MNSSampleCommon
from mns.account import Account
from mns.queue import *

import pprint

accid, acckey, endpoint, token = MNSSampleCommon.LoadConfig()

my_account = Account(endpoint, accid, acckey, token)
queue_name = "cdn-logs-cn"
my_queue = my_account.get_queue(queue_name)

a = 0
cnt = 0
for object_info in oss2.ObjectIterator(bucket):
    if re.search(date_pattern, object_info.key):
        a += 1
        print(object_info.key)

print("Runnning")
for object_info in oss2.ObjectIterator(bucket):
    if re.search(date_pattern, object_info.key):
        print(object_info.key)
        cnt += 1

print(cnt)
        # try:
        #     msg_body = {
        #         "events": [
        #             {
        #                 "eventName": "ObjectCreated:UploadPart",
        #                 "eventSource": "acs:oss",
        #                 "eventTime":   "2022-12-01T10:37:50.000Z",
        #                 "eventVersion": "1.0",
        #                 "oss": {
        #                     "bucket": {
        #                         "arn": "acs:oss:cn-shanghai:5799651911776518:cdn-logs-prod-cn-shanghai",
        #                         "name": "cdn-logs-prod-cn-shanghai",
        #                         "ownerIdentity": "5799651911776518",
        #                         "virtualBucket": "",
        #                     },
        #                     "object": {
        #                         "eTag": object_info.etag[:-2],
        #                         "key":  object_info.key,
        #                         "size": object_info.size,
        #                     },
        #                     "ossSchemaVersion": "1.0",
        #                     "ruleId": "CDNLogsProd",
        #                 },
        #                 "region": "cn-shanghai",
        #                 "requestParameters": {"sourceIPAddress": "172.20.17.185"},
        #                 "responseElements": {"requestId": "62B0DD0A41D0953233373852"},
        #                 "userIdentity": {"principalId": "338061994335656722"},
        #             }
        #         ]
        #     }

        #     import json

        #     msg_body = json.dumps(msg_body)

        #     msg = Message(msg_body)
        #     re_msg = my_queue.send_message(msg)
        #     pprint.pp(msg_body)

        # except MNSExceptionBase as e:
        #     if e.type == "QueueNotExist":
        #         print("Queue not exist, please create queue before send message.")
        #         sys.exit(0)
        #     print("Send Message Fail! Exception:%s\n" % e)
