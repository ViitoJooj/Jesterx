import "./styles.scss"
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { Footer } from './components/Footer/Footer'
import { Header } from './components/Header/Header'
import { I18nProvider } from './i18n'
import { Home } from './pages/home'
import { NotFound } from './pages/not-found'
import { Showcase } from './pages/showcase'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <I18nProvider>
      <BrowserRouter>
          <Header />
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/design" element={<Showcase />} />
            <Route path="*" element={<NotFound />} />
          </Routes>
          <Footer />
      </BrowserRouter>
    </I18nProvider>
  </StrictMode>,
)
