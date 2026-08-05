
## quem Irá utilizar ?:
> verkoupe é um saas multi tenant para empresas e usuarios, com foco ao uso global.
> a ideia é ser similar a uma shopfy, vamos atender qualquer um que queira criar
> um ecommerce, no futuro vamos expandir pra quem quer criar landing pages, criar uma plataforma de cursos,

## Requisitos funcionais:

- RF01 - Autenticação completa.
- RF02 - Criar produtos.
- RF03 - Ter mais de um website.
- RF04 - Ter inumeros usuarios por empresa
- RF05 - Conseguir Criar roles customizadas e definir oq cada usuario pode fazer (rbac)
- RF06 - Os websites criados pelos usuarios terão auth próprio, vai funcionar básicamente como o verkoupe.
- RF07 - Suporte a vários tipos de site como SvelteKit, React, HTML&CSS e Elementor.
- RF08 - Usaremos Blob Storage para guardar os arquivos do código do usuario, porém se ele importar do github, vamos utilizar o proprio github, também séria possivel utilizar o driver caso ele permitisse e estivesse com o oauth do google ativo.
- RF09 - Cookies, JWT, Passeto e etc, jamais vão ser expostos, n podem de jeito nenhum ficar no localStorage ou Variavel, é pra sempre está protegido via https.
- RF10 - A regra de JWT é Refresh-token e Access-token para ao maximo de segurança, o token curto terá 15 min e o longo terá 90 dias.
- RF11 - Se o envio não houver X-Website o envio vai ser para a verkoupe

## Requisitos não funcionais:

- RNF01 - Rate limit.
- RNF02 - Velocidade ou desempenho do software
- RNF03 - Versão mobile.
- RNF04 - Responsividade no site,
- RNF05 - Fazer o kubernetes funcionar agora.
- RNF06 - Limpeza de banco para usuarios que apagaram a conta ou estão inativos durante um longo tempo.
- RNF07 - Design do site principal.