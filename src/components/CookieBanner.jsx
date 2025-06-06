import { useState, useEffect } from 'react';

const CookieBanner = () => {
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    const consent = localStorage.getItem('cookieConsent');
    if (!consent) {
      setTimeout(() => setVisible(true), 1000);
    }
  }, []);

  const acceptCookies = () => {
    localStorage.setItem('cookieConsent', 'accepted');
    setVisible(false);
  };

  if (!visible) return null;

  return (
    <div className="cookie-banner">
      <p>Мы используем куки, чтобы улучшить ваш опыт. Продолжая использовать сайт, вы соглашаетесь с этим.</p>
      <button onClick={acceptCookies}>Понятно</button>
    </div>
  );
};

export default CookieBanner;