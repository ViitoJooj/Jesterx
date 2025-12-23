package templates

func GetTemplateByType(pageType string) string {
	templates := map[string]string{
		"landing": `
			<script>
				let title = "Minha Landing Page";
			</script>
			<h1>{title}</h1>
			<p>Template de landing page</p>
		`,
		"ecommerce": `
			<script>
				let products = [];
			</script>
			<div class="ecommerce">
				<h1>Minha Loja</h1>
				<!-- Lista de produtos -->
			</div>
		`,
		"software": `
			<script>
				let features = [];
			</script>
			<div class="software-page">
				<h1>Meu Software</h1>
			</div>
		`,
		"video": `
			<script>
				let videoUrl = "";
			</script>
			<div class="video-page">
				<h1>Página de Vídeos</h1>
			</div>
		`,
	}

	if template, ok := templates[pageType]; ok {
		return template
	}
	return templates["landing"] // fallback
}
