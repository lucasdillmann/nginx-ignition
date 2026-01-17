import React from "react"
import { LoadingOutlined } from "@ant-design/icons"

export default class FullPagePreloader extends React.Component {
    render() {
        return (
            <div className="preloader-container">
                <div className="preloader-body">
                    <div className="preloader-spinner">
                        <LoadingOutlined />
                    </div>
                </div>
            </div>
        )
    }
}
