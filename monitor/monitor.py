import time
import pandas as pd
import streamlit as st
from sqlalchemy import create_engine, text, Connection
from streamlit_autorefresh import st_autorefresh


def get_source(conn: Connection) -> pd.DataFrame:
    return pd.read_sql_table('person', conn)

def get_reference(conn: Connection) -> pd.DataFrame:
    return pd.read_sql_table('company', conn)

def get_target(conn: Connection) -> pd.DataFrame:
    sql = text('select * from person order by person_id')
    return pd.read_sql_query(sql, conn)


def main():
    st.set_page_config(page_title="NRT CDC STREAM", layout="wide")
    st.title("NRT CDC STREAM")

    count = st_autorefresh(interval=2000, limit=10000, key="refresh")

    mysql_conn = create_engine("mysql+pymysql://root:password@localhost/cdc_source").connect()
    mariadb_conn = create_engine("mysql+pymysql://root:password@localhost:3307/cdc_reference").connect()
    postgres_conn = create_engine("postgresql+psycopg2://admin:password@localhost/cdc_target").connect()
    
    col1, col2, col3 = st.columns([3, 1, 3])

    try: 
        if count:
            col1.subheader("Source (mysql:8)")
            col1.table(get_source(mysql_conn))
            col2.subheader("Reference (mariadb:10)")
            col2.table(get_reference(mariadb_conn))
            col3.subheader("Target (postgres:14)")
            col3.table(get_target(postgres_conn))
    except KeyboardInterrupt:
        mysql_conn.close()
        mariadb_conn.close()
        postgres_conn.close()
    


if __name__ == '__main__':
    main()

