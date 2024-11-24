import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import NginxIgnition from './domain/NginxIgnition';
import reportWebVitals from "./web-vitals";

const reactRoot = ReactDOM.createRoot(document.body)
const nginxIgnition = <NginxIgnition />
reactRoot.render(nginxIgnition);
reportWebVitals();
