import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom'; // Добавлен Link
import PageLayout from '../components/PageLayout';

const AddAdvertPage = ({ logo }) => {
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    price: '',
    category: 'garage',
    images: []
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();

    // Валидация полей
    if (!formData.title.trim() || !formData.description.trim() || !formData.price.trim()) {
      alert('Пожалуйста, заполните все обязательные поля');
      return;
    }

    setIsSubmitting(true);

    try {
      const response = await fetch('http://localhost:4242/api/adverts', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          title: formData.title,
          description: formData.description,
          price: formData.price,
          category: formData.category,
          // Добавляем временные ссылки на изображения
          images: formData.images
        })
      });

      const result = await response.json();

      if (!response.ok) {
        throw new Error(result.message || 'Ошибка сохранения');
      }

      navigate('/', {
        state: {
          message: `Объявление "${result.title}" успешно создано!`,
          type: 'success'
        }
      });

    } catch (error) {
      console.error('Ошибка:', error);
      alert(error.message || 'Не удалось сохранить объявление. Проверьте консоль для подробностей.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleImageUpload = (e) => {
    const files = Array.from(e.target.files);
    const imageUrls = files.map(file => URL.createObjectURL(file));
    setFormData(prev => ({
      ...prev,
      images: [...prev.images, ...imageUrls].slice(0, 10)
    }));
  };

  const removeImage = (index) => {
    setFormData(prev => ({
      ...prev,
      images: prev.images.filter((_, i) => i !== index)
    }));
  };

  return (
    <PageLayout logo={logo}>
      <div className="text-content">
        <div className="form-header">
          <h1>Создать объявление</h1>
          <Link to="/" className="back-button">
            ← Назад
          </Link>
        </div>

        <form onSubmit={handleSubmit} className="advert-form">
          <div className="form-group">
            <label>Название*</label>
            <input
              type="text"
              value={formData.title}
              onChange={(e) => setFormData({...formData, title: e.target.value})}
              required
              maxLength={100}
            />
          </div>

          <div className="form-group">
            <label>Описание*</label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData({...formData, description: e.target.value})}
              required
              rows={5}
            />
          </div>

          <div className="form-row">
            <div className="form-group">
              <label>Цена*</label>
              <input
                type="text"
                value={formData.price}
                onChange={(e) => setFormData({...formData, price: e.target.value})}
                required
              />
            </div>
            <div className="form-group">
              <label>Категория*</label>
              <select
                value={formData.category}
                onChange={(e) => setFormData({...formData, category: e.target.value})}
              >
                <option value="garage">Гараж</option>
                <option value="electronics">Электроника</option>
                <option value="other">Другое</option>
              </select>
            </div>
          </div>

          <div className="form-group">
            <label>Фотографии</label>
            <input
              type="file"
              multiple
              accept="image/*"
              onChange={handleImageUpload}
              disabled={formData.images.length >= 10}
            />
            <div className="image-previews">
              {formData.images.map((img, index) => (
                <div key={index} className="image-preview-container">
                  <img src={img} alt={`Preview ${index}`} className="image-preview" />
                  <button
                    type="button"
                    className="remove-image-button"
                    onClick={() => removeImage(index)}
                  >
                    ×
                  </button>
                </div>
              ))}
            </div>
          </div>

          <button 
            type="submit" 
            className="submit-btn"
            disabled={isSubmitting}
          >
            {isSubmitting ? 'Сохранение...' : 'Опубликовать'}
          </button>
        </form>
      </div>
    </PageLayout>
  );
};

export default AddAdvertPage;