import random
import time

from faker import Faker
from sqlalchemy import create_engine, text

fake = Faker()

engine = create_engine(
    "mysql+pymysql://root:password@localhost/cdc_source", pool_recycle=3600, echo=True)

NUM_ROW = 10
conn = engine.connect()
while True:
    try:
        upsert_template = '''REPLACE INTO person(person_id, fullname , num_transaction , company_id)  
        VALUES(:person_id, :fullname, :num_transaction, :company_id);'''
        conn.execute(text(upsert_template), [{
            "person_id": random.randint(1, NUM_ROW),
            "fullname": fake.name(),
            "num_transaction": random.randint(1, 100),
            "company_id": random.randint(1, 4)
        }])
        conn.execute(text(upsert_template), [{
            "person_id": random.randint(1, NUM_ROW),
            "fullname": fake.name(),
            "num_transaction": random.randint(1, 100),
            "company_id": random.randint(1, 4)
        }])
        delete_template = 'DELETE FROM person WHERE person_id=:person_id'
        conn.execute(text(delete_template), [{
            "person_id": random.randint(1, NUM_ROW)
        }])
        conn.commit()
        time.sleep(1)
    except KeyboardInterrupt:
        break

print('closing connection')
conn.close()
