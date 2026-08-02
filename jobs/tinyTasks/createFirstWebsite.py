def run(cursor):
    cursor.execute("""
        INSERT INTO websites (owner_uuid, owner_type, label, url, write_in, description)
        VALUES ('00000000-0000-0000-0000-000000000000', 'User', 'Jesterx', 'https://jesterx.com', 'React', 'Jesterx platform')
        RETURNING uuid
    """)
    return cursor.fetchone()[0]
