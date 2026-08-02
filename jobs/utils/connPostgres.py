import psycopg2
import os
from dotenv import load_dotenv

def db(sql, parametros=None):
    with psycopg2.connect(
        host = os.Getenv("POSTGRES_HOST"),
        database = os.Getenv("POSTGRES_DB"),
        user = os.Getenv("POSTGRES_USER"),
        password = os.Getenv("POSTGRES_PASSWORD"),
        port = os.Getenv("POSTGRES_PORT")
    ) as conn:
        with conn.cursor() as cursor:
                cursor.execute(sql, parametros)
                if cursor.description:
                    return cursor.fetchall()