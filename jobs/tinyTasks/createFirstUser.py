import bcrypt

def run(cursor, website_uuid, rbac_uuid, email, password):
    hashed = bcrypt.hashpw(password.encode(), bcrypt.gensalt()).decode()

    cursor.execute("""
        INSERT INTO users (website_uuid, name, email, role, password)
        VALUES (%s, 'Admin', %s, %s, %s)
        RETURNING uuid
    """, (website_uuid, email, rbac_uuid, hashed))

    return cursor.fetchone()[0]
