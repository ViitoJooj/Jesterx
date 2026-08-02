import os
from dotenv import load_dotenv

load_dotenv(dotenv_path='./.env')

from utils.connPostgresRaw import conn
from tinyTasks import createFirstWebsite, createFirstRbac, createFirstUser, updateWebsiteOwner

conexao = conn()
cursor = conexao.cursor()

try:
    website_uuid = createFirstWebsite.run(cursor)
    print(f"website created: {website_uuid}")

    rbac_uuid = createFirstRbac.run(cursor, website_uuid)
    print(f"rbac created: {rbac_uuid}")

    user_uuid = createFirstUser.run(cursor, website_uuid, rbac_uuid, "admin@jesterx.com", "admin123")
    print(f"user created: {user_uuid}")

    updateWebsiteOwner.run(cursor, website_uuid, user_uuid)
    print(f"owner updated: {website_uuid} -> {user_uuid}")

    conexao.commit()
    print("bootstrap completed")
except Exception as e:
    conexao.rollback()
    print(f"bootstrap failed: {e}")
    raise
finally:
    cursor.close()
    conexao.close()
