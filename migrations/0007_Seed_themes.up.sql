INSERT INTO themes (id, name, description, category, preview_url, source_type, source, active) VALUES

-- 1. Landing Page Simples
('th-0001-landing-simples',
 'Landing Page Simples',
 'Página de entrada limpa com chamada para ação destacada, ideal para campanhas e captação de leads.',
 'landing',
 '',
 'ELEMENTOR_JSON',
 '{"pages":[{"name":"Home","blocks":[{"id":"b1","type":"heading","text":"Conquiste mais clientes hoje","x":60,"y":40,"w":880,"h":60,"z":1,"rotation":0,"style":{"color":"#1a2740","fontSize":"40px","fontWeight":"700","textAlign":"center"}},{"id":"b2","type":"paragraph","text":"Solução completa para impulsionar seu negócio com resultados reais e mensuráveis.","x":120,"y":120,"w":760,"h":50,"z":1,"rotation":0,"style":{"color":"#5a6379","fontSize":"18px","textAlign":"center"}},{"id":"b3","type":"button","label":"Começar agora","action_type":"navigate","action_target":"/contato","x":380,"y":200,"w":240,"h":48,"z":1,"rotation":0,"style":{"background":"#ff5d1f","color":"#ffffff","borderRadius":"8px","fontSize":"16px","fontWeight":"600"}},{"id":"b4","type":"divider","x":60,"y":280,"w":880,"h":2,"z":1,"rotation":0,"style":{"background":"#e2e8f0"}}]}],"header":{"blocks":[]},"footer":{"blocks":[]},"popups":[]}',
 true),

-- 2. Loja Virtual
('th-0002-loja-virtual',
 'Loja Virtual',
 'Template completo para e-commerce com vitrine de produtos, destaque de ofertas e carrinho.',
 'ecommerce',
 '',
 'ELEMENTOR_JSON',
 '{"pages":[{"name":"Loja","blocks":[{"id":"b1","type":"heading","text":"Nossos Produtos","x":60,"y":30,"w":880,"h":55,"z":1,"rotation":0,"style":{"color":"#1a2740","fontSize":"36px","fontWeight":"700"}},{"id":"b2","type":"product_list","x":60,"y":110,"w":880,"h":420,"z":1,"rotation":0,"style":{"columns":3,"gap":"24px","background":"#f8fafc","borderRadius":"12px"}},{"id":"b3","type":"cart_items","x":720,"y":560,"w":220,"h":52,"z":2,"rotation":0,"style":{"background":"#1a2740","color":"#ffffff","borderRadius":"8px","fontSize":"15px","fontWeight":"600"}}]}],"header":{"blocks":[]},"footer":{"blocks":[]},"popups":[]}',
 true),

-- 3. Blog / Portfólio
('th-0003-blog-portfolio',
 'Blog / Portfólio',
 'Layout editorial para apresentar artigos, projetos e trabalhos criativos com carrossel em destaque.',
 'blog',
 '',
 'ELEMENTOR_JSON',
 '{"pages":[{"name":"Blog","blocks":[{"id":"b1","type":"heading","text":"Meu Portfólio","x":60,"y":30,"w":880,"h":55,"z":1,"rotation":0,"style":{"color":"#0f172a","fontSize":"38px","fontWeight":"800"}},{"id":"b2","type":"carousel","x":60,"y":110,"w":880,"h":340,"z":1,"rotation":0,"style":{"borderRadius":"16px","gap":"20px","autoplay":true}},{"id":"b3","type":"paragraph","text":"Projetos desenvolvidos com foco em experiência do usuário, performance e boas práticas de design.","x":60,"y":480,"w":700,"h":50,"z":1,"rotation":0,"style":{"color":"#475569","fontSize":"17px"}},{"id":"b4","type":"button","label":"Ver todos os projetos","action_type":"navigate","action_target":"/projetos","x":60,"y":555,"w":220,"h":44,"z":1,"rotation":0,"style":{"background":"#6366f1","color":"#ffffff","borderRadius":"8px","fontSize":"15px"}}]}],"header":{"blocks":[]},"footer":{"blocks":[]},"popups":[]}',
 true),

-- 4. Perfil Pessoal
('th-0004-perfil-pessoal',
 'Perfil Pessoal',
 'Página de apresentação pessoal ou profissional com card de perfil, bio e links de contato.',
 'profile',
 '',
 'ELEMENTOR_JSON',
 '{"pages":[{"name":"Perfil","blocks":[{"id":"b1","type":"profile_card","x":340,"y":30,"w":320,"h":200,"z":1,"rotation":0,"style":{"borderRadius":"50%","border":"4px solid #6366f1","background":"#f1f5f9"}},{"id":"b2","type":"heading","text":"Olá, eu sou João Silva","x":60,"y":260,"w":880,"h":50,"z":1,"rotation":0,"style":{"color":"#0f172a","fontSize":"34px","fontWeight":"700","textAlign":"center"}},{"id":"b3","type":"paragraph","text":"Desenvolvedor Full-Stack apaixonado por criar experiências digitais incríveis. +5 anos de experiência em produtos SaaS.","x":120,"y":330,"w":760,"h":60,"z":1,"rotation":0,"style":{"color":"#64748b","fontSize":"17px","textAlign":"center"}},{"id":"b4","type":"button","label":"Entre em contato","action_type":"navigate","action_target":"/contato","x":360,"y":420,"w":280,"h":48,"z":1,"rotation":0,"style":{"background":"#6366f1","color":"#ffffff","borderRadius":"24px","fontSize":"16px","fontWeight":"600"}}]}],"header":{"blocks":[]},"footer":{"blocks":[]},"popups":[]}',
 true),

-- 5. SaaS / Software
('th-0005-saas-software',
 'SaaS / Software',
 'Template moderno para produtos de software com hero section, benefícios e chamada para teste grátis.',
 'saas',
 '',
 'ELEMENTOR_JSON',
 '{"pages":[{"name":"Home","blocks":[{"id":"b1","type":"heading","text":"O software que seu time precisa","x":60,"y":40,"w":880,"h":60,"z":1,"rotation":0,"style":{"color":"#0f172a","fontSize":"42px","fontWeight":"800","textAlign":"center"}},{"id":"b2","type":"paragraph","text":"Automatize processos, aumente a produtividade e tome decisões baseadas em dados com nossa plataforma.","x":100,"y":125,"w":800,"h":55,"z":1,"rotation":0,"style":{"color":"#475569","fontSize":"18px","textAlign":"center"}},{"id":"b3","type":"button","label":"Teste grátis por 14 dias","action_type":"navigate","action_target":"/signup","x":300,"y":210,"w":400,"h":52,"z":1,"rotation":0,"style":{"background":"#0ea5e9","color":"#ffffff","borderRadius":"10px","fontSize":"17px","fontWeight":"700"}},{"id":"b4","type":"divider","x":60,"y":295,"w":880,"h":2,"z":1,"rotation":0,"style":{"background":"#e2e8f0"}},{"id":"b5","type":"paragraph","text":"✓ Sem cartão de crédito  ✓ Configuração em 5 minutos  ✓ Suporte 24/7","x":100,"y":320,"w":800,"h":40,"z":1,"rotation":0,"style":{"color":"#64748b","fontSize":"15px","textAlign":"center"}}]}],"header":{"blocks":[]},"footer":{"blocks":[]},"popups":[]}',
 true);
