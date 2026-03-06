INSERT INTO websites (id, website_type, image, name, short_description, description, creator_id, banned, updated_at, created_at)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'JESTERX',
    NULL,
    'JesterX',
    'Plataforma de criação e publicação de sites.',
    'JesterX é a plataforma principal para criação, gestão e publicação de sites e lojas virtuais.',
    '00000000-0000-0000-0000-000000000000',
    false,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO plans (name, description, description_md, price, max_sites, max_routes, billing_cycle, active)
VALUES
    (
        'Starter',
        'Para criadores independentes que querem mais controle.',
        '## Starter&#10;&#10;- 3 sites&#10;- Até 20 rotas por site&#10;- Suporte por e-mail',
        29.90, 1, 7, 'monthly', true
    ),
    (
        'Pro',
        'Para profissionais e pequenas agências.',
        '## Pro&#10;&#10;- 10 sites&#10;- Até 100 rotas por site&#10;- Suporte prioritário',
        79.90, 5, 35, 'monthly', true
    ),
    (
        'Business',
        'Para agências e empresas com alto volume.',
        '## Business&#10;&#10;- Sites ilimitados&#10;- Rotas ilimitadas&#10;- SLA e suporte dedicado',
        199.90, 15, 105, 'monthly', true
    )
ON CONFLICT (name) DO NOTHING;