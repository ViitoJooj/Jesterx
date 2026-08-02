def run(cursor, website_uuid, user_uuid):
    cursor.execute("""
        UPDATE websites SET owner_uuid = %s WHERE uuid = %s
    """, (user_uuid, website_uuid))
