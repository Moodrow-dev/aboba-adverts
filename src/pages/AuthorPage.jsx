import PageLayout from '../components/PageLayout';

const authorImg = require('../assets/author.jpg');
const AuthorPage = ({logo}) => {
  return (
    <PageLayout logo={logo}>
      <div className="text-content">
        <h1>Об авторе</h1>
        <h3>Макаров Антон Владимирович</h3>
        <img width="500" height="500" src={authorImg} alt="user" className="user__photo" />
        <h3>Статус: Гаражный сеньор</h3>
        <a href="mailto:author@aboba.ru"><button id="contact__btn">Связаться</button></a>
        <br />
        <a href="https://t.me/mynameisasskiss"><button id="contact__btn">Написать</button></a>
        <br/>
        Знакомьтесь, это Антон — автор сайта "Абоба". Он настолько любит гаражи, что спит в обнимку с чертежами и называет свой сайт "гаражным Tinder" — свайп вправо, если хочешь бетонный дворец для своего 'ласточки'. Говорит, что идея сайта родилась, когда он в третий раз за месяц застрял в гараже с пивом и шашлыками, потому что ключ потерял. Теперь он продаёт гаражи всем, кто готов доверить ему свою машину, а заодно и свои нервы. Если что, Антон клянётся, что каждый гараж проверяет лично — сидит внутри с фонариком и ищет, где бы ещё приварить полку для банок с солёными огурцами.
      </div>
    </PageLayout>
  );
};

export default AuthorPage;