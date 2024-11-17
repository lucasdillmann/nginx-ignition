import React from "react";
import {Navigate} from "react-router-dom";

export default class HomePage extends React.PureComponent {
    render() {
        return (
            <Navigate to="/hosts" />
        )
    }
}
