import React, {Component} from 'react';
import './App.css';

import {createSpeechObservable} from './speech_recognition';
import {postRawText} from './api';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      recording: false,
    };
    this.handleClick = this.handleClick.bind(this);
  }

  handleClick() {
    if (this.state.recording) {
      alert('Cannot start multiple recorder at the same time');
      return;
    }
    const observable = createSpeechObservable();
    this.setState({recording: true});
    observable.subscribe(
      async result => {
        console.log(result);
        this.setState({recording: false});
        await postRawText(result.transcript);
      },
      error => {
        console.log(error);
        this.setState({recording: false});
      },
      x => {
        this.setState({recording: false});
      },
    );
  }

  render() {
    return (
      <div className="App">
        <button className="button" onClick={this.handleClick}>
          Kakeibo
          {this.state.recording && 'ing...'}
        </button>
      </div>
    );
  }
}

export default App;
