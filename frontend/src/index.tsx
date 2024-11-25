import React from 'react';
import ReactDOM from 'react-dom/client';
import NginxIgnition from './domain/NginxIgnition';
import reportWebVitals from "./web-vitals";

const rootElement = document.getElementById('nginx-ignition-root') as HTMLElement
const reactRoot = ReactDOM.createRoot(rootElement)
const nginxIgnition = <NginxIgnition />
reactRoot.render(nginxIgnition);
reportWebVitals();
