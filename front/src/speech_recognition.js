import {Observable} from 'rxjs';

const SpeechRecognition =
  window.SpeechRecognition || window.webkitSpeechRecognition;

export const createSpeechObservable = () => {
  return Observable.create(observer => {
    const recognition = new SpeechRecognition();
    recognition.continuous = false;
    recognition.interimResults = false;
    recognition.lang = 'ja-JP';
    recognition.onresult = event => {
      return observer.next(event.results[0][0]);
    };
    recognition.onerror = error => {
      observer.error(error);
    };
    recognition.start();
  });
};
