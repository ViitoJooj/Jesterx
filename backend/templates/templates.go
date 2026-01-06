package templates

// GetDefaultTemplate returns a default HTML template for a new page
func GetDefaultTemplate(pageType string) string {
switch pageType {
case "landing":
return `<!DOCTYPE html><html><head><title>Landing Page</title></head><body><h1>Welcome</h1></body></html>`
case "ecommerce":
return `<!DOCTYPE html><html><head><title>E-commerce</title></head><body><h1>Shop</h1></body></html>`
case "software":
return `<!DOCTYPE html><html><head><title>Software</title></head><body><h1>Software</h1></body></html>`
case "video":
return `<!DOCTYPE html><html><head><title>Video Page</title></head><body><h1>Video</h1></body></html>`
default:
return `<!DOCTYPE html><html><head><title>Page</title></head><body><h1>Page</h1></body></html>`
}
}

// GetTemplateByType is an alias for GetDefaultTemplate for backward compatibility
func GetTemplateByType(pageType string) string {
return GetDefaultTemplate(pageType)
}
