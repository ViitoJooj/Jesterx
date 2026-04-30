import "./main.scss"
import "./variables.scss"
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from "react-router-dom"

import { Header } from "./components/Header/header"
import { Footer } from "./components/Footer/footer"
import { Home } from "./pages/home/home"


createRoot(document.getElementById('root')).render(
  <StrictMode>
    <BrowserRouter>
      <Header />
      <Home />
      <Footer />
    </BrowserRouter>
  </StrictMode>
)
