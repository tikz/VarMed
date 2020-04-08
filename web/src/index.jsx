import React from 'react';
import ReactDOM from 'react-dom';
import Results from './components/Results';
import './styles/app.scss';

require('file-loader?name=[name].[ext]!./index.html');

ReactDOM.render(
    <Results />,
    document.getElementById('app')
);