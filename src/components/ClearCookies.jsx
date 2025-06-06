const ClearCookies = () => {
    const handleClick = () => {
      if (window.confirm('Точно очистить куки?')) {
        document.cookie.split(';').forEach(cookie => {
          const name = cookie.trim().split('=')[0];
          document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/`;
        });
        alert('Куки очищены!');
        window.location.reload();
      }
    };
  
    return (
      <button 
        onClick={handleClick}
        style={{ 
          marginLeft: '10px',
          background: '#ff4d4f',
          color: 'white',
          border: 'none',
          padding: '5px 10px',
          borderRadius: '4px'
        }}
      >
        Очистить куки
      </button>
    );
  };
  export default ClearCookies;