import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';


import $ from 'jquery';
import Popper from 'popper.js';
import 'bootstrap/dist/js/bootstrap.bundle.min';
import 'bootstrap/dist/css/bootstrap.min.css';
import "./App.css"

import 'moment/locale/de';
import 'moment/locale/de';
import * as serviceWorker from './serviceWorker';
import moment from "moment";



import { Provider } from 'react-redux'
import {applyMiddleware, createStore} from "redux";

import createSagaMiddleware from 'redux-saga'
import {createLogger} from 'redux-logger';
import {MainReducer} from "./redux/reducers/MainReducer";
import {rootSaga} from "./redux/effects";


moment.locale("de");


const logger = createLogger();
const sagaMiddleware = createSagaMiddleware();
const store = createStore(MainReducer,applyMiddleware(sagaMiddleware, logger));


sagaMiddleware.run(rootSaga);

ReactDOM.render(<Provider store={store}>
    <App />
</Provider>, document.getElementById('root'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
