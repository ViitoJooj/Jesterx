import "./styles/global.scss";
import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Header } from "./components/Header";
import { Login } from "./pages/Login";
import { Register } from "./pages/Register";
import { UserProvider } from "./config/UserContext";
import { ProtectedRoute } from "./config/ProtectedRoute";
import { Footer } from "./components/Footer";
import { NotFound } from "./pages/Not-found";
import { HomePage } from "./pages/HomePage";
import { SpeedInsights } from "@vercel/speed-insights/react";
import { Pages } from "./pages/Pages";
import { Products } from "./pages/Products";
import { Pricing } from "./pages/Pricing";
import { PaymentSuccess } from "./pages/PaymentSuccess";
import { PaymentCancel } from "./pages/PaymentCancel";

const root = ReactDOM.createRoot(document.getElementById("root") as HTMLElement);

root.render(
  <BrowserRouter>
    <UserProvider>
      <Header />
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/404" element={<NotFound />} />
        <Route path="*" element={<NotFound />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/pages" element={<Pages />} />
        <Route path="/products" element={<Products />} />
        <Route path="/pricing" element={<Pricing />} />
        <Route path="/payment-success" element={<PaymentSuccess />} />
        <Route path="/payment-cancel" element={<PaymentCancel />} />
      </Routes>
      <Footer />
    </UserProvider>
  </BrowserRouter>
);
