import { Link } from 'react-router-dom';

const PageLayout = ({ children, logo }) => {
  return (
    <>
      <header className="site-header">
        <Link to="/" className="logo-link">
          <img className="text-box logo" src={logo} alt="Логотип" />
        </Link>
        <nav className="header-nav">
          <Link to="/add-advert" className="create-ad-button">
            Создать объявление
          </Link>
          <Link to="/profile" className="back-button">
            Личный кабинет
          </Link>
        </nav>
      </header>
      
      <div className="blur-box roboto-regular text-container">
        {children}
      </div>
    </>
  );
};


export default PageLayout;