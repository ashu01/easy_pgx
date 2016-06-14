import psycopg2
import random
import requests
import json

postgres = psycopg2.connect(database = 'tododb', host = 'localhost', user = 'todo', password = 'todo123')
cursor = postgres.cursor()
cursor.execute("""
   SELECT EXISTS(SELECT 0 FROM pg_class WHERE relname = 'todo')
""") 
# it returns true if there is at least row in database and it will be a tuple
presence = cursor.fetchone()[0]
print(presence)

if presence:
    print('todo table exists, deleting here')
    cursor.execute("""
        DROP TABLE todo
    """)

cursor.execute("""
    CREATE TABLE todo(
        id text CONSTRAINT todo_pk PRIMARY KEY,
        task text,
        complete bool
    )
""")

try:
    postgres.commit()
except Exception as e:
    postgres.rollback()
    print(e)


