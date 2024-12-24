import React from "react"
import ReactDOM from "react-dom/client"
import NginxIgnition from "./domain/NginxIgnition"
import "@ant-design/v5-patch-for-react-19"

const rootElement = document.getElementById("nginx-ignition-root") as HTMLElement
const reactRoot = ReactDOM.createRoot(rootElement)
const nginxIgnition = <NginxIgnition />
reactRoot.render(nginxIgnition)
