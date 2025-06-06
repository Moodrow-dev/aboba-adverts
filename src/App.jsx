import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { useState, useEffect } from 'react';
import CookieBanner from './components/CookieBanner';
import MainPage from './pages/MainPage';
import SellerPage from './pages/SellerPage';
import AdvertPage from './pages/AdvertPage';
import AuthorPage from './pages/AuthorPage';
import ErrorPage from './pages/ErrorPage';
import AddAdvertPage from './pages/AddAdvertPage';
import UserProfile from './pages/UserProfile';
import logo from './assets/service_logo.png';
import './App.css';

function App() {
  const [viewedAds, setViewedAds] = useState([]);

  useEffect(() => {
    const ads = JSON.parse(localStorage.getItem('viewedAds') || '[]');
    setViewedAds(ads);
  }, []);

  const trackAdView = (ad) => {
    const updated = [ad, ...viewedAds.filter(a => a.id !== ad.id)].slice(0, 5);
    setViewedAds(updated);
    localStorage.setItem('viewedAds', JSON.stringify(updated));
  };

  return (
    <Router>
      <div className="app">
        <CookieBanner />
        <Routes>
          <Route path="/" element={<MainPage viewedAds={viewedAds} logo={logo} />} />
          <Route path="/seller/vasya" element={<SellerPage logo={logo} />} />
          <Route path="/advert/:id" element={<AdvertPage trackAdView={trackAdView} logo={logo} />} />
          <Route path="/author" element={<AuthorPage logo={logo} />} />
          <Route path="/add-advert" element={<AddAdvertPage logo={logo} />} />
          <Route path="*" element={<ErrorPage logo={logo} />} />
          <Route path="/profile" element={<UserProfile />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;