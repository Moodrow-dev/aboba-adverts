import { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { Button, Image, Spin, message } from 'antd';
import Cookies from 'js-cookie';
import PageLayout from '../components/PageLayout';
import MarkdownContent from '../components/MarkdownContent';

const AdvertPage = ({ logo }) => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [advert, setAdvert] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [previewImage, setPreviewImage] = useState('');
  const [previewVisible, setPreviewVisible] = useState(false);

  // Fetch advert data
  useEffect(() => {
    const fetchAdvert = async () => {
      try {
        setLoading(true);
        const response = await fetch(`http://localhost:4242/api/adverts/${id}`);

        if (!response.ok) {
          if (response.status === 404) {
            navigate('/not-found', { replace: true });
            return;
          }
          throw new Error('Ошибка загрузки объявления');
        }

        const data = await response.json();
        console.log('Advert Data:', data); // Debug log
        setAdvert(data);

        // Add to viewed ads
        addToViewedAds(data);
      } catch (err) {
        console.error('Ошибка:', err);
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchAdvert();
  }, [id, navigate]);

  // Delete advert
  const handleDelete = async () => {
    try {
      const response = await fetch(`http://localhost:4242/api/adverts/${id}`, {
        method: 'DELETE',
      });

      if (!response.ok) {
        if (response.status === 404) {
          throw new Error('Объявление не найдено');
        }
        throw new Error('Ошибка удаления объявления');
      }

      message.success('Объявление успешно удалено');
      navigate('/'); // Redirect to home page after deletion
    } catch (err) {
      console.error('Ошибка удаления:', err);
      message.error(err.message);
    }
  };

  // Add advert to viewed ads
  const addToViewedAds = (advertData) => {
    try {
      const currentAds = getViewedAds();

      // Remove duplicates
      const filteredAds = currentAds.filter(ad => ad.id !== advertData.id);

      // Add new advert to the start of the list
      const updatedAds = [
        {
          id: advertData.id || 0,
          title: advertData.title || 'Без названия',
          seller: advertData.seller?.name || 'Продавец',
          sellerId: advertData.seller?.id?.toString() || 'seller',
          price: advertData.price || '0',
          time: new Date().toISOString(),
          image: advertData.images?.[0] || null,
        },
        ...filteredAds,
      ].slice(0, 10); // Limit to 10 items

      Cookies.set('viewed_ads', JSON.stringify(updatedAds), {
        expires: 30, // Store for 30 days
        path: '/',
      });
    } catch (e) {
      console.error('Ошибка сохранения истории просмотров:', e);
    }
  };

  // Get viewed ads from cookies
  const getViewedAds = () => {
    try {
      const viewedAds = Cookies.get('viewed_ads');
      return viewedAds ? JSON.parse(viewedAds) : [];
    } catch (e) {
      console.error('Ошибка чтения куки:', e);
      Cookies.remove('viewed_ads');
      return [];
    }
  };

  if (loading) {
    return (
        <PageLayout logo={logo}>
          <div style={{ display: 'flex', justifyContent: 'center', padding: '50px' }}>
            <Spin size="large" />
          </div>
        </PageLayout>
    );
  }

  if (error) {
    return (
        <PageLayout logo={logo}>
          <div style={{ padding: '20px', textAlign: 'center' }}>
            <h2>Произошла ошибка</h2>
            <p>{error}</p>
            <Button type="primary" onClick={() => window.location.reload()}>
              Попробовать снова
            </Button>
          </div>
        </PageLayout>
    );
  }

  if (!advert) {
    return (
        <PageLayout logo={logo}>
          <div style={{ padding: '20px', textAlign: 'center' }}>
            <h2>Объявление не найдено</h2>
            <Link to="/">
              <Button type="primary">На главную</Button>
            </Link>
          </div>
        </PageLayout>
    );
  }

  return (
      <PageLayout logo={logo}>
        <div className="text-content" style={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
          <h1 style={{ marginBottom: '20px' }}>{advert.title || 'Без названия'}</h1>

          {/* Image gallery */}
          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fill, minmax(200px, 1fr))',
            gap: '10px',
            marginBottom: '20px',
          }}>
            {Array.isArray(advert.images) && advert.images.length > 0 ? (
                advert.images.map((img, index) => (
                    <img
                        key={index}
                        src={img}
                        alt={`Фото объявления ${index + 1}`}
                        style={{
                          width: '100%',
                          height: '200px',
                          objectFit: 'cover',
                          borderRadius: '4px',
                          cursor: 'pointer',
                          boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
                        }}
                        onClick={() => {
                          setPreviewImage(img);
                          setPreviewVisible(true);
                        }}
                        loading="lazy"
                    />
                ))
            ) : (
                <p>Изображения отсутствуют</p>
            )}
          </div>

          {/* Image preview */}
          <Image
              preview={{
                visible: previewVisible,
                onVisibleChange: (visible) => setPreviewVisible(visible),
                src: previewImage,
              }}
              style={{ display: 'none' }}
          />

          {/* Seller info */}
          <div style={{
            backgroundColor: '#f5f5f5',
            padding: '15px',
            borderRadius: '4px',
            marginBottom: '20px',
          }}>
            <h3 style={{ marginBottom: '10px' }}>
              Продавец: <Link to={`/seller/${advert.seller?.id || '1'}`} style={{ color: '#1890ff' }}>
              {advert.seller?.name || 'Продавец'}
            </Link>
            </h3>
            <h3 style={{ marginBottom: '10px' }}>Категория: {advert.category || 'Без категории'}</h3>
            <h3 style={{ marginBottom: '0' }}>Цена: {advert.price || '0'}</h3>
          </div>

          {/* Action buttons */}
          <div style={{ display: 'flex', gap: '10px', marginBottom: '25px' }}>
            <Button
                type="primary"
                size="large"
                style={{ flex: 1, height: '50px', fontSize: '18px' }}
            >
              Купить
            </Button>
            <Button
                type="default"
                size="large"
                style={{ flex: 1, height: '50px', fontSize: '18px' }}
            >
              Написать продавцу
            </Button>
            <Button
                type="danger"
                size="large"
                style={{ flex: 1, height: '50px', fontSize: '18px' }}
                onClick={handleDelete}
            >
              Удалить объявление
            </Button>
          </div>

          {/* Description */}
          <div style={{ marginBottom: '30px' }}>
            <h2 style={{ marginBottom: '15px' }}>Описание</h2>
            <MarkdownContent content={advert.description || 'Описание отсутствует'} />
          </div>

          {/* Additional info */}
          <div style={{
            backgroundColor: '#f9f9f9',
            padding: '15px',
            borderRadius: '4px',
            marginBottom: '30px',
          }}>
            <h3 style={{ marginBottom: '10px' }}>Дополнительная информация</h3>
            <p>Дата публикации: {advert.created_at ? new Date(advert.created_at).toLocaleDateString() : 'Неизвестно'}</p>
            <p>ID объявления: {advert.id || 'Неизвестно'}</p>
          </div>

          {/* Back button */}
          <div style={{ textAlign: 'center' }}>
            <Button type="default" onClick={() => navigate(-1)}>
              Вернуться назад
            </Button>
          </div>
        </div>
      </PageLayout>
  );
};

export default AdvertPage;
