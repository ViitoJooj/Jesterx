BEGIN;

INSERT INTO users (id, email, password, plan)
VALUES
  ('e18c1e8e-e3e1-4e70-a975-9656f88b5d88', 'ana.souza@example.com',
   '$2b$10$YvE4W0r3xAMPLEhashANA000000000000000000000000000000', 'pro'),
  ('cb360a70-9662-4e2d-8428-68ef596deb19', 'bruno.lima@example.com',
   '$2b$10$YvE4W0r3xAMPLEhashBRUNO0000000000000000000000000000', 'free'),
  ('f29e049d-bb2d-4952-bb82-517207e45760', 'carla.santos@example.com',
   '$2b$10$YvE4W0r3xAMPLEhashCARLA0000000000000000000000000000', 'business'),
  ('621fab73-7581-40ee-a341-15b545740791', 'daniel.costa@example.com',
   '$2b$10$YvE4W0r3xAMPLEhashDANIEL000000000000000000000000000', 'pro'),
  ('e9c7c519-7ce0-407f-b10d-501d9dfcfcb0', 'emily.johnson@example.com',
   '$2b$10$YvE4W0r3xAMPLEhashEMILY0000000000000000000000000000', 'enterprise'),
  ('736577d0-6494-4f57-872a-b206a8d2a43e', 'john.smith@example.com',
   '$2b$10$YvE4W0r3xAMPLEhashJOHN00000000000000000000000000000', 'free');

INSERT INTO user_profiles (
    id, user_id, profile_img,
    first_name, last_name,
    phone_area_code, country_code, phone
) VALUES
  ('54916f5e-e0d4-44f1-9974-a0f6b68b6cd1',
   'e18c1e8e-e3e1-4e70-a975-9656f88b5d88',
   'https://cdn.example.com/avatars/ana.jpg',
   'Ana', 'Souza',
   '11', '+55', '998877665'),
  ('8e132024-ae8b-478e-8a35-7600d5ec64df',
   'cb360a70-9662-4e2d-8428-68ef596deb19',
   'https://cdn.example.com/avatars/bruno.jpg',
   'Bruno', 'Lima',
   '21', '+55', '997766554'),
  ('a5b9e734-96a1-4036-bc0b-4e0f5824cef1',
   'f29e049d-bb2d-4952-bb82-517207e45760',
   'https://cdn.example.com/avatars/carla.jpg',
   'Carla', 'Santos',
   '31', '+55', '996655443'),
  ('2039aa7b-9c2a-4741-853e-e5cf037b9aa8',
   '621fab73-7581-40ee-a341-15b545740791',
   'https://cdn.example.com/avatars/daniel.jpg',
   'Daniel', 'Costa',
   '41', '+55', '995544332'),
  ('9513132e-f90c-414d-9459-5a8e1aac054a',
   'e9c7c519-7ce0-407f-b10d-501d9dfcfcb0',
   'https://cdn.example.com/avatars/emily.jpg',
   'Emily', 'Johnson',
   '415', '+1', '5550134'),
  ('c2abf0c1-2ede-4d79-bbe4-a18adf2a39dd',
   '736577d0-6494-4f57-872a-b206a8d2a43e',
   'https://cdn.example.com/avatars/john.jpg',
   'John', 'Smith',
   '212', '+1', '5550172');
INSERT INTO tenants (id, name, page_id)
VALUES
  ('e1d331dd-5486-4b4b-8987-2f1f455e8bd9', 'Acme Marketing', 'acme-marketing'),
  ('5b3a4a0b-0ad5-4e92-9adb-f3210bd80436', 'Loja BoaCompra', 'lojaboa'),
  ('aa66cfdb-2ca3-4a38-aa12-6e8efbb1d3c5', 'DevSchool Online', 'devschool');
INSERT INTO tenant_users (tenant_id, user_id, role)
VALUES
  ('e1d331dd-5486-4b4b-8987-2f1f455e8bd9', 'e18c1e8e-e3e1-4e70-a975-9656f88b5d88', 'owner'),  -- Ana
  ('e1d331dd-5486-4b4b-8987-2f1f455e8bd9', 'cb360a70-9662-4e2d-8428-68ef596deb19', 'admin'),  -- Bruno
  ('e1d331dd-5486-4b4b-8987-2f1f455e8bd9', '621fab73-7581-40ee-a341-15b545740791', 'user'),   -- Daniel
  ('5b3a4a0b-0ad5-4e92-9adb-f3210bd80436', 'f29e049d-bb2d-4952-bb82-517207e45760', 'owner'),  -- Carla
  ('5b3a4a0b-0ad5-4e92-9adb-f3210bd80436', 'cb360a70-9662-4e2d-8428-68ef596deb19', 'user'),   -- Bruno
  ('aa66cfdb-2ca3-4a38-aa12-6e8efbb1d3c5', 'e9c7c519-7ce0-407f-b10d-501d9dfcfcb0', 'owner'),  -- Emily
  ('aa66cfdb-2ca3-4a38-aa12-6e8efbb1d3c5', '736577d0-6494-4f57-872a-b206a8d2a43e', 'admin');  -- John
INSERT INTO pages (
    id, tenant_id, name, page_id,
    domain, theme_id
) VALUES
  ('86bb188f-9e47-4535-980d-ac99567f04ff',
   'e1d331dd-5486-4b4b-8987-2f1f455e8bd9',
   'Página Principal', 'acme-home',
   'acmemkt.com', 'default'),
  ('7d3f0e50-dc04-44d1-9583-da07c200ebe3',
   'e1d331dd-5486-4b4b-8987-2f1f455e8bd9',
   'Campanha Black Friday', 'acme-blackfriday',
   'promo.acmemkt.com', 'dark'),
  ('6f7119f6-6d13-400f-80eb-dc7ac1736fb8',
   '5b3a4a0b-0ad5-4e92-9adb-f3210bd80436',
   'Loja Virtual', 'lojaboa-store',
   'lojaboa.com.br', 'light'),
  ('1a66dc65-50aa-48de-b384-802934419df2',
   '5b3a4a0b-0ad5-4e92-9adb-f3210bd80436',
   'Landing de Frete Grátis', 'lojaboa-frete-gratis',
   NULL, 'default'),

  ('d8795196-b82f-4317-af50-364aa905bcc2',
   'aa66cfdb-2ca3-4a38-aa12-6e8efbb1d3c5',
   'Home', 'devschool-home',
   'devschool.io', 'dark'),
  ('9513132e-f90c-414d-9459-5a8e1aac054a',
   'aa66cfdb-2ca3-4a38-aa12-6e8efbb1d3c5',
   'Página do Curso React', 'devschool-react',
   'react.devschool.io', 'code'),
  ('c2abf0c1-2ede-4d79-bbe4-a18adf2a39dd',
   'aa66cfdb-2ca3-4a38-aa12-6e8efbb1d3c5',
   'Página do Curso Node.js', 'devschool-node',
   'node.devschool.io', 'code');

COMMIT;
