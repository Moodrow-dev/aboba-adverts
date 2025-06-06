import ReactMarkdown from 'react-markdown';
import './MarkdownContent.css';

const MarkdownContent = ({ content }) => {
  return (
    <div className="markdown-content">
      <ReactMarkdown>{content}</ReactMarkdown>
    </div>
  );
};

export default MarkdownContent;