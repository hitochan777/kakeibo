import axios from 'axios';

const sleep = milliseconds => {
  return new Promise(resolve => setTimeout(resolve, milliseconds));
};

export const postRawText = async text => {
  console.log(text);
  await sleep(1000); // TODO: post text to converter server
};
