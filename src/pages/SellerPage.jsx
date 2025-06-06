import { useState, useEffect } from 'react';
import { Link, useParams } from 'react-router-dom';
import PageLayout from '../components/PageLayout';
import AdvertList from '../components/AdvertList';

const userImg = require('../assets/user.jpg');

const SellerPage = ({ logo }) => {
    const { sellerId } = useParams();
    const [seller, setSeller] = useState(null);
    const [adverts, setAdverts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchSellerData = async () => {
            try {
                setLoading(true);

                // Загрузка информации о продавце
                const sellerResponse = await fetch(`http://localhost:4242/api/sellers/${sellerId}`);
                if (!sellerResponse.ok) {
                    throw new Error('Продавец не найден');
                }
                const sellerData = await sellerResponse.json();

                // Загрузка объявлений продавца
                const advertsResponse = await fetch(`http://localhost:4242/api/adverts?sellerId=${sellerId}`);
                if (!advertsResponse.ok) {
                    throw new Error('Ошибка загрузки объявлений');
                }
                const advertsData = await advertsResponse.json();

                setSeller(sellerData);
                setAdverts(advertsData);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchSellerData();
    }, [sellerId]);

    if (loading) {
        return (
            <PageLayout logo={logo}>
                <div className="text-content">
                    <p>Загрузка информации о продавце...</p>
                </div>
            </PageLayout>
        );
    }

    if (error) {
        return (
            <PageLayout logo={logo}>
                <div className="text-content">
                    <h1>Ошибка</h1>
                    <p>{error}</p>
                    <Link to="/">Вернуться на главную</Link>
                </div>
            </PageLayout>
        );
    }

    if (!seller) {
        return (
            <PageLayout logo={logo}>
                <div className="text-content">
                    <h1>Продавец не найден</h1>
                    <Link to="/">Вернуться на главную</Link>
                </div>
            </PageLayout>
        );
    }

    return (
        <PageLayout logo={logo}>
            <div className="text-content">
                <h1>{seller.name} - {seller.title}</h1>
                <img width="25%" src={seller.photo || userImg} alt={seller.name} className="user__photo" />

                <div className="seller-info">
                    <h3>Продавец: {seller.name} ({seller.nickname})</h3>
                    <h3>Статус: {seller.status}</h3>
                    <h3>Город: {seller.city}</h3>
                    <a href={`mailto:${seller.email}`}>
                        <button id="contact__btn">Связаться</button>
                    </a>
                </div>

                <h2>О себе</h2>
                <p id="seller__info">
                    {seller.about}
                </p>

                <h2>Активные объявления</h2>
                {adverts.length > 0 ? (
                    <AdvertList adverts={adverts.map(ad => ({
                        id: ad.id,
                        title: ad.title,
                        seller: seller.name,
                        sellerId: seller.id,
                        time: new Date(ad.created_at).toLocaleString(),
                        link: `/advert/${ad.id}`,
                        price: ad.price,
                        category: ad.category
                    }))} />
                ) : (
                    <p>У этого продавца пока нет активных объявлений</p>
                )}
            </div>
        </PageLayout>
    );
};

export default SellerPage;