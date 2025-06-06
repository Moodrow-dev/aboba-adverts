import { Link } from 'react-router-dom';
import './AdvertList.css';

const AdvertList = ({ adverts = [] }) => {
  return (
    <div className="advert-list">
      {adverts.map(ad => (
        <div key={ad.id} className="advert-item">
          <h4>
            <Link to={ad.link || '#'} className="advert-title">
              {ad.title || 'Без названия'}
            </Link>
          </h4>
          <div className="advert-meta">
            <Link 
              to={`/seller/${ad.sellerId || 'unknown'}`} 
              className="advert-seller"
            >
              {ad.seller || 'Неизвестный продавец'}
            </Link>
            <span className="advert-time">{ad.time || ''}</span>
          </div>
        </div>
      ))}
    </div>
  );
};

export default AdvertList;