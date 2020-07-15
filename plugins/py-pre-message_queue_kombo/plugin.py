from __future__ import absolute_import, unicode_literals

import datetime

from kombu import Connection

conn, q = None, None

queueName = 'myqueue'
connectionString = 'amqp://aaa:bbb@server/abc'

def setup():
    global conn, q
    conn = Connection(connectionString)
    conn.connect()
    q = conn.SimpleQueue(queueName)

def put(message):
    q.put(message)
