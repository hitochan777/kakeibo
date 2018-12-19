import axios from 'axios';

export const postRawText = async text => {
  await axios.post('http://localhost:8080', text, {
    'Content-Type': 'text/plain',
  });
};
