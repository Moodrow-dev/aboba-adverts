import { Link } from 'react-router-dom';
import PageLayout from '../components/PageLayout';

const ErrorPage = ({logo}) => {
  return (
    <PageLayout logo={logo}>
      <div className="text centered text-content">
        <h1>Ошибка 404</h1>
        <h2>
          Ой-ой-ой, страничка убежала!<br />
          404 — это когда ты хотел что-то найти, а оно такое: "Пока, я в Воронеже, ищу гараж!" <br />
          Попробуй кнопочку обновить, может, страничка просто спит, как дядя Петя на охране.<br />
          Или возвращайся на <Link to="/">главную</Link> — там всё целое, обещаем!
        </h2>
      </div>
    </PageLayout>
  );
};

export default ErrorPage;