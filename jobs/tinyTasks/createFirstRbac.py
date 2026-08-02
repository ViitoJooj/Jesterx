from utils.connPostgres import db

def run(cursor, website_uuid):
    cursor.execute("""
        INSERT INTO rbac (website_uuid, label, can_read, can_write, can_update, can_upgrade, can_delete)
        VALUES (%s, 'admin', TRUE, TRUE, TRUE, TRUE, TRUE)
        RETURNING uuid
    """, (website_uuid,))
    return cursor.fetchone()[0]
