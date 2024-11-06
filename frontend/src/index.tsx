import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import NginxIgnite from './domain/NginxIgnite';
import reportWebVitals from "./web-vitals";

const rootElement = document.getElementById('root') as HTMLElement
const reactRoot = ReactDOM.createRoot(rootElement)
const nginxIgnite = <NginxIgnite />
reactRoot.render(nginxIgnite);
reportWebVitals();
