import React from "react";
import {HostRoute} from "../model/HostRequest";
import ValidationResult from "../../../core/validation/ValidationResult";

export interface HostRoutesProps {
    routes: HostRoute[]
    validationResult: ValidationResult
}

export default class HostRoutes extends React.Component<HostRoutesProps> {
    render() {
        return "TODO"
    }
}
