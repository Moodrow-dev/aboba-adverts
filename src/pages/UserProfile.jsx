import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Card, Button } from 'antd';
import ClearCookies from '../components/ClearCookies';

const UserProfile = () => {
  const [userData, setUserData] = useState(null);

  useEffect(() => {
    // Загрузка данных пользователя
    const loadUserData = () => {
      try {
        const fakeUserData = {
          name: "Гость",
          email: "guest@example.com",
          adsCount: 0
        };
        setUserData(fakeUserData);
      } catch (error) {
        console.error('Ошибка загрузки данных:', error);
      }
    };

    loadUserData();
  }, []);

  if (!userData) return <div>Загрузка...</div>;

  return (
    <div className="user-profile">
      <Card title="Личный кабинет" className="profile-card">
        <div className="user-info">
          <h2>{userData.name}</h2>
          {userData.email && <p>Email: {userData.email}</p>}
          {userData.adsCount !== undefined && <p>Объявлений: {userData.adsCount}</p>}
        </div>

        <div className="action-buttons">
          <Link to="/add-advert">
            <Button type="primary" className="add-ad-btn">
              Добавить объявление
            </Button>
          </Link>

          <Link to="/my-ads">
            <Button type="link">Мои объявления</Button>
          </Link>

          {/* Компонент очистки кук */}
          <ClearCookies />
        </div>
      </Card>
    </div>
  );
};

export default UserProfile;