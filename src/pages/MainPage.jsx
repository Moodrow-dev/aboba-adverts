import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import Cookies from 'js-cookie';
import PageLayout from '../components/PageLayout';
import AdvertList from '../components/AdvertList';
import MarkdownContent from '../components/MarkdownContent';

const MainPage = ({ logo }) => {
  const [viewedAds, setViewedAds] = useState([]);
  const [hotAdverts, setHotAdverts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    // Load viewed ads from cookies
    const loadViewedAds = () => {
      try {
        const adsFromCookie = Cookies.get('viewed_ads');
        if (adsFromCookie) {
          const parsedAds = JSON.parse(adsFromCookie);
          const validAds = parsedAds.filter(ad => (
              ad?.id &&
              ad?.title &&
              ad?.seller &&
              ad?.time
          ));
          console.log('Viewed Ads:', validAds); // Debug log
          setViewedAds(validAds);
        }
      } catch (e) {
        console.error('Ошибка загрузки просмотренных объявлений:', e);
        Cookies.remove('viewed_ads');
      }
    };

    // Fetch hot adverts from server
    const fetchHotAdverts = async () => {
      try {
        setLoading(true);
        const response = await fetch('http://localhost:4242/api/adverts');
        if (!response.ok) {
          throw new Error('Ошибка загрузки объявлений');
        }
        const data = await response.json();
        console.log('Hot Adverts:', data); // Debug log
        setHotAdverts(Array.isArray(data) ? data : []);
      } catch (err) {
        console.error('Ошибка:', err);
        setError(err.message);
        setHotAdverts([]);
      } finally {
        setLoading(false);
      }
    };

    loadViewedAds();
    fetchHotAdverts();
  }, []);

  // Normalize advert data with fallback values
  const normalizeAd = (ad) => ({
    id: ad.id || 0,
    title: ad.title || 'Без названия',
    seller: ad.seller?.name || 'Продавец',
    sellerId: ad.seller?.id?.toString() || 'seller',
    time: ad.created_at ? new Date(ad.created_at).toLocaleString() : 'недавно',
    link: Number(ad.id) > 0 ? `/advert/${ad.id}` : '/advert/0',
    price: ad.price || '0',
    category: ad.category || 'Без категории',
  });

  if (loading) {
    return (
        <PageLayout logo={logo}>
          <div className="text-content">
            <p>Загрузка объявлений...</p>
          </div>
        </PageLayout>
    );
  }

  if (error) {
    return (
        <PageLayout logo={logo}>
          <div className="text-content">
            <p>Ошибка: {error}</p>
          </div>
        </PageLayout>
    );
  }

  return (
      <PageLayout logo={logo}>
        <div className="text-content">
          <div className="action-buttons">
            <Link to="/add-advert" className="create-ad-button main-page-button">
              Создать объявление
            </Link>
          </div>

          <MarkdownContent content={`# Добро пожаловать!`} />

          <h3>Горячие предложения:</h3>
          {hotAdverts.length > 0 ? (
              <AdvertList adverts={hotAdverts.map(normalizeAd)} />
          ) : (
              <p>Нет доступных объявлений</p>
          )}

          {viewedAds.length > 0 && (
              <>
                <h3 style={{ marginTop: '30px' }}>Вы недавно смотрели:</h3>
                <AdvertList adverts={viewedAds.map(normalizeAd)} />
              </>
          )}

          <h5 style={{ marginTop: '40px' }}>
            <Link to="/author">Об авторе</Link>
          </h5>
        </div>
      </PageLayout>
  );
};

export default MainPage;