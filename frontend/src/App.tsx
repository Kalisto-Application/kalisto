import {useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {SpecFromProto, SendGrpc} from "../wailsjs/go/api/Api";

function App() {
    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo"/>
            <div id="result" className="result">"hey dude"</div>
        </div>
    )
}

export default App
